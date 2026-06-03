package bibtex

import "testing"

func TestSplitName(t *testing.T) {
	tests := []struct {
		in   string
		want Name
	}{
		// "First von Last"
		{"John Smith", Name{First: "John", Last: "Smith"}},
		{"John Paul Smith", Name{First: "John Paul", Last: "Smith"}},
		{"Smith", Name{Last: "Smith"}},
		{"Ludwig van Beethoven", Name{First: "Ludwig", Von: "van", Last: "Beethoven"}},
		{
			"Charles Louis Xavier Joseph de la Vallee Poussin",
			Name{First: "Charles Louis Xavier Joseph", Von: "de la", Last: "Vallee Poussin"},
		},
		// "von Last, First"
		{"van Beethoven, Ludwig", Name{Von: "van", Last: "Beethoven", First: "Ludwig"}},
		{"Smith, John", Name{Last: "Smith", First: "John"}},
		// "von Last, Jr, First"
		{"von Last, Jr, First", Name{Von: "von", Last: "Last", Jr: "Jr", First: "First"}},
		{"Doe, Sr, John", Name{Last: "Doe", Jr: "Sr", First: "John"}},
		// braces protect separators and case
		{"{Barnes and Noble, Inc.}", Name{Last: "{Barnes and Noble, Inc.}"}},
		{"{de la} Cruz, Maria", Name{Last: "{de la} Cruz", First: "Maria"}},
		// an accent macro in a von-eligible (non-final) position must not be
		// read as von via its macro letter (\v -> Š, uppercase)
		{`Jean \v{S}imon Vidal`, Name{First: `Jean \v{S}imon`, Last: "Vidal"}},
	}

	for _, tt := range tests {
		if got := SplitName(tt.in); got != tt.want {
			t.Errorf("SplitName(%q) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}
