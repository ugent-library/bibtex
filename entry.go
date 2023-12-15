package bibtex

type Entry struct {
	Raw        string
	Type       string
	Key        string
	Fields     []Field
	RawAuthors []string
	Authors    []string
	RawEditors []string
	Editors    []string
}

type Field struct {
	Name     string
	RawValue string
	Value    string
}
