package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	reStripComment = regexp.MustCompile(`^%.*$`)
	namePattern    = `[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+`
	reAtName       = regexp.MustCompile(`@(` + namePattern + `)`)
	reKey          = regexp.MustCompile(`s*\{\s*(` + namePattern + `)\s*,[\s\n]*|\s+\r?\s*`)
	reField        = regexp.MustCompile(`[\s\n]*(` + namePattern + `)[\s\n]*=[\s\n]*`)
	reDigits       = regexp.MustCompile(`^\d+`)
	reName         = regexp.MustCompile(`^` + namePattern)
	reQuotedString = regexp.MustCompile(`^"(([^"\\]*(\\.)*[^\\"]*)*)"`)
	reConcatString = regexp.MustCompile(`^\s*#\s*`)
	reWhitespace   = regexp.MustCompile(`^\s*`)
	reEscape       = regexp.MustCompile(`^\\.`)
	reStringVal    = regexp.MustCompile(`^[^\\\{\}]+`) // TODO better name
)

type Parser struct {
	scanner *bufio.Scanner
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *Parser) Next() (*Entry, error) {
	scanner := p.scanner
	buf := strings.Builder{}

	for scanner.Scan() {
		line := reStripComment.ReplaceAllString(scanner.Text(), "")
		if line != "" {
			buf.WriteString(line + "\n")
		}
		if !strings.Contains(line, "@") {
			continue
		}

		e := &Entry{}

		eStr := buf.String()

		// get type
		// TODO do we need the match across multiple lines flag like im perl? -> (?s)
		m := reAtName.FindStringSubmatchIndex(eStr)
		if m == nil {
			// include more info (see perl)
			return nil, errors.New("type not found")
		}
		e.Type = strings.ToUpper(eStr[m[2]:m[3]])

		// read rest of entry (matches braces)
		startPos := m[0] - 1
		// count braces
		braceLevel := strings.Count(eStr, "{") - strings.Count(eStr, "}")

		for braceLevel != 0 {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			braceLevel = braceLevel + strings.Count(line, "{") - strings.Count(line, "}")
			buf.WriteString(line + "\n")
		}

		eStr = buf.String()

		// add pre to entry
		e.Pre = strings.TrimSpace(eStr[:startPos])

		// raw bibtex
		e.Raw = strings.TrimSpace(eStr[startPos:])

		// TOOD handle STRING, COMMENT, PREAMBLE

		// advance
		eStr = eStr[m[1]:]
		m = reKey.FindStringSubmatchIndex(eStr)
		if m == nil {
			// include more info (see perl)
			// TODO slurp close bracket
			return nil, errors.New("malformed entry")
		}

		e.Key = eStr[m[2]:m[3]]

		// advance
		eStr = eStr[m[1]:]
		for m = reField.FindStringSubmatchIndex(eStr); m != nil; m = reField.FindStringSubmatchIndex(eStr) {
			field := Field{Name: eStr[m[2]:m[3]]}
			// advance
			eStr = eStr[m[1]:]
			newEStr, val, err := p.parseString(eStr)
			if err != nil {
				return nil, err
			}
			eStr = newEStr
			field.Value = val
			e.Fields = append(e.Fields, field)
			// skip past next comma
			if idx := strings.Index(eStr, ","); idx > -1 {
				eStr = eStr[idx+1:]
			}
		}

		return e, nil
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *Parser) parseString(eStr string) (string, string, error) {
	str := ""

	for {
		if m := reDigits.FindStringIndex(eStr); m != nil {
			log.Print("EXTRACT DIGITS")
			str += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
		} else if m := reName.FindStringIndex(eStr); m != nil {
			log.Print("EXTRACT NAME")
			// TODO look up string in strings map
			str += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
		} else if m := reQuotedString.FindStringSubmatchIndex(eStr); m != nil {
			log.Print("EXTRACT QUOTED")
			str += eStr[m[2]:m[3]]
			eStr = eStr[m[1]:]
		} else {
			log.Print("EXTRACT BRACKETED")
			newEStr, val := p.extractBracketedValue(eStr)
			// TODO remove brackets
			str += val
			eStr = newEStr
		}

		if m := reConcatString.FindStringIndex(eStr); m != nil {
			eStr = eStr[m[1]:]
			continue
		}
		break
	}

	// TODO replace newlines? see perl

	return eStr, str, nil
}

func (p *Parser) extractBracketedValue(eStr string) (string, string) {
	val := ""

	// skip whitespace
	if m := reWhitespace.FindStringIndex(eStr); m != nil {
		eStr = eStr[m[1]:]
	}

	var depth int
	for {
		if m := reEscape.FindStringIndex(eStr); m != nil {
			val += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
			continue
		} else if strings.HasPrefix(eStr, "{") {
			val += "{"
			eStr = eStr[1:]
			depth++
			continue
		} else if strings.HasPrefix(eStr, "}") {
			val += "}"
			eStr = eStr[1:]
			depth--
			if depth > 0 {
				continue
			}
			break
		} else if m := reStringVal.FindStringIndex(eStr); m != nil {
			val += eStr[m[0]:m[1]]
			eStr = eStr[m[1]:]
			continue
		}
		break
	}

	return eStr, val
}

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

func main() {
	bib := `
	PREAMBLE

	@article{01HHFFDETR97WHGHW9JFA8X6K0,
		abstract     = {{The European Court of Human Rights (ECtHR) in a judgment of 3 October 2023 found a violation by the Turkish authorities of the right to freedom of expression via social media as guaranteed by Article 10 of the European Convention on Human Rights (ECHR). The case concerns the conviction and prison sentences of two persons, Mr. Baran Durukan and Mrs. İlknur Birol, on account of content they posted on Facebook and Twitter. Although the effects of their convictions were suspended, subject to probation periods of three and five years, the ECtHR considered the convictions and their suspension in view of their potentially chilling effect as unjustified interferences with Durukan’s and Birol’s rights under Article 10 ECHR. According to the ECtHR the interferences did not afford the requisite protection against arbitrary abuse by the public authorities of the rights guaranteed under the ECHR.}},
		author       = {{Voorhoof, Dirk}},
		issn         = {{2078-6158}},
		journal      = {{Iris, Legal newsletter of the European Audiovisual Observatory}},
		keywords     = {{freedom of expression, chilling effect, Facebook}},
		language     = {{eng}},
		number       = {{IRIS 2023-10:1/22}},
		pages        = {{2}},
		title        = {{European Court of Human Rights : Durukan and Birol v. Turkey}},
		url          = {{https://merlin.obs.coe.int/article/9895}},
		year         = {{2023}},
	  }

	@book{01HAKX6YXA1C4ZNYK54G1HMXEF,
	abstract     = {{The so-called Mamluk sultans who ruled Egypt and Syria between the late thirteenth and early sixteenth centuries AD have often been portrayed as lacking in legitimacy due to their background as slave soldiers. Sultanic biographies written by chancery officials in the early period of the sultanate have been read as part of an effort of these sultans to legitimise their position on the throne. This book reconsiders the main corpus of six such biographies written by the historians Ibn ʿAbd al-Ẓāhir (d. 1293) and his nephew Shāfiʿ ibn ʿAlī (d. 1330) and argues that these were in fact far more complex texts. An understanding of their discourses of legitimisation needs to be embedded within a broader understanding of the multi-directional discourses operating across the texts. The study proposes to interpret these texts as "spectacles", in which authors emplotted the reign of a sultan in thoroughly literary and rhetorical fashion, making especially extensive use of textual forms prevalent in the chancery. In doing so the authors reimagined the format of the biography as a performative vehicle for displaying their literary credentials and helping them negotiate positions in the chancery and the wider courtly orbit.}},
	author       = {{Van Den Bossche, Gowaart}},
	isbn         = {{9783110752243}},
	keywords     = {{Historiography,Islamic history,Mamluk,sultanate}},
	language     = {{eng}},
	pages        = {{231}},
	publisher    = {{De Gruyter}},
	title        = {{Literary spectacles of sultanship : historiography, the chancery, and social practice in late medieval Egypt}},
	url          = {{http://doi.org/10.1515/9783110753028}},
	volume       = {{10}},
	year         = {{2023}},
	}
	`
	p := NewParser(strings.NewReader(bib))
	for {
		e, err := p.Next()
		if err != nil {
			log.Fatal(err)
		}
		if e == nil {
			break
		}
		log.Printf("%+v", e.Fields)
	}
}
