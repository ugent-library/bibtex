package bibtex

// TODO method to split an author/editor name into its parts (von, last, first)

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
