// Package latex converts the subset of LaTeX markup that appears in BibTeX
// field values into Unicode: accents (\'e -> é), diacritics (\v{S} -> Š),
// named macros and symbols (\alpha -> α, \textendash -> –), and the --/---
// ligatures (en/em dashes). Unrecognised macros are passed through unchanged.
//
// It is intentionally scoped to bibliographic input and is not a general LaTeX
// renderer. It does NOT handle:
//
//   - math mode: $...$ and \(...\) are left verbatim, including super- and
//     subscripts such as $^{13}$C or $T_c$;
//   - \textsuperscript and \textsubscript;
//   - \ding{} dingbats and other symbol or decorative fonts;
//   - \not-composed negations (\not\equiv); only precomposed forms present in
//     the macro table (e.g. \neq -> ≠) are converted;
//   - context awareness: the --/--- ligatures are applied everywhere, so "--"
//     inside a url or doi field is also turned into a dash.
//
// Based on https://metacpan.org/release/FIRMICUS/LaTeX-Decode-0.05
// See also https://metacpan.org/release/BORISV/LaTeX-ToUnicode-0.54
package latex

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reNormalize1           = regexp.MustCompile(`(\\[a-zA-Z]+)\\(\s+)`) // \foo\ bar -> \foo{} bar
	reNormalize2           = regexp.MustCompile(`([^{]\\\w)([;,.:%])`)  //} Aaaa\o, -> Aaaa\o{},
	accentsPattern         = `[\^\.` + "`" + `'"~=]`
	baseDiacPattern        = `r|b|B|c|d|G|H|k|M|t|u|v`
	reNormalize3           = regexp.MustCompile(`(\\(?:` + baseDiacPattern + `|` + accentsPattern + `))\{\\i\}`) // special cases such as '\={\i}' -> '\={i}' -> "i\x{304}"
	reAccents1             = regexp.MustCompile(`\\(` + accentsPattern + `)\{(\p{L}\p{M}*)\}`)
	reAccents2             = regexp.MustCompile(`\\(` + accentsPattern + `)(\p{L}\p{M}*)`)
	reBracedAccentedLetter = regexp.MustCompile(`{(\PM\pM+)}`)
	reEmDash               = regexp.MustCompile(`---`) // em dash; matched before en dash
	reEnDash               = regexp.MustCompile(`--`)

	// need init
	diacPattern string
	reMacros    *regexp.Regexp
	reDiac1     *regexp.Regexp
	reDiac2     *regexp.Regexp
)

func init() {
	// diacritics
	diacNames := make([]string, 0, len(diacritics))
	for k := range diacritics {
		diacNames = append(diacNames, k)
	}
	sort.Slice(diacNames, func(i, j int) bool {
		return len(diacNames[i]) > len(diacNames[j])
	})
	diacPattern = strings.Join(diacNames, "|")

	reDiac1 = regexp.MustCompile(`\\(` + diacPattern + `)\s*\{(\p{L}\p{M}*)\}`)
	reDiac2 = regexp.MustCompile(`\\(` + diacPattern + `)\s+(\p{L}\p{M}*)`)

	// macros
	macroNames := make([]string, 0, len(macros))
	for k := range macros {
		macroNames = append(macroNames, k)
	}
	sort.Slice(macroNames, func(i, j int) bool {
		return len(macroNames[i]) > len(macroNames[j])
	})
	reMacros = regexp.MustCompile(`\\(` + strings.Join(macroNames, "|") + `)(?:\{\}|\s+|\b)`)
}

func Decode(str string) string {
	str = reNormalize1.ReplaceAllString(str, "$1{}$2")
	str = reNormalize2.ReplaceAllString(str, "$1{}$2")
	str = reNormalize3.ReplaceAllString(str, "$1{i}")

	str = reMacros.ReplaceAllStringFunc(str, func(macro string) string {
		m := reMacros.FindStringSubmatch(macro)
		return macros[m[1]]
	})

	// run twice
	for i := 0; i < 2; i++ {
		str = reAccents1.ReplaceAllStringFunc(str, func(s string) string {
			m := reAccents1.FindStringSubmatch(s)
			return m[2] + accents[m[1]]
		})
		str = reAccents2.ReplaceAllStringFunc(str, func(s string) string {
			m := reAccents2.FindStringSubmatch(s)
			return m[2] + accents[m[1]]
		})
		str = reDiac1.ReplaceAllStringFunc(str, func(s string) string {
			m := reDiac1.FindStringSubmatch(s)
			return m[2] + diacritics[m[1]]
		})
		str = reDiac2.ReplaceAllStringFunc(str, func(s string) string {
			m := reDiac2.FindStringSubmatch(s)
			return m[2] + diacritics[m[1]]
		})
	}

	// remove {} around letter+combining mark(s)
	// the perl version skips this by default, because it destroys constructions like \foo{\`e}
	str = reBracedAccentedLetter.ReplaceAllString(str, "$1")

	// --/--- ligatures -> en/em dashes (e.g. page ranges "168--198"). Em first
	// so "---" is not consumed as "--" + "-".
	str = reEmDash.ReplaceAllString(str, "—")
	str = reEnDash.ReplaceAllString(str, "–")

	return str
}
