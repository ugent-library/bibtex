package bibtex

import (
	"regexp"
	"strings"

	"github.com/ugent-library/bibtex/latex"
)

var reSplitAuthorEditor = regexp.MustCompile(`(?i)\s+and\s+`)
var reAuthorEditor = regexp.MustCompile(`(?i)(.*?)(\{|\s+and\s+)`)

type Entry struct {
	Raw    string
	Type   string
	Key    string
	Fields []Field
}

func (e *Entry) Author() []string {
	for _, f := range e.Fields {
		if f.Name == "author" {
			return splitAuthorEditor(f.Value)
		}
	}
	return nil
}

func (e *Entry) Editor() []string {
	for _, f := range e.Fields {
		if f.Name == "editor" {
			return splitAuthorEditor(f.Value)
		}
	}
	return nil
}

type Field struct {
	Name  string
	Value string
}

func (f Field) DecodeValue() string {
	return latex.Decode(f.Value)
}

// Split the str using reSplitAuthorEditor as a delimiter with
// each part having balanced braces (so the pattern
// does NOT split).
// Return empty list if unmatched braces
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

		if reSplitAuthorEditor.MatchString(secondMatch) {
			buf += firstMatch
			tokens = append(tokens, buf)
			buf = ""
		} else if strings.Contains(secondMatch, "{") {
			buf += firstMatch
			buf += "{"
			numBraces := 1
			for numBraces != 0 && str != "" {
				sym := str[0:1]
				buf += sym
				if sym == "{" {
					numBraces++
				} else if sym == "}" {
					numBraces--
				}
				str = str[1:]
			}
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
