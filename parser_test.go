package bibtex

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	shouldParse := []struct {
		File       string
		NumEntries int
	}{
		{"examples/biblatex.bib", 92},
		{"examples/biblio.bib", 23},
		{"examples/kul.bib", 46},
		{"examples/scopus_old.bib", 35},
		{"examples/scopus_recent.bib", 3},
		{"examples/scopus_tidy.bib", 35},
		{"examples/scopus.bib", 20},
		{"examples/ua.bib", 37},
	}

	for _, data := range shouldParse {
		f, err := os.Open(data.File)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		p := NewParser(f)

		var entries []*Entry

		for {
			e, err := p.Next()
			if err != nil {
				t.Error(err)
			}
			if e == nil {
				break
			}
			entries = append(entries, e)
		}

		if len(entries) != data.NumEntries {
			t.Errorf("%q: expected %d entries, got %d", data.File, data.NumEntries, len(entries))
		}
	}
}
