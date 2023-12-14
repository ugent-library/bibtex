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

	@string{eng="English"}

	@book{01HAKX6YXA1C4ZNYK54G1HMXEF,
	abstract     = {{The so-called \textyen{}Mamluk\textdollar sultans \downarrow\ who ruled Egypt and Syria between the late thirteenth and early sixteenth centuries AD have often been portrayed as lacking in legitimacy due to their background as slave soldiers. Sultanic biographies written by chancery officials in the early period of the sultanate have been read as part of an effort of these sultans to legitimise their position on the throne. This book reconsiders the main corpus of six such biographies written by the historians Ibn ʿAbd al-Ẓāhir (d. 1293) and his nephew Shāfiʿ ibn ʿAlī (d. 1330) and argues that these were in fact far more complex texts. An understanding of their discourses of legitimisation needs to be embedded within a broader understanding of the multi-directional discourses operating across the texts. The study proposes to interpret these texts as "spectacles", in which authors emplotted the reign of a sultan in thoroughly literary and rhetorical fashion, making especially extensive use of textual forms prevalent in the chancery. In doing so the authors reimagined the format of the biography as a performative vehicle for displaying their literary credentials and helping them negotiate positions in the chancery and the wider courtly orbit.}},
	author       = {{Van Den Bossche, Gowaart}},
	isbn         = {{9783110752243}},
	keywords     = {{Historiography,Islamic history,Mamluk,sultanate}},
	language     = eng,
	pages        = {{231}},
	publisher    = {{De Gruyter}},
	title        = {{Literary spectacles of sultanship : historiography, the chancery, and social practice in late medieval Egypt}},
	url          = {{http://doi.org/10.1515/9783110753028}},
	volume       = {{10}},
	year         = {{2023}},
	}

	@comment{
		bla bla bla
	}

	Scopus
EXPORT DATE: 07 September 2023
@article{  Van       Haute  [2]    ,
	title        = {Author Correction: Prediction of essential oil content in spearmint (Mentha spicata) via near-infrared hyperspectral imaging and chemometrics (Scientific Reports, (2023), 13, 1, (4261), 10.1038/s41598-023-31517-8)},
	author       = {Van Haute, Sam and Nikkhah, Amin and Malavi, Derick and Kiani, Sajad},
	year         = 2023,
	journal      = {Scientific Reports},
	publisher    = {Nature Research},
	volume       = 13,
	number       = 1,
	doi          = {10.1038/s41598-023-36583-6},
	issn         = 20452322,
	url          = {https://www.scopus.com/inward/record.uri?eid=2-s2.0-85161672817\&doi=10.1038\%2fs41598-023-36583-6\&partnerID=40\&md5=9e8487bfdcfa36111b9255ac3e6c38e1},
	note         = {Cited by: 0; All Open Access, Gold Open Access, Green Open Access},
	affiliations = {Department of Food Technology, Safety and Health, Faculty of Bioscience Engineering, Ghent University, Coupure Links 653, Ghent, 9000, Belgium; Department of Molecular Biotechnology, Environmental Technology, and Food Technology, Ghent University Global Campus, 119, Songdomunhwa-Ro, Yeonsu-Gu, Incheon, 21985, South Korea; Friedman School of Nutrition Science and Policy, Tufts University, Boston, MA, United States; Biosystems Engineering Department, Sari Agricultural Sciences and Natural Resources University, Sari, Iran},
	abstract     = {The Acknowledgments section in the original version of this Article was incomplete. ` + "``" + `The authors gratefully acknowledge the financial support provided by Ghent University Global Campus.'' now reads: ` + "``" + `The authors gratefully acknowledge the financial support provided by Ghent University Global Campus and Sari Agricultural Sciences and Natural Resources University (grant number: 02-1400-02).'' The original Article has been corrected. \textcopyright{} 2023, The Author(s).},
	keywords     = {erratum},
	correspondence_address = {S. Van Haute; Department of Food Technology, Safety and Health, Faculty of Bioscience Engineering, Ghent University, Ghent, Coupure Links 653, 9000, Belgium; email: sam.vanhaute@ghent.ac.kr; S. Kiani; Biosystems Engineering Department, Sari Agricultural Sciences and Natural Resources University, Sari, Iran; email: s.kiani\@sanru.ac.ir},
	pmid         = 37296160,
	language     = {English},
	abbrev_source_title = {Sci. Rep.},
	type         = {Erratum},
	publication_stage = {Final},
	source       = {Scopus}
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
		if e.Type == "COMMENT" || e.Type == "PREAMBLE" {
			log.Printf("raw: %s", e.Raw)
		} else {
			log.Printf("key: %s", e.Key)
			for _, f := range e.Fields {
				log.Printf("field %s: %s", f.Name, f.DecodeValue())
			}
			log.Printf("author: %s", strings.Join(e.Author(), ";"))
			log.Printf("editor: %s", strings.Join(e.Editor(), ";"))
		}
	}
}
