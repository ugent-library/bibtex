package latex

import (
	"regexp"
	"strings"
)

var reNormalizeMacro1 = regexp.MustCompile(`(\\[a-zA-Z]+)\\(\s+)`) // \foo\ bar -> \foo{} bar
var reNormalizeMacro2 = regexp.MustCompile(`([^{]\\\w)([;,.:%])`)  //} Aaaa\o, -> Aaaa\o{},
var remacros *regexp.Regexp

func init() {
	macroNames := make([]string, 0, len(macros))
	for k := range macros {
		macroNames = append(macroNames, k)
	}
	remacros = regexp.MustCompile(`\\(` + strings.Join(macroNames, "|") + `)(?:\{\}|\s+|\b)`)
	// remacros = regexp.MustCompile(`\\(` + strings.Join(macroNames, "|") + `)(?:\{\})?`)
}

func Decode(str string) string {
	str = reNormalizeMacro1.ReplaceAllString(str, "$1{}$2")
	// log.Printf("str: %s", str)
	str = reNormalizeMacro2.ReplaceAllString(str, "$1{}$2")
	// log.Printf("str: %s", str)

	// TODO this doesn't always work?
	str = remacros.ReplaceAllStringFunc(str, func(macro string) string {
		m := remacros.FindStringSubmatch(macro)
		return macros[m[1]]
	})
	return str
}

