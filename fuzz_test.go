package bibtex

import (
	"strings"
	"testing"
)

// FuzzNext asserts the forgiving parser never panics or hangs on arbitrary
// input: every call must eventually return (nil, nil) or an error. Robustness
// on malformed input is the whole point of this package.
func FuzzNext(f *testing.F) {
	seeds := []string{
		"",
		`@article{k, title = {t}}`,
		`@article{k, title = {t}}@misc{j, year = 2020}`,
		`junk before @string{x = "y"} @book{b, author = {A and B}}`,
		`@misc{trunc, title = {no closing brace`,
		`@@@{{{,,,###"""}}}`,
		"@article{í, a = b # c}",
		`@string{a = "x"} @x{k, f = a # a # a}`,
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		p := NewParser(strings.NewReader(string(data)))
		// Bound the loop so a regression that fails to make progress fails the
		// test instead of hanging the fuzzer.
		for i := 0; i < 1_000_000; i++ {
			e, err := p.Next()
			if err != nil || e == nil {
				return
			}
		}
		t.Fatalf("Next did not terminate on %q", data)
	})
}
