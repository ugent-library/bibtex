package bibtex

import (
	"testing"
)

// TestExamples is a smoke test over the real-world .bib files in examples/.
// It guards that the parser keeps consuming actual exports end to end without
// erroring, returns the expected number of entries, and gives every entry a
// type. The exact field contents are pinned by TestWeird instead.
func TestExamples(t *testing.T) {
	files := []struct {
		File       string
		NumEntries int
	}{
		{"examples/biblatex.bib", 100}, // 92 records + 8 @string entries
		{"examples/biblio.bib", 23},
		{"examples/kul.bib", 46},
		{"examples/scopus_old.bib", 35},
		{"examples/scopus_recent.bib", 3},
		{"examples/scopus_tidy.bib", 35},
		{"examples/scopus.bib", 20},
		{"examples/ua.bib", 37},
	}

	for _, f := range files {
		t.Run(f.File, func(t *testing.T) {
			entries := parseAll(t, f.File)

			if len(entries) != f.NumEntries {
				t.Errorf("got %d entries, want %d", len(entries), f.NumEntries)
			}
			for i, e := range entries {
				if e.Type == "" {
					t.Errorf("entry %d (line %d) has no type: %.60q", i, e.Line, e.Raw)
				}
			}
		})
	}
}

// TestForgivesWeirdKeys checks the two things "real" parsers most often choke
// on, against actual export data: junk before the first entry and unusual keys.
func TestForgivesWeirdKeys(t *testing.T) {
	tests := []struct {
		name string
		file string
		keys []string
	}{
		{
			// scopus_old.bib has no junk prefix but plenty of non-ASCII keys.
			name: "non-ASCII keys",
			file: "examples/scopus_old.bib",
			keys: []string{"Blåsjö2022587", "Cè2022411", "Ramírez-Giraldo2022199"},
		},
		{
			// ua.bib glues entries together with no separator and uses keys
			// full of colons and underscores.
			name: "colon/underscore keys in glued entries",
			file: "examples/ua.bib",
			keys: []string{"moto:c:irua:174579_de_is", "moto:c:irua:174580_vers_can"},
		},
		{
			// scopus_tidy.bib (text header) and the BOM-prefixed files prove
			// junk before the first @ is skipped rather than rejected.
			name: "junk/BOM prefix before first entry",
			file: "examples/scopus_tidy.bib",
			keys: []string{"Karamanou2022168"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := map[string]bool{}
			for _, e := range parseAll(t, tt.file) {
				seen[e.Key] = true
			}
			for _, k := range tt.keys {
				if !seen[k] {
					t.Errorf("key %q not parsed from %s", k, tt.file)
				}
			}
		})
	}
}
