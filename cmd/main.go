package main

import (
	"log"
	"strings"

	"github.com/ugent-library/bibtex"
)

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
	abstract     = {{The so-called \textyen{}Mamluk\textdollar sultans\downarrowwho ruled Egypt and Syria between the late thirteenth and early sixteenth centuries AD have often been portrayed as lacking in legitimacy due to their background as slave soldiers. Sultanic biographies written by chancery officials in the early period of the sultanate have been read as part of an effort of these sultans to legitimise their position on the throne. This book reconsiders the main corpus of six such biographies written by the historians Ibn ʿAbd al-Ẓāhir (d. 1293) and his nephew Shāfiʿ ibn ʿAlī (d. 1330) and argues that these were in fact far more complex texts. An understanding of their discourses of legitimisation needs to be embedded within a broader understanding of the multi-directional discourses operating across the texts. The study proposes to interpret these texts as "spectacles", in which authors emplotted the reign of a sultan in thoroughly literary and rhetorical fashion, making especially extensive use of textual forms prevalent in the chancery. In doing so the authors reimagined the format of the biography as a performative vehicle for displaying their literary credentials and helping them negotiate positions in the chancery and the wider courtly orbit.}},
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
	p := bibtex.NewParser(strings.NewReader(bib))
	for {
		e, err := p.Next()
		if err != nil {
			log.Fatal(err)
		}
		if e == nil {
			break
		}
		log.Print("-----------")
		log.Printf("type: %s", e.Type)
		log.Printf("key: %s", e.Key)
		for _, f := range e.Fields {
			log.Printf("field %s: %s", f.Name, f.DecodeValue())
		}
	}
}
