package bibtex

// Parser based on https://metacpan.org/release/BORISV/BibTeX-Parser-1.04
// Some useful links:
// format description https://maverick.inria.fr/~Xavier.Decoret/resources/xdkbibtex/bibtex_summary.html

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/ugent-library/bibtex/latex"
)

var (
	reStripComment = regexp.MustCompile(`^%.*$`)
	namePattern    = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+`
	inKeyPattern   = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\' ]+` // we allow spaces in identifiers (scopus)
	reAtName       = regexp.MustCompile(`@(` + namePattern + `)`)
	// TODO match this more efficiently
	reKey             = regexp.MustCompile(`s*\{\s*(` + namePattern + inKeyPattern + namePattern + `)\s*,[\s\n]*|\s+\r?\s*`)
	reFieldName       = regexp.MustCompile(`[\s\n]*(` + namePattern + `)[\s\n]*=[\s\n]*`)
	reDigits          = regexp.MustCompile(`^\d+`)
	reName            = regexp.MustCompile(`^` + namePattern)
	reStringName      = regexp.MustCompile(`\{\s*(` + namePattern + `)\s*=\s*`)
	reQuotedString    = regexp.MustCompile(`^"(([^"\\]*(\\.)*[^\\"]*)*)"`)
	reConcatString    = regexp.MustCompile(`^\s*#\s*`)
	reWhitespace      = regexp.MustCompile(`^\s*`)
	reEscape          = regexp.MustCompile(`^\\.`)
	reStringValue     = regexp.MustCompile(`^[^\\\{\}]+`)
	reLeftBrackets    = regexp.MustCompile(`^\{+`)
	reRightBrackets   = regexp.MustCompile(`\}+$`)
	reAuthorEditorSep = regexp.MustCompile(`(?i)\s+and\s+`)
	reAuthorEditor    = regexp.MustCompile(`(?i)(.*?)(\{|\s+and\s+)`)
)

type Parser struct {
	scanner *bufio.Scanner
	strings map[string]string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		strings: make(map[string]string),
	}
}

