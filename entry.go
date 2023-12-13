package main

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
