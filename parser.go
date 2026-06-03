package bibtex

// Parser based on https://metacpan.org/release/BORISV/BibTeX-Parser-1.04
// Some useful links:
// format description https://maverick.inria.fr/~Xavier.Decoret/resources/xdkbibtex/bibtex_summary.html

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/ugent-library/bibtex/latex"
)

var (
	namePattern       = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+`
	reAtName          = regexp.MustCompile(`@(` + namePattern + `)`)
	reKey             = regexp.MustCompile(`^\s*\{\s*([^\s,]+[^,]*?)\s*,[\s\n]*`)
	reFieldName       = regexp.MustCompile(`^[\s\n]*(` + namePattern + `)[\s\n]*=[\s\n]*`)
	reDigits          = regexp.MustCompile(`^\d+`)
	reName            = regexp.MustCompile(`^` + namePattern)
	reStringName      = regexp.MustCompile(`\{\s*(` + namePattern + `)\s*=\s*`)
	reQuotedString    = regexp.MustCompile(`^"(([^"\\]*(\\.)*[^\\"]*)*)"`)
	reConcatString    = regexp.MustCompile(`^\s*#\s*`)
	reWhitespace      = regexp.MustCompile(`^\s*`)
	reEscape          = regexp.MustCompile(`^\\.`)
	reStringValue     = regexp.MustCompile(`^[^\\\{\}]+`)
	reAuthorEditorSep = regexp.MustCompile(`(?i)\s+and\s+`)
	reAuthorEditor    = regexp.MustCompile(`(?i)(.*?)(\{|\s+and\s+)`)
)

type Parser struct {
	r       *bufio.Reader
	line    int
	strings map[string]string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r:       bufio.NewReader(r),
		line:    1,
		strings: make(map[string]string),
	}
}

// Next returns the next entry in the stream, or (nil, nil) at end of input.
//
// The parser is deliberately forgiving. It skips any junk before an entry (a
// byte-order mark, export headers, blank lines), splits entries purely by
// matching braces (so entries glued together without a separator still parse),
// and best-effort returns a truncated final entry rather than erroring. A
// malformed or missing key yields an empty Key instead of an error.
//
// @string, @comment and @preamble are returned as entries with Type set; an
// @string definition is additionally recorded so that later entries can
// resolve references to it.
func (p *Parser) Next() (*Entry, error) {
	for {
		buf := &strings.Builder{}

		// Skip anything until the next entry marker. EOF here means we're done.
		if err := p.skipUntil(buf, '@'); err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		startLine := p.line

		// Read up to and including the opening brace. An EOF before we find the
		// brace means a truncated marker at end of input; parse what we have.
		if err := p.writeWhitespace(buf); err != nil && err != io.EOF {
			return nil, err
		}
		if err := p.writeUntil(buf, '{'); err != nil && err != io.EOF {
			return nil, err
		}

		// Consume the body by matching braces. EOF mid-body is best-effort: we
		// keep what we read and parse it anyway.
		braceLevel := 1
		for braceLevel != 0 {
			c, err := p.readRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if c == '{' {
				braceLevel++
			} else if c == '}' {
				braceLevel--
			}
			if _, err := buf.WriteRune(c); err != nil {
				return nil, err
			}
		}

		if e := p.parse(buf.String(), startLine); e != nil {
			return e, nil
		}
		// The chunk held no recognisable @type (e.g. a stray '@'); keep scanning.
	}
}

// parse turns a raw "@type{...}" chunk into an Entry, or returns nil when the
// chunk holds no recognisable @type.
func (p *Parser) parse(raw string, line int) *Entry {
	eStr := raw

	m := reAtName.FindStringSubmatchIndex(eStr)
	if m == nil {
		return nil
	}

	e := &Entry{
		Raw:  strings.TrimSpace(raw),
		Type: strings.ToLower(eStr[m[2]:m[3]]),
		Line: line,
	}
	eStr = eStr[m[1]:] // advance past @type

	switch e.Type {
	case "comment", "preamble":
		// Freeform body; expose type and raw only.
		return e
	case "string":
		m = reStringName.FindStringSubmatchIndex(eStr)
		if m == nil {
			return e // malformed @string; best-effort
		}
		key := eStr[m[2]:m[3]]
		_, val := p.parseString(eStr[m[1]:])
		p.strings[key] = val
		e.Key = key
		e.Fields = append(e.Fields, Field{
			Name:     key,
			RawValue: val,
			Value:    latex.Decode(val),
		})
		return e
	}

	// Key. Forgiving: if it doesn't match, leave Key empty but still try to
	// reach the fields by skipping past the opening brace.
	if m = reKey.FindStringSubmatchIndex(eStr); m != nil {
		e.Key = eStr[m[2]:m[3]]
		eStr = eStr[m[1]:]
	} else if i := strings.IndexByte(eStr, '{'); i >= 0 {
		eStr = eStr[i+1:]
	}

	// Fields.
	for m = reFieldName.FindStringSubmatchIndex(eStr); m != nil; m = reFieldName.FindStringSubmatchIndex(eStr) {
		field := Field{Name: strings.ToLower(eStr[m[2]:m[3]])}
		eStr = eStr[m[1]:] // advance past "name ="

		var val string
		eStr, val = p.parseString(eStr)
		field.RawValue = val
		field.Value = latex.Decode(val)
		e.Fields = append(e.Fields, field)

		switch field.Name {
		case "author":
			e.RawAuthors = splitAuthorEditor(val)
			e.Authors = decodeAll(e.RawAuthors)
		case "editor":
			e.RawEditors = splitAuthorEditor(val)
			e.Editors = decodeAll(e.RawEditors)
		}

		// Skip past the field separator.
		if idx := strings.IndexByte(eStr, ','); idx >= 0 {
			eStr = eStr[idx+1:]
		}
	}

	return e
}

func decodeAll(names []string) []string {
	if names == nil {
		return nil
	}
	out := make([]string, len(names))
	for i, name := range names {
		out[i] = latex.Decode(name)
	}
	return out
}

// skipUntil discards runes until char is found, writing char to buf.
func (p *Parser) skipUntil(buf *strings.Builder, char rune) error {
	for {
		c, err := p.readRune()
		if err != nil {
			return err
		}
		if c == char {
			_, err := buf.WriteRune(c)
			return err
		}
	}
}

func (p *Parser) writeWhitespace(buf *strings.Builder) error {
	for {
		c, err := p.readRune()
		if err != nil {
			return err
		}
		if unicode.IsSpace(c) {
			if _, err := buf.WriteRune(c); err != nil {
				return err
			}
		} else {
			return p.r.UnreadRune()
		}
	}
}

func (p *Parser) writeUntil(buf *strings.Builder, char rune) error {
	for {
		c, err := p.readRune()
		if err != nil {
			return err
		}
		if _, err := buf.WriteRune(c); err != nil {
			return err
		}
		if c == char {
			return nil
		}
	}
}

func (p *Parser) readRune() (rune, error) {
	c, _, err := p.r.ReadRune()
	if err != nil {
		return 0, err
	}
	// Count lines on '\n' only, so that "\r\n" counts once.
	if c == '\n' {
		p.line++
	}
	return c, nil
}

func (p *Parser) parseString(eStr string) (string, string) {
	buf := &strings.Builder{}

	for {
		if m := reDigits.FindStringIndex(eStr); m != nil {
			buf.WriteString(eStr[m[0]:m[1]])
			eStr = eStr[m[1]:] // advance
		} else if m := reName.FindStringIndex(eStr); m != nil {
			// a bare word is a reference to an @string macro
			buf.WriteString(p.strings[eStr[m[0]:m[1]]])
			eStr = eStr[m[1]:] // advance
		} else if m := reQuotedString.FindStringSubmatchIndex(eStr); m != nil {
			buf.WriteString(eStr[m[2]:m[3]])
			eStr = eStr[m[1]:] // advance
		} else {
			var val string
			eStr, val = p.extractBracketedValue(eStr)
			buf.WriteString(val)
		}

		// resolve "a" # "b" concatenation
		if m := reConcatString.FindStringIndex(eStr); m != nil {
			eStr = eStr[m[1]:] // advance
			continue
		}
		break
	}

	return eStr, buf.String()
}

// extractBracketedValue reads a {braced} value and returns its content with the
// single outermost brace pair removed, preserving any inner (case-protecting)
// braces. So "{168--198}" -> "168--198" and "{{IBM} Research}" -> "{IBM} Research".
func (p *Parser) extractBracketedValue(eStr string) (string, string) {
	var val strings.Builder

	// skip leading whitespace
	if m := reWhitespace.FindStringIndex(eStr); m != nil {
		eStr = eStr[m[1]:]
	}

	depth := 0
	for {
		if m := reEscape.FindStringIndex(eStr); m != nil {
			val.WriteString(eStr[m[0]:m[1]])
			eStr = eStr[m[1]:] // advance
		} else if strings.HasPrefix(eStr, "{") {
			depth++
			if depth > 1 { // keep inner braces, drop the outermost open
				val.WriteByte('{')
			}
			eStr = eStr[1:] // advance
		} else if strings.HasPrefix(eStr, "}") {
			depth--
			if depth > 0 { // keep inner braces, drop the outermost close
				val.WriteByte('}')
			}
			eStr = eStr[1:] // advance
			if depth <= 0 {
				break
			}
		} else if m := reStringValue.FindStringIndex(eStr); m != nil {
			val.WriteString(eStr[m[0]:m[1]])
			eStr = eStr[m[1]:] // advance
		} else {
			break
		}
	}

	return eStr, val.String()
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
			// best-effort: on unbalanced braces keep what we parsed so far
			if numBraces != 0 {
				break
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
