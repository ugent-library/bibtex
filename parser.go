package bibtex

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

var (
	reStripComment = regexp.MustCompile(`^%.*$`)
	namePattern    = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+`
	inKeyPattern   = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\' ]+` // allow spaces in identifiers (scopus)
	reAtName       = regexp.MustCompile(`@(` + namePattern + `)`)
	// TODO match this more efficiently
	reKey           = regexp.MustCompile(`s*\{\s*(` + namePattern + inKeyPattern + namePattern + `)\s*,[\s\n]*|\s+\r?\s*`)
	reField         = regexp.MustCompile(`[\s\n]*(` + namePattern + `)[\s\n]*=[\s\n]*`)
	reDigits        = regexp.MustCompile(`^\d+`)
	reName          = regexp.MustCompile(`^` + namePattern)
	reQuotedString  = regexp.MustCompile(`^"(([^"\\]*(\\.)*[^\\"]*)*)"`)
	reConcatString  = regexp.MustCompile(`^\s*#\s*`)
	reWhitespace    = regexp.MustCompile(`^\s*`)
	reEscape        = regexp.MustCompile(`^\\.`)
	reStringVal     = regexp.MustCompile(`^[^\\\{\}]+`) // TODO better name
	reLeftBrackets  = regexp.MustCompile(`^\{+`)
	reRightBrackets = regexp.MustCompile(`\}+$`)
)

type Parser struct {
	scanner *bufio.Scanner
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *Parser) Next() (*Entry, error) {
	scanner := p.scanner
	buf := strings.Builder{}

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
		e.Type = strings.ToUpper(eStr[m[2]:m[3]])

		// read rest of entry (matches braces)
		startPos := m[0] - 1
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

		// add pre to entry
		// TODO ignore pre?
		e.Pre = strings.TrimSpace(eStr[:startPos])

		// raw bibtex
		e.Raw = strings.TrimSpace(eStr[startPos:])

		// advance
		eStr = eStr[m[1]:]

		// TOOD handle STRING

		// only add raw COMMENT and PREAMBLE
		if e.Type == "COMMENT" || e.Type == "PREAMBLE" {
			return e, nil
		}

		m = reKey.FindStringSubmatchIndex(eStr)
		if m == nil {
			// include more info (see perl)
			// TODO slurp close bracket
			return nil, errors.New("malformed entry")
		}

		e.Key = eStr[m[2]:m[3]]

		// advance
		eStr = eStr[m[1]:]
		for m = reField.FindStringSubmatchIndex(eStr); m != nil; m = reField.FindStringSubmatchIndex(eStr) {
			field := Field{Name: strings.ToLower(eStr[m[2]:m[3]])}
			// advance
			eStr = eStr[m[1]:]
			newEStr, val, err := p.parseString(eStr)
			if err != nil {
				return nil, err
			}
			eStr = newEStr
			field.Value = val
			e.Fields = append(e.Fields, field)
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
	str := ""

	for {
		if m := reDigits.FindStringIndex(eStr); m != nil {
			str += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
		} else if m := reName.FindStringIndex(eStr); m != nil {
			// TODO look up string in strings map
			str += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
		} else if m := reQuotedString.FindStringSubmatchIndex(eStr); m != nil {
			str += eStr[m[2]:m[3]]
			eStr = eStr[m[1]:]
		} else {
			newEStr, val := p.extractBracketedValue(eStr)
			// TODO remove brackets
			str += val
			eStr = newEStr
		}

		if m := reConcatString.FindStringIndex(eStr); m != nil {
			eStr = eStr[m[1]:]
			continue
		}
		break
	}

	// TODO replace newlines? see perl

	return eStr, str, nil
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
			eStr = eStr[m[1]:]
			continue
		} else if strings.HasPrefix(eStr, "{") {
			val += "{"
			eStr = eStr[1:]
			depth++
			continue
		} else if strings.HasPrefix(eStr, "}") {
			val += "}"
			eStr = eStr[1:]
			depth--
			if depth > 0 {
				continue
			}
			break
		} else if m := reStringVal.FindStringIndex(eStr); m != nil {
			val += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
			continue
		}
		break
	}

	val = reLeftBrackets.ReplaceAllString(val, "")
	val = reRightBrackets.ReplaceAllString(val, "")

	return eStr, val
}
