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
	// reName         = regexp.MustCompile(`[a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+`)
	reAtName = regexp.MustCompile(`@([a-zA-Z0-9\!\$\&\*\+\-\.\/\:\;\<\>\?\[\]\^\_\` + "`" + `\|\']+)`)
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

		e.Type = strings.ToUpper(eStr[m[0]+1 : m[1]])

		// read rest of entry (matches braces)
		startPos := m[0] - 1
		// count braces
		braceLevel := strings.Count(eStr, "{") - strings.Count(eStr, "}")

		for braceLevel != 0 {
			// pos := m[1]
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

		return e, nil
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

type Entry struct {
	Type string
	Pre  string
	Raw  string
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
		log.Printf("%+v", e)
	}
}
