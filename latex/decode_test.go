package latex

import "testing"

// TestDecode pins the behaviour of Decode on representative LaTeX constructs.
// Accents and diacritics are emitted in decomposed (combining-mark) form, so
// the expected values use explicit \u escapes for the marks rather than
// precomposed characters; this keeps the test stable regardless of how an
// editor might normalise the source file.
func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text untouched", `plain text`, "plain text"},
		{"acute accent", `\'e`, "é"},
		{"umlaut", `\"o`, "ö"},
		{"braced accented letter loses braces", `{\'e}`, "é"},
		{"accent over braced letter", `\"{o}`, "ö"},
		{"circumflex over dotless i", `\^{\i}`, "î"},
		{"macron over dotless i", `\={\i}`, "ī"}, // the \={\i} special case
		{"cedilla", `\c{c}`, "ç"},
		{"caron", `\v{s}`, "š"},
		{"breve", `\u{a}`, "ă"},
		{"hungarumlaut", `\H{o}`, "ő"},
		{"eszett macro", `\ss`, "ß"},
		{"ae macro", `\ae`, "æ"},
		{"ring-A macro", `\AA`, "Å"},
		{"unknown macro untouched", `\textbf`, `\textbf`},
		{"accent inside a word", `Sch\"onberg`, "Schönberg"},
		{"cedilla inside a word", `Cura\c{c}ao`, "Curaçao"},
		{"accent on a capital", `\'Angel`, "Ángel"},
		{"protective braces around a word are kept", `{IBM} Research`, "{IBM} Research"},
		{"en dash in a page range", `168--198`, "168–198"},
		{"em dash", `foo---bar`, "foo—bar"},
		{"single hyphen untouched", `well-known`, "well-known"},
		{"em then hyphen for four", `a----b`, "a—-b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.in); got != tt.want {
				t.Errorf("Decode(%q) = %q (% x), want %q (% x)",
					tt.in, got, []byte(got), tt.want, []byte(tt.want))
			}
		})
	}
}
