package latex

// Based on https://metacpan.org/release/FIRMICUS/LaTeX-Decode-0.05
// See also https://metacpan.org/release/BORISV/LaTeX-ToUnicode-0.54

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

	// need init
	diacPattern string
	reMacros    *regexp.Regexp
	reDiac1     *regexp.Regexp
	reDiac2     *regexp.Regexp
)

// TODO use code generation
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

// TODO superscript, dings, negations
func Decode(str string) string {
	str = reNormalize1.ReplaceAllString(str, "$1{}$2")
	str = reNormalize2.ReplaceAllString(str, "$1{}$2")
	str = reNormalize3.ReplaceAllString(str, "$1{i}")
	// TODO do we need this?
	//remove {} around macros that print one character
	// by default we skip that, as it would break constructions like \foo{\i}
	// if ($strip_outer_braces) {
	//     $text =~ s/ \{\\($WORDMAC_RE)\} / $WORDMAC{$1} /gxe;
	// }

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

	return str
}
