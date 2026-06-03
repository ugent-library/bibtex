package bibtex

import (
	"strings"
	"unicode"

	"github.com/ugent-library/bibtex/latex"
)

type Entry struct {
	Raw        string   `json:"-"`
	Type       string   `json:"type"`
	Key        string   `json:"key"`
	Line       int      `json:"line,omitempty"` // source line where the entry starts
	Fields     []Field  `json:"fields,omitempty"`
	RawAuthors []string `json:"-"`
	Authors    []string `json:"authors,omitempty"`
	RawEditors []string `json:"-"`
	Editors    []string `json:"editors,omitempty"`
}

type Field struct {
	Name     string `json:"name"`
	RawValue string `json:"-"`
	Value    string `json:"value"`
}

// Name holds the components of a single author or editor name following the
// BibTeX name grammar.
type Name struct {
	First string `json:"first,omitempty"`
	Von   string `json:"von,omitempty"`
	Last  string `json:"last,omitempty"`
	Jr    string `json:"jr,omitempty"`
}

// SplitName parses a single name (one element of an author/editor list, after
// it has been split on " and ") into its First/von/Last/Jr parts.
//
// It recognises the three BibTeX name forms, distinguished by how many
// top-level (brace depth 0) commas the name contains:
//
//	"First von Last"        (no comma)
//	"von Last, First"       (one comma)
//	"von Last, Jr, First"   (two commas)
//
// The "von" part is the run of tokens delimited by lowercase-initial tokens;
// the final Last token is never absorbed into von. Braces protect their
// contents, so "{Barnes and Noble}" stays a single Last token and a braced
// group is treated as uppercase for von detection.
func SplitName(name string) Name {
	sections := splitAtDepth0(name, ',')
	for i := range sections {
		sections[i] = strings.TrimSpace(sections[i])
	}

	switch len(sections) {
	case 1: // "First von Last"
		tokens := tokenizeName(sections[0])
		n := len(tokens)
		if n == 0 {
			return Name{}
		}
		if n == 1 {
			return Name{Last: tokens[0]}
		}
		firstVon, lastVon := -1, -1
		for i := 0; i < n-1; i++ { // the last token is always Last
			if isVonToken(tokens[i]) {
				if firstVon == -1 {
					firstVon = i
				}
				lastVon = i
			}
		}
		if firstVon == -1 {
			return Name{
				First: join(tokens[:n-1]),
				Last:  tokens[n-1],
			}
		}
		return Name{
			First: join(tokens[:firstVon]),
			Von:   join(tokens[firstVon : lastVon+1]),
			Last:  join(tokens[lastVon+1:]),
		}

	case 2: // "von Last, First"
		von, last := splitVonLast(tokenizeName(sections[0]))
		return Name{Von: von, Last: last, First: sections[1]}

	default: // "von Last, Jr, First" (extra commas fold into First)
		von, last := splitVonLast(tokenizeName(sections[0]))
		return Name{
			Von:   von,
			Last:  last,
			Jr:    sections[1],
			First: strings.Join(sections[2:], ", "),
		}
	}
}

// splitVonLast splits the "von Last" portion of a comma form: von is every
// token up to and including the last lowercase-initial token, the rest is Last.
func splitVonLast(tokens []string) (von, last string) {
	n := len(tokens)
	if n == 0 {
		return "", ""
	}
	if n == 1 {
		return "", tokens[0]
	}
	lastVon := -1
	for i := 0; i < n-1; i++ { // the last token is always Last
		if isVonToken(tokens[i]) {
			lastVon = i
		}
	}
	if lastVon == -1 {
		return "", join(tokens)
	}
	return join(tokens[:lastVon+1]), join(tokens[lastVon+1:])
}

// isVonToken reports whether a token belongs in the von part: its first real
// letter is lowercase. A token wrapped in braces is protected (treated as
// uppercase). The token is LaTeX-decoded for the test so accented forms like
// "\v{S}imon" (Š) or "\O" (Ø) are judged on the resolved letter, not the macro
// name — "\v" or "\c" would otherwise read as a lowercase 'v'/'c'.
func isVonToken(tok string) bool {
	if strings.HasPrefix(tok, "{") {
		return false
	}
	for _, r := range latex.Decode(tok) {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}
	return false
}

// splitAtDepth0 splits s on sep, ignoring separators inside braces.
func splitAtDepth0(s string, sep byte) []string {
	var parts []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		case sep:
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, s[start:])
}

// tokenizeName splits a name section into whitespace-separated tokens, treating
// the TeX tie "~" as whitespace and never splitting inside braces.
func tokenizeName(s string) []string {
	var tokens []string
	depth, start := 0, -1
	flush := func(end int) {
		if start != -1 {
			tokens = append(tokens, s[start:end])
			start = -1
		}
	}
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '{':
			if start == -1 {
				start = i
			}
			depth++
		case '}':
			if start == -1 {
				start = i
			}
			if depth > 0 {
				depth--
			}
		case ' ', '\t', '\n', '\r', '~':
			if depth == 0 {
				flush(i)
				continue
			}
			if start == -1 {
				start = i
			}
		default:
			if start == -1 {
				start = i
			}
		}
	}
	flush(len(s))
	return tokens
}

func join(tokens []string) string {
	return strings.Join(tokens, " ")
}