func (p *Parser) Next() (*Entry, error) {
	scanner := p.scanner
	buf := &strings.Builder{}

	for scanner.Scan() {
		line := reStripComment.ReplaceAllString(scanner.Text(), "")
		if line != "" {
			buf.WriteString(line + "\n")
		}
		if !strings.Contains(line, "@") {
			continue
		}

		e := &Entry{}

		eStr := buf.String()

		// get type
		m := reAtName.FindStringSubmatchIndex(eStr)
		if m == nil {
			// include more info (see perl)
			return nil, errors.New("type not found")
		}
		e.Type = strings.ToLower(eStr[m[2]:m[3]])

		// read rest of entry (matches braces)
		startPos := m[0] - 1
		// TODO why?
		if startPos < 0 {
			startPos = 0
		}
		// count braces
		braceLevel := strings.Count(eStr, "{") - strings.Count(eStr, "}")

		for braceLevel != 0 {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			braceLevel = braceLevel + strings.Count(line, "{") - strings.Count(line, "}")
			buf.WriteString(line + "\n")
		}

		eStr = buf.String()

		// raw bibtex
		e.Raw = strings.TrimSpace(eStr[startPos:])

		eStr = eStr[m[1]:] // advance

		// skip @comment and @preamble
		if e.Type == "comment" || e.Type == "preamble" {
			return p.Next()
		}

		// handle @string
		if e.Type == "string" {
			m = reStringName.FindStringSubmatchIndex(eStr)
			if m == nil {
				return nil, errors.New("malformed string") // TODO include more info
			}
			key := eStr[m[2]:m[3]]

			eStr = eStr[m[1]:] // advance
			_, val, err := p.parseString(eStr)
			if err != nil {
				return nil, err
			}

			p.strings[key] = val

			return p.Next()
		}

		// handle normal entry
		m = reKey.FindStringSubmatchIndex(eStr)
		if m == nil {
			// include more info (see perl)
			// TODO slurp close bracket
			return nil, errors.New("malformed entry")
		}

		e.Key = eStr[m[2]:m[3]]

		eStr = eStr[m[1]:] // advance
		for m = reFieldName.FindStringSubmatchIndex(eStr); m != nil; m = reFieldName.FindStringSubmatchIndex(eStr) {
			field := Field{Name: strings.ToLower(eStr[m[2]:m[3]])}
			eStr = eStr[m[1]:] // advance
			newEStr, val, err := p.parseString(eStr)
			if err != nil {
				return nil, err
			}
			eStr = newEStr
			field.RawValue = val
			field.Value = latex.Decode(val)
			e.Fields = append(e.Fields, field)

			if field.Name == "author" {
				e.RawAuthors = splitAuthorEditor(val)
				e.Authors = make([]string, len(e.RawAuthors))
				for i, name := range e.RawAuthors {
					e.Authors[i] = latex.Decode(name)
				}
			} else if field.Name == "editor" {
				e.RawEditors = splitAuthorEditor(val)
				e.Editors = make([]string, len(e.RawEditors))
				for i, name := range e.RawEditors {
					e.Editors[i] = latex.Decode(name)
				}
			}

			// skip past next comma
			if idx := strings.Index(eStr, ","); idx > -1 {
				eStr = eStr[idx+1:]
			}
		}

		return e, nil
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *Parser) parseString(eStr string) (string, string, error) {
	buf := &strings.Builder{}

	for {
		if m := reDigits.FindStringIndex(eStr); m != nil {
			buf.WriteString(eStr[m[0]:m[1]])
			eStr = eStr[m[1]:] //advance
		} else if m := reName.FindStringIndex(eStr); m != nil {
			key := eStr[m[0]:m[1]]
			buf.WriteString(p.strings[key])
			eStr = eStr[m[1]:] // advance
		} else if m := reQuotedString.FindStringSubmatchIndex(eStr); m != nil {
			buf.WriteString(eStr[m[2]:m[3]])
			eStr = eStr[m[1]:] //advance
		} else {
			newEStr, val := p.extractBracketedValue(eStr)
			buf.WriteString(val)
			eStr = newEStr
		}

		if m := reConcatString.FindStringIndex(eStr); m != nil {
			eStr = eStr[m[1]:] //advance
			continue
		}
		break
	}

	// TODO replace newlines? see perl

	return eStr, buf.String(), nil
}

func (p *Parser) extractBracketedValue(eStr string) (string, string) {
	val := ""

	// skip whitespace
	if m := reWhitespace.FindStringIndex(eStr); m != nil {
		eStr = eStr[m[1]:]
	}

	var depth int
	for {
		if m := reEscape.FindStringIndex(eStr); m != nil {
			val += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:] // advance
			continue
		} else if strings.HasPrefix(eStr, "{") {
			val += "{"
			eStr = eStr[1:] // advance
			depth++
			continue
		} else if strings.HasPrefix(eStr, "}") {
			val += "}"
			eStr = eStr[1:] // advance
			depth--
			if depth > 0 {
				continue
			}
			break
		} else if m := reStringValue.FindStringIndex(eStr); m != nil {
			val += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:] // advance
			continue
		}
		break
	}

	// remove brackets
	val = reLeftBrackets.ReplaceAllString(val, "")
	val = reRightBrackets.ReplaceAllString(val, "")

	return eStr, val
}

func splitAuthorEditor(str string) []string {
	str = strings.TrimSpace(str)

	var tokens []string

	buf := ""
	for str != "" {
		m := reAuthorEditor.FindStringSubmatchIndex(str)

		if m == nil {
			buf += str
			break
		}

		firstMatch := str[m[2]:m[3]]
		secondMatch := str[m[4]:m[5]]
		str = str[m[1]:] // advance

		if reAuthorEditorSep.MatchString(secondMatch) {
			buf += firstMatch
			tokens = append(tokens, buf)
			buf = ""
		} else if strings.Contains(secondMatch, "{") {
			buf += firstMatch
			buf += "{"
			numBraces := 1
			for numBraces != 0 && str != "" {
				sym := str[0:1] // peek
				buf += sym
				if sym == "{" {
					numBraces++
				} else if sym == "}" {
					numBraces--
				}
				str = str[1:] // advance
			}
			// return nil when braces are unbalanced
			if numBraces != 0 {
				return nil
			}
		} else {
			buf += firstMatch
		}
	}

	if buf != "" {
		tokens = append(tokens, buf)
	}

	return tokens
}
