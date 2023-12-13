package bibtex

import "github.com/ugent-library/bibtex/latex"

type Entry struct {
	Pre    string
	Raw    string
	Type   string
	Key    string
	Fields []Field
}

type Field struct {
	Name  string
	Value string
}

func (f Field) DecodeValue() string {
	return latex.Decode(f.Value)
}
