[![Go Reference](https://pkg.go.dev/badge/github.com/ugent-library/bibtex.svg)](https://pkg.go.dev/github.com/ugent-library/bibtex)

# ugent-library/bibtex

Robust Golang BibTeX parser

## Examples

```go
p := bibtex.NewParser(os.Stdin)
for {
    e, err := p.Next()
    if err != nil {
        log.Fatal(err)
    }
    if e == nil {
        break
    }
    
    fmt.Printf("bibtex entry: %s", e.Raw)
    fmt.Printf("type: %s", e.Type)
    fmt.Printf("key: %s", e.Key)
    for _, f := range e.Fields {
        fmt.Printf("field %s: %s", f.Name, f.Value)
    }
    for _, a := range e.Authors {
        fmt.Printf("author: %s", a)
    }
    for _, a := range e.Editors {
        fmt.Printf("editor: %s", a)
    }
}
```