var macros = map[string]string{
	"AA":                   "\u00C5",
	"aa":                   "\u00E5",
	"AE":                   "\u00C6",
	"ae":                   "\u00E6",
	"ain":                  "\u02BF",
	"angle":                "\u2220",
	"approx":               "\u2248",
	"approxeq":             "\u224A",
	"ast":                  "\u2217",
	"asymp":                "\u224D",
	"ayn":                  "\u02BF",
	"backsim":              "\u223D",
	"backsimeq":            "\u22CD",
	"barwedge":             "\u22BC",
	"because":              "\u2235",
	"between":              "\u226C",
	"bigcap":               "\u22C2",
	"bigcup":               "\u22C3",
	"bigvee":               "\u22C1",
	"bigwedge":             "\u22C0",
	"bot":                  "\u22A5",
	"bowtie":               "\u22C8",
	"Box":                  "\u25A1",
	"boxdot":               "\u22A1",
	"boxminus":             "\u229F",
	"boxplus":              "\u229E",
	"boxtimes":             "\u22A0",
	"bullet":               "\u2219",
	"Bumpeq":               "\u224E",
	"bumpeq":               "\u224F",
	"cap":                  "\u2229",
	"Cap":                  "\u22D2",
	"cdot":                 "\u22C5",
	"cdots":                "\u22EF",
	"circ":                 "\u2218",
	"circeq":               "\u2257",
	"circledast":           "\u229B",
	"circledcirc":          "\u229A",
	"circleddash":          "\u229D",
	"clubsuit":             "\u2663",
	"complement":           "\u2201",
	"cong":                 "\u2245",
	"coprod":               "\u2210",
	"copyright":            "\u00A9",
	"cup":                  "\u222A",
	"Cup":                  "\u22D3",
	"curlyeqprec":          "\u22DE",
	"curlyeqsucc":          "\u22DF",
	"curlyvee":             "\u22CE",
	"curlywedge":           "\u22CF",
	"dag":                  "\u2020",
	"dashv":                "\u22A3",
	"ddag":                 "\u2021",
	"ddots":                "\u22F1",
	"Delta":                "\u2206",
	"DH":                   "\u00D0",
	"dh":                   "\u00F0",
	"diamond":              "\u22C4",
	"diamondsuit":          "\u2662",
	"div":                  "\u00F7",
	"divideontimes":        "\u22C7",
	"DJ":                   "\u0110",
	"dj":                   "\u0111",
	"doteq":                "\u2250",
	"doteqdot":             "\u2251",
	"dotplus":              "\u2214",
	"dots":                 "\u2026",
	"downarrow":            "\u2193",
	"eqcirc":               "\u2256",
	"equiv":                "\u2261",
	"ESH":                  "\u01A9",
	"euro":                 "\u20AC",
	"exists":               "\u2203",
	"fallingdotseq":        "\u2252",
	"flat":                 "\u266D",
	"forall":               "\u2200",
	"geq":                  "\u2265",
	"geqq":                 "\u2267",
	"gg":                   "\u226B",
	"ggg":                  "\u22D9",
	"gneqq":                "\u2269",
	"gnsim":                "\u22E7",
	"gtrdot":               "\u22D7",
	"gtreqless":            "\u22DB",
	"gtrless":              "\u2277",
	"gtrsim":               "\u2273",
	"guillemotleft":        "\u00AB",
	"guillemotright":       "\u00BB",
	"guilsinglleft":        "\u2039",
	"guilsinglright":       "\u203A",
	"hamza":                "\u02BE",
	"heartsuit":            "\u2661",
	"hv":                   "\u0195",
	"i":                    "\u0131",
	"iiint":                "\u222D",
	"iint":                 "\u222C",
	"IJ":                   "\u0132",
	"ij":                   "\u0133",
	"in":                   "\u2208",
	"infty":                "\u221E",
	"int":                  "\u222B",
	"intercal":             "\u22BA",
	"L":                    "\u0141",
	"l":                    "\u0142",
	"langle":               "\u27E8",
	"lceil":                "\u2308",
	"leadsto":              "\u219D",
	"leftarrow":            "\u2190",
	"leftrightarrow":       "\u2194",
	"Leftrightarrow":       "\u21D4",
	"leftthreetimes":       "\u22CB",
	"leq":                  "\u2264",
	"leqq":                 "\u2266",
	"lessdot":              "\u22D6",
	"lesseqgtr":            "\u22DA",
	"lessgtr":              "\u2276",
	"lesssim":              "\u2272",
	"lfloor":               "\u230A",
	"lhd":                  "\u22B2",
	"ll":                   "\u226A",
	"lll":                  "\u22D8",
	"lneqq":                "\u2268",
	"lnot":                 "\u00AC",
	"lnsim":                "\u22E6",
	"ltimes":               "\u22C9",
	"measuredangle":        "\u2221",
	"mid":                  "\u2223",
	"mp":                   "\u2213",
	"mu":                   "\u00B5",
	"multimap":             "\u22B8",
	"nabla":                "\u2207",
	"natural":              "\u266E",
	"ncong":                "\u2247",
	"neq":                  "\u2260",
	"nexists":              "\u2204",
	"NG":                   "\u014A",
	"ng":                   "\u014B",
	"ngeq":                 "\u2271",
	"ngtr":                 "\u226F",
	"ni":                   "\u220B",
	"nleq":                 "\u2270",
	"nless":                "\u226E",
	"nmid":                 "\u2224",
	"nobreakspace":         "\u00A0",
	"notin":                "\u2209",
	"nparallel":            "\u2226",
	"nprec":                "\u2280",
	"nsim":                 "\u2241",
	"nsubseteq":            "\u2288",
	"nsucc":                "\u2281",
	"nsupseteq":            "\u2289",
	"ntriangleleft":        "\u22EA",
	"ntrianglelefteq":      "\u22EC",
	"ntriangleright":       "\u22EB",
	"ntrianglerighteq":     "\u22ED",
	"nVdash":               "\u22AE",
	"O":                    "\u00D8",
	"o":                    "\u00F8",
	"odot":                 "\u2299",
	"OE":                   "\u0152",
	"oe":                   "\u0153",
	"OHORN":                "\u01A0",
	"ohorn":                "\u01A1",
	"oint":                 "\u222E",
	"ominus":               "\u2296",
	"oplus":                "\u2295",
	"oslash":               "\u2298",
	"otimes":               "\u2297",
	"P":                    "\u00B6",
	"parallel":             "\u2225",
	"partial":              "\u2202",
	"pitchfork":            "\u22D4",
	"pm":                   "\u00B1",
	"pounds":               "\u00A3",
	"prec":                 "\u227A",
	"preccurlyeq":          "\u227C",
	"precnsim":             "\u22E8",
	"precsim":              "\u227E",
	"prod":                 "\u220F",
	"propto":               "\u221D",
	"quotedblbase":         "\u201E",
	"quotesinglbase":       "\u201A",
	"rangle":               "\u27E9",
	"rceil":                "\u2309",
	"rfloor":               "\u230B",
	"rhd":                  "\u22B3",
	"rightarrow":           "\u2192",
	"Rightarrow":           "\u21D2",
	"rightleftharpoons":    "\u21CC",
	"rightthreetimes":      "\u22CC",
	"risingdotseq":         "\u2253",
	"rtimes":               "\u22CA",
	"S":                    "\u00A7",
	"set":                  "\u2205",
	"setminus":             "\u2216",
	"sharp":                "\u266F",
	"sim":                  "\u223C",
	"simeq":                "\u2243",
	"spadesuit":            "\u2660",
	"sphericalangle":       "\u2222",
	"sqcap":                "\u2293",
	"sqcup":                "\u2294",
	"sqsubset":             "\u228F",
	"sqsubseteq":           "\u2291",
	"sqsupset":             "\u2290",
	"sqsupseteq":           "\u2292",
	"ss":                   "\u00DF",
	"star":                 "\u22C6",
	"subset":               "\u2282",
	"Subset":               "\u22D0",
	"subseteq":             "\u2286",
	"subsetneq":            "\u228A",
	"succ":                 "\u227B",
	"succcurlyeq":          "\u227D",
	"succnsim":             "\u22E9",
	"succsim":              "\u227F",
	"sum":                  "\u2211",
	"supset":               "\u2283",
	"Supset":               "\u22D1",
	"supseteq":             "\u2287",
	"supsetneq":            "\u228B",
	"surd":                 "\u221A",
	"textampersand":        "\u0026",
	"textasciiacute":       "\u00B4",
	"textasciicedilla":     "\u00B8",
	"textasciicircum":      "\u005E",
	"textasciidieresis":    "\u00A8",
	"textasciigrave":       "\u0060",
	"textasciimacron":      "\u00AF",
	"textasciitilde":       "\u007E",
	"textasteriskcentered": "\u002A",
	"textbackslash":        "\u005C",
	"textbar":              "\u007C",
	"textbardotlessj":      "\u025F",
	"textbarglotstop":      "\u02A1",
	"textbari":             "\u0268",
	"textbarl":             "\u0142",
	"textbaro":             "\u0275",
	"textbarrevglotstop":   "\u02A2",
	"textbaru":             "\u0289",
	"textbeltl":            "\u026C",
	"textBhook":            "\u0181",
	"textbhook":            "\u0253",
	"textbraceleft":        "\u007B",
	"textbraceright":       "\u007D",
	"textbrokenbar":        "\u00A6",
	"textbullet":           "\u2022",
	"textbullseye":         "\u0298",
	"textcent":             "\u00A2",
	"textcentereddot":      "\u00B7",
	"textChook":            "\u0187",
	"textchook":            "\u0188",
	"textcloseepsilon":     "\u029A",
	"textcloseomega":       "\u0277",
	"textcloserevepsilon":  "\u025E",
	"textcolonmonetary":    "\u20A1",
	"textcopyright":        "\u00A9",
	"textcrb":              "\u0180",
	"textcrd":              "\u0111",
	"textcrh":              "\u0127",
	"textcrlambda":         "\u019B",
	"textctc":              "\u0255",
	"textctesh":            "\u0286",
	"textctj":              "\u029D",
	"textctyogh":           "\u0293",
	"textctz":              "\u0291",
	"textcurrency":         "\u00A4",
	"textDafrican":         "\u0189",
	"textdctzlig":          "\u02A5",
	"textdegree":           "\u00B0",
	"textDhook":            "\u018A",
	"textdhook":            "\u0257",
	"textdiv":              "\u00F7",
	"textdollar":           "\u0024",
	"textdong":             "\u20AB",
	"textdtail":            "\u0256",
	"textdyoghlig":         "\u02A4",
	"textdzlig":            "\u02A3",
	"textemdash":           "\u2014",
	"textendash":           "\u2013",
	"textEopen":            "\u0190",
	"texteopen":            "\u025B",
	"textepsilon":          "\u025B",
	"textequals":           "\u003D",
	"textEreversed":        "\u018E",
	"textEsh":              "\u01A9",
	"textesh":              "\u0283",
	"texteturned":          "\u01DD",
	"texteuro":             "\u20AC",
	"textexclamdown":       "\u00A1",
	"textEzh":              "\u01B7",
	"textezh":              "\u0292",
	"textFhook":            "\u0191",
	"textfishhookr":        "\u027E",
	"textflorin":           "\u0192",
	"textg":                "\u0067",
	"textgamma":            "\u0263",
	"textGammaafrican":     "\u0194",
	"textgammalatinsmall":  "\u0263",
	"textglotstop":         "\u0294",
	"textgreater":          "\u003E",
	"texthash":             "\u0023",
	"textHbar":             "\u0126",
	"texthbar":             "\u0127",
	"texthtb":              "\u0253",
	"texthtbardotlessj":    "\u0284",
	"texthtc":              "\u0188",
	"texthtd":              "\u0257",
	"texthtg":              "\u0260",
	"texthth":              "\u0266",
	"texththeng":           "\u0267",
	"texthtk":              "\u0199",
	"texthtp":              "\u01A5",
	"texthtq":              "\u02A0",
	"texthtscg":            "\u029B",
	"texthtt":              "\u01AD",
	"texthvlig":            "\u0195",
	"textinterrobang":      "\u203D",
	"textinvglotstop":      "\u0296",
	"textinvscr":           "\u0281",
	"textiota":             "\u0269",
	"textIotaafrican":      "\u0196",
	"textiotalatin":        "\u0269",
	"textKhook":            "\u0198",
	"textkhook":            "\u0199",
	"textkra":              "\u0138",
	"textlengthmark":       "\u02D0",
	"textless":             "\u003C",
	"textlhti":             "\u027F", // ??
	"textlira":             "\u20A4",
	"textlogicalnot":       "\u00AC",
	"textlonglegr":         "\u027C",
	"textlooptoprevesh":    "\u01AA",
	"textltailm":           "\u0271",
	"textltailn":           "\u0272",
	"textltilde":           "\u026B",
	"textlyoghlig":         "\u026E",
	"textminus":            "\u2212",
	"textmu":               "\u00B5",
	"textnaira":            "\u20A6",
	"textNhookleft":        "\u019D",
	"textnhookleft":        "\u0272",
	"textnumero":           "\u2116",
	"textonehalf":          "\u00BD",
	"textonequarter":       "\u00BC",
	"textonesuperior":      "\u00B9",
	"textOopen":            "\u0186",
	"textoopen":            "\u0254",
	"textopeno":            "\u0254",
	"textordfeminine":      "\u00AA",
	"textordmasculine":     "\u00BA",
	"textoverline":         "\u203E",
	"textpalhookbelow":     "\u01AB",
	"textparagraph":        "\u00B6",
	"textpercent":          "\u0025",
	"textperiodcentered":   "\u00B7",
	"textpertenthousand":   "\u2031",
	"textperthousand":      "\u2030",
	"textphi":              "\u0278",
	"textPhook":            "\u01A4",
	"textphook":            "\u01A5",
	"textpm":               "\u00B1",
	"textprimstress":       "\u02C8",
	"textquestiondown":     "\u00BF",
	"textquotedbl":         "\u0022",
	"textquotedblleft":     "\u201C",
	"textquotedblright":    "\u201D",
	"textquoteleft":        "\u2018",
	"textquoteright":       "\u2019",
	"textquotesingle":      "\u0027",
	"textraisevibyi":       "\u0285", // ??
	"textramshorns":        "\u0264",
	"textreferencemark":    "\u203B",
	"textregistered":       "\u00AE",
	"textreve":             "\u0258",
	"textrevepsilon":       "\u025C",
	"textrevglotstop":      "\u0295",
	"textrhookrevepsilon":  "\u025D",
	"textrhookschwa":       "\u025A",
	"textrtaild":           "\u0256",
	"textrtaill":           "\u026D",
	"textrtailn":           "\u0273",
	"textrtailr":           "\u027D",
	"textrtails":           "\u0282",
	"textrtailt":           "\u0288",
	"textrtailz":           "\u0290",
	"textscb":              "\u0299",
	"textscg":              "\u0262",
	"textsch":              "\u029C",
	"textschwa":            "\u0259",
	"textsci":              "\u026A",
	"textscl":              "\u029F",
	"textscn":              "\u0274",
	"textscoelig":          "\u0276",
	"textscr":              "\u0280",
	"textscripta":          "\u0251",
	"textscriptg":          "\u0261",
	"textscriptv":          "\u028B",
	"textscy":              "\u028F",
	"textsection":          "\u00A7",
	"textsterling":         "\u00A3",
	"textstretchc":         "\u0297",
	"textTbar":             "\u0166",
	"texttbar":             "\u0167",
	"texttctclig":          "\u02A8",
	"texttesh":             "\u02A7",
	"textteshlig":          "\u02A7",
	"textThook":            "\u01AC",
	"textthook":            "\u01AD",
	"textthorn":            "\u00FE",
	"textthornvari":        "\u00FE",
	"textthornvarii":       "\u00FE",
	"textthornvariii":      "\u00FE",
	"textthornvariv":       "\u00FE",
	"textthreequarters":    "\u00BE",
	"textthreesuperior":    "\u00B3",
	"texttimes":            "\u00D7",
	"texttrademark":        "\u2122",
	"textTretroflexhook":   "\u01AE",
	"texttretroflexhook":   "\u0288",
	"texttslig":            "\u02A6",
	"textTstroke":          "\u0166",
	"texttstroke":          "\u0167",
	"textturna":            "\u0250",
	"textturnh":            "\u0265",
	"textturnk":            "\u029E",
	"textturnlonglegr":     "\u027A",
	"textturnm":            "\u026F",
	"textturnmrleg":        "\u0270",
	"textturnr":            "\u0279",
	"textturnrrtail":       "\u027B",
	"textturnscripta":      "\u0252",
	"textturnt":            "\u0287",
	"textturnv":            "\u028C",
	"textturnw":            "\u028D",
	"textturny":            "\u028E",
	"texttwosuperior":      "\u00B2",
	"textunderscore":       "\u005F",
	"textupsilon":          "\u028A",
	"textVhook":            "\u01B2",
	"textvhook":            "\u028B",
	"textwon":              "\u20A9",
	"textyen":              "\u00A5",
	"textYhook":            "\u01B3",
	"textyhook":            "\u01B4",
	"textyogh":             "\u0292",
	"TH":                   "\u00DE",
	"th":                   "\u00FE",
	"therefore":            "\u2234",
	"Thorn":                "\u00DE",
	"times":                "\u00D7",
	"tone1":                "\u02E9",
	"tone2":                "\u02E8",
	"tone3":                "\u02E7",
	"tone4":                "\u02E6",
	"tone5":                "\u02E5",
	"top":                  "\u22A4",
	"triangleq":            "\u225C",
	"UHORN":                "\u01AF",
	"uhorn":                "\u01B0",
	"unlhd":                "\u22B4",
	"unrhd":                "\u22B5",
	"uparrow":              "\u2191",
	"updownarrow":          "\u2195",
	"uplus":                "\u228E",
	"vdash":                "\u22A2",
	"Vdash":                "\u22A9",
	"vdots":                "\u22EE",
	"vee":                  "\u2228",
	"veebar":               "\u22BB",
	"Vvdash":               "\u22AA",
	"wedge":                "\u2227",
	"wr":                   "\u2240",
}
