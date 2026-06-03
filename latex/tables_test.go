package latex

import "testing"

// TestNoTableKeyCollisions guards the invariant that a control-sequence name
// lives in at most one table. Decode applies macros before accents and
// diacritics, so a name shared between tables would make the accent/diacritic
// entry unreachable (this is what happened with "textpalhookbelow", now
// dropped from diacritics). A straight translation of the upstream tables is
// the kind of edit that can reintroduce such a clash.
func TestNoTableKeyCollisions(t *testing.T) {
	tables := []struct {
		name string
		m    map[string]string
	}{
		{"accents", accents},
		{"diacritics", diacritics},
		{"macros", macros},
	}
	for i := range tables {
		for j := i + 1; j < len(tables); j++ {
			for k := range tables[i].m {
				if _, ok := tables[j].m[k]; ok {
					t.Errorf("key %q is in both %s and %s; one shadows the other", k, tables[i].name, tables[j].name)
				}
			}
		}
	}
}

// TestEveryTableEntryDecodes probes that every table entry is actually
// reachable through Decode and produces its mapped value. It is the regression
// net for the "is everything used?" audit: any entry that becomes unreachable
// (a clash, a prefix shadow, a bad escape) fails here, named explicitly.
func TestEveryTableEntryDecodes(t *testing.T) {
	for k, v := range accents {
		// \^e, \'e, ... -> e + combining accent
		if in, want := `\`+k+"e", "e"+v; Decode(in) != want {
			t.Errorf("accent %q: Decode(%q) = %q, want %q", k, in, Decode(in), want)
		}
	}
	for k, v := range diacritics {
		// \r{e}, \capitalring{e}, ... -> e + combining diacritic
		if in, want := `\`+k+"{e}", "e"+v; Decode(in) != want {
			t.Errorf("diacritic %q: Decode(%q) = %q, want %q", k, in, Decode(in), want)
		}
	}
	for k, v := range macros {
		// \ss{}, \euro{}, ... -> the mapped character
		if in := `\` + k + "{}"; Decode(in) != v {
			t.Errorf("macro %q: Decode(%q) = %q, want %q", k, in, Decode(in), v)
		}
	}
}
