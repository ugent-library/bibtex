package bibtex

import (
	"os"
	"reflect"
	"testing"
)

// TestWeird parses testdata/weird.bib, a single stream that bundles every
// awkward construct this parser is meant to forgive, and checks the exact
// sequence of entries it produces. Each property below documents one thing
// that grammar-based parsers typically reject.
//
// The file deliberately combines, in order:
//   - a UTF-8 BOM and an export header before the first entry (junk prefix)
//   - @string / @comment / @preamble returned as entries
//   - an uppercase @ARTICLE type
//   - a key containing colons and underscores
//   - a non-ASCII key (Blåsjö2022)
//   - two entries glued together with no separator (}@)
//   - a bare numeric value, an @ inside a value, an @string macro reference,
//     and string concatenation with #
//   - doubled {{...}} and nested {{x} y} braces
//   - an empty {} value
//   - a truncated final entry with no closing brace
func TestWeird(t *testing.T) {
	want := []Entry{
		{
			Type: "string", Key: "pub", Line: 4,
			Fields: []Field{{Name: "pub", RawValue: "Acme Press", Value: "Acme Press"}},
		},
		{Type: "comment", Line: 6},
		{Type: "preamble", Line: 8},
		{
			Type: "article", Key: "key:with:colons_1", Line: 10,
			Fields: []Field{
				{Name: "title", RawValue: "The {DNA} of Things", Value: "The {DNA} of Things"},
				{Name: "author", RawValue: "Doe, Jane and Smith, John", Value: "Doe, Jane and Smith, John"},
				{Name: "year", RawValue: "2020", Value: "2020"},
				{Name: "email", RawValue: "jane@example.com", Value: "jane@example.com"},
				{Name: "publisher", RawValue: "Acme Press", Value: "Acme Press"},                  // macro resolved
				{Name: "note", RawValue: "part one and part two", Value: "part one and part two"}, // # concat
			},
			RawAuthors: []string{"Doe, Jane", "Smith, John"},
			Authors:    []string{"Doe, Jane", "Smith, John"},
		},
		{
			Type: "article", Key: "Blåsjö2022", Line: 17,
			Fields: []Field{
				{Name: "title", RawValue: "{Plurale tantum}", Value: "{Plurale tantum}"},
				{Name: "editor", RawValue: "{Acme Corporation}", Value: "{Acme Corporation}"},
				{Name: "pages", RawValue: "{10--20}", Value: "{10--20}"},
				{Name: "year", RawValue: "", Value: ""},
			},
			RawEditors: []string{"{Acme Corporation}"},
			Editors:    []string{"{Acme Corporation}"},
		},
		{
			Type: "misc", Key: "truncated", Line: 24,
			Fields: []Field{
				{Name: "title", RawValue: "No closing brace here", Value: "No closing brace here"},
			},
		},
	}

	got := parseAll(t, "testdata/weird.bib")

	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d", len(got), len(want))
	}
	for i := range want {
		// Ignore Raw: it is the exact source slice and not the point of these checks.
		got[i].Raw = ""
		if !reflect.DeepEqual(*got[i], want[i]) {
			t.Errorf("entry %d:\n got  %+v\n want %+v", i, *got[i], want[i])
		}
	}
}

// parseAll reads every entry from a file, failing the test on any parse error.
func parseAll(t *testing.T, path string) []*Entry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	p := NewParser(f)
	var entries []*Entry
	for {
		e, err := p.Next()
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", path, err)
		}
		if e == nil {
			break
		}
		entries = append(entries, e)
	}
	return entries
}
