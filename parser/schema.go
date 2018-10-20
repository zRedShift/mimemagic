package main

import (
	"fmt"
	"math"
	"strings"
)

type mimeInfo struct {
	XMLName  struct{}    `xml:"http://www.freedesktop.org/standards/shared-mime-info mime-info"`
	MIMEType []*mimeType `xml:"mime-type"`
}

type mimeType struct {
	Type            string           `xml:"type,attr"`
	Comment         []*comment       `xml:"comment"`
	Acronym         *acronym         `xml:"acronym,omitempty"`
	ExpandedAcronym *expandedAcronym `xml:"expanded-acronym,omitempty"`
	Icon            *icon            `xml:"icon,omitempty"`
	GenericIcon     *genericIcon     `xml:"generic-icon,omitempty"`
	Glob            []*glob          `xml:"glob,omitempty"`
	Magic           []*magic         `xml:"magic,omitempty"`
	TreeMagic       []*treeMagic     `xml:"treemagic,omitempty"`
	RootXML         []*rootXML       `xml:"root-XML,omitempty"`
	Alias           []*alias         `xml:"alias,omitempty"`
	SubClassOf      []*subClassOf    `xml:"sub-class-of,omitempty"`
}

type comment struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

type acronym struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}
type expandedAcronym struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

type icon struct {
	Name string `xml:"name,attr"`
}

type genericIcon struct {
	Name string `xml:"name,attr"`
}

type glob struct {
	Pattern       string `xml:"pattern,attr"`
	Weight        *int   `xml:"weight,attr,omitempty"`
	CaseSensitive bool   `xml:"case-sensitive,attr,omitempty"`
}

type magic struct {
	Match    []*match `xml:"match"`
	Priority *int     `xml:"priority,attr,omitempty"`
}

type match struct {
	Match  []*match `xml:"match,omitempty"`
	Offset string   `xml:"offset,attr"`
	Type   string   `xml:"type,attr"`
	Value  string   `xml:"value,attr"`
	Mask   string   `xml:"mask,attr,omitempty"`
}

type treeMagic struct {
	TreeMatch []*treeMatch `xml:"treematch"`
	Priority  *int         `xml:"priority,attr,omitempty"`
}

type treeMatch struct {
	TreeMatch  []*treeMatch `xml:"treematch,omitempty"`
	Path       string       `xml:"path,attr"`
	Type       string       `xml:"type,attr,omitempty"`
	MatchCase  bool         `xml:"match-case,attr,omitempty"`
	Executable bool         `xml:"executable,attr,omitempty"`
	NonEmpty   bool         `xml:"non-empty,attr,omitempty"`
	MIMEType   string       `xml:"mimetype,attr,omitempty"`
}

type rootXML struct {
	NamespaceURI string `xml:"namespaceURI,attr"`
	LocalName    string `xml:"localName,attr"`
}

type alias struct {
	Type string `xml:"type,attr"`
}
type subClassOf struct {
	Type string `xml:"type,attr"`
}

type parsedMIMEInfo []*parsedMIMEType

type parsedMIMEType struct {
	Media, Subtype, Comment, Acronym, ExpandedAcronym, Icon, GenericIcon string
	Alias, SubClassOf, Extension                                         []string
	Glob                                                                 []*parsedGlob
	Magic                                                                []*parsedMagic
	TreeMagic                                                            []*parsedTreeMagic
	RootXML                                                              []*parsedRootXML
	Lexicographic                                                        int
}

const nilString = "nil"

func (p *parsedMIMEType) String() string {
	alias, subclass, ext, subint := nilString, nilString, nilString, nilString
	if len(p.Alias) > 0 {
		alias = fmt.Sprintf("%#v", p.Alias)
		s := make([]int, 0, len(alias))
		for _, a := range p.Alias {
			if _, ok := types[a]; !ok {
				continue
			}
			n := types[a].Lexicographic
			unique := true
			for _, ss := range s {
				if ss == n {
					unique = false
					break
				}
			}
			if unique {
				s = append(s, n)
			}
		}
	}
	if len(p.SubClassOf) > 0 {
		subclass = fmt.Sprintf("%#v", p.SubClassOf)
		s := make([]int, 0, len(p.SubClassOf))
		for _, sb := range p.SubClassOf {
			if _, ok := types[sb]; !ok {
				if sb, ok = aliases[sb]; !ok {
					continue
				}
			}
			n := types[sb].Lexicographic
			unique := true
			for _, ss := range s {
				if ss == n {
					unique = false
					break
				}
			}
			if unique {
				s = append(s, n)
			}
		}
		if len(s) > 0 {
			subint = fmt.Sprintf("%#v", s)
		}
	}
	if len(p.Extension) > 0 {
		ext = fmt.Sprintf("%#v", p.Extension)
	}
	return fmt.Sprintf("{%q, %q, %q, %q, %q, %q, %q, %s, %s, %s, %s}", p.Media, p.Subtype, p.Comment, p.Acronym, p.ExpandedAcronym, p.Icon, p.GenericIcon, alias, subclass, ext, subint)
}

type parsedGlob struct {
	Pattern       string
	Weight        int
	CaseSensitive bool
}

type parsedMagic struct {
	Priority int
	MIMEType int
	Match    []*parsedMatch
}

func (p *parsedMagic) MaxLen() int {
	max := 0
	for _, pp := range p.Match {
		if nmax := pp.MaxLen(); nmax > max {
			max = nmax
		}
	}
	return max
}

func (p *parsedMagic) TestNum() int {
	t := 0
	for _, pp := range p.Match {
		t += pp.TestNum()
	}
	return t
}

func (p *parsedMagic) MinPatternLen() int {
	min := 0
	for _, pp := range p.Match {
		if min == 0 || min > pp.MinPatternLen() {
			min = pp.MinPatternLen()
		}
	}
	return min
}

func (p *parsedMagic) String() string {
	s := make([]string, 0, len(p.Match))
	for _, pp := range p.Match {
		s = append(s, pp.String())
	}
	pMatch := fmt.Sprintf("[]*magicMatch{%s}", strings.Join(s, ", "))
	return fmt.Sprintf("{%d, %s}", p.MIMEType, pMatch)
}

type magicSliceType []*parsedMagic

func (p magicSliceType) Len() int      { return len(p) }
func (p magicSliceType) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p magicSliceType) Less(j, i int) bool {
	switch {
	case p[i].Priority < p[j].Priority:
		return true
	case p[i].Priority > p[j].Priority:
		return false
	case p[i].MinPatternLen() < p[j].MinPatternLen():
		return true
	case p[i].MinPatternLen() > p[j].MinPatternLen():
		return false
	default:
		return p[i].TestNum() < p[j].TestNum()
	}
}

type parsedMatch struct {
	RangeStart, RangeLength int
	Data, Mask              []byte
	Match                   []*parsedMatch
}

func (p *parsedMatch) MaxLen() int {
	max := p.RangeStart + p.RangeLength + len(p.Data)
	for _, pp := range p.Match {
		if nmax := pp.MaxLen(); nmax > max {
			max = nmax
		}
	}
	return max
}

func (p *parsedMatch) TestNum() int {
	t := 1
	for _, pp := range p.Match {
		t += pp.TestNum()
	}
	return t
}

func (p *parsedMatch) MinPatternLen() int {
	t := len(p.Data)
	min := 0
	for _, pp := range p.Match {
		if min == 0 || min > pp.MinPatternLen() {
			min = pp.MinPatternLen()
		}
	}
	return t + min
}

func (p *parsedMatch) String() string {
	pMatch := nilString
	if len(p.Match) > 0 {
		s := make([]string, 0, len(p.Match))
		for _, pp := range p.Match {
			s = append(s, pp.String())
		}
		pMatch = fmt.Sprintf("[]*magicMatch{%s}", strings.Join(s, ", "))
	}
	pMask := nilString
	if len(p.Mask) > 0 {
		pMask = fmt.Sprintf("%#v", p.Mask)
	}
	return fmt.Sprintf("{%d, %d, %#v, %s, %s}", p.RangeStart, p.RangeLength, p.Data, pMask, pMatch)
}

type parsedTreeMagic struct {
	Priority, MIMEType int
	TreeMatch          []*parsedTreeMatch
}

func (p *parsedTreeMagic) String() string {
	s := make([]string, 0, len(p.TreeMatch))
	for _, pp := range p.TreeMatch {
		s = append(s, pp.String())
	}
	pMatch := fmt.Sprintf("[]treeMatch{%s}", strings.Join(s, ", "))
	return fmt.Sprintf("{%d, %s}", p.MIMEType, pMatch)
}

func (p *parsedTreeMagic) TestNum() int {
	t := 0
	for _, pp := range p.TreeMatch {
		t += pp.TestNum()
	}
	return t
}

type treeMagicSliceType []*parsedTreeMagic

func (p treeMagicSliceType) Len() int      { return len(p) }
func (p treeMagicSliceType) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p treeMagicSliceType) Less(j, i int) bool {
	switch {
	case p[i].Priority < p[j].Priority:
		return true
	case p[i].Priority > p[j].Priority:
		return false
	default:
		return p[i].TestNum() < p[j].TestNum()
	}
}

type parsedTreeMatch struct {
	Path, MIMEType                  string
	Type                            int
	MatchCase, Executable, NonEmpty bool
	TreeMatch                       []*parsedTreeMatch
}

func (p *parsedTreeMatch) TestNum() int {
	t := 1
	for _, pp := range p.TreeMatch {
		t += pp.TestNum()
	}
	return t
}

func (p *parsedTreeMatch) String() string {
	pMatch := nilString
	if len(p.TreeMatch) > 0 {
		s := make([]string, 0, len(p.TreeMatch))
		for _, pp := range p.TreeMatch {
			s = append(s, pp.String())
		}
		pMatch = fmt.Sprintf("[]treeMatch{%s}", strings.Join(s, ", "))
	}
	mType := -1
	if p, ok := types[p.MIMEType]; ok {
		mType = p.Lexicographic
	}
	var pType string
	switch p.Type {
	case 0:
		pType = "anyType"
	case 1:
		pType = "fileType"
	case 2:
		pType = "directoryType"
	case 3:
		pType = "linkType"
	}
	return fmt.Sprintf("{%q, %d, %s, %t, %t, %t, %s}", p.Path, mType, pType, p.MatchCase, p.Executable, p.NonEmpty, pMatch)
}

type parsedRootXML struct {
	NamespaceURI, LocalName string
	MIMEType                int
}

func (p *parsedRootXML) String() string {
	return fmt.Sprintf("{%q, %q, %d}", p.NamespaceURI, p.LocalName, p.MIMEType)
}

type rootXMLSliceType []*parsedRootXML

func (p rootXMLSliceType) Len() int      { return len(p) }
func (p rootXMLSliceType) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p rootXMLSliceType) Less(i, j int) bool {
	return p[i].NamespaceURI < p[j].NamespaceURI
}

type identifierSliceType []identifier

func (p identifierSliceType) Len() int { return len(p) }

func (p identifierSliceType) Less(j, i int) bool {
	switch {
	case p[i].weight() < p[j].weight():
		return true
	case p[i].weight() > p[j].weight():
		return false
	case p[i].length() < p[j].length():
		return true
	case p[i].length() > p[j].length():
		return false
	case !p[i].isUpperCase() && p[j].isUpperCase():
		return true
	case p[i].isUpperCase() && !p[j].isUpperCase():
		return false
	default:
		return p[i].mimeType() > p[j].mimeType()
	}
}

func (p identifierSliceType) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p identifierSliceType) GenerateMaps() (maxLen int) {
	for _, id := range p {
		if l := id.length(); l > maxLen {
			maxLen = l
		}
		switch id := id.(type) {
		case simpleSuffix:
			if id.CaseSensitive {
				csSuffixMap[string(id.suffix)] = append(csSuffixMap[string(id.suffix)], weightedMIME{id.W, id.MIMEType})
			} else {
				suffixMap[string(id.suffix)] = append(suffixMap[string(id.suffix)], weightedMIME{id.W, id.MIMEType})
			}
		case simplePrefix:
			if id.CaseSensitive {
				csPrefixMap[string(id.prefix)] = append(csPrefixMap[string(id.prefix)], weightedMIME{id.W, id.MIMEType})
			} else {
				prefixMap[string(id.prefix)] = append(prefixMap[string(id.prefix)], weightedMIME{id.W, id.MIMEType})
			}
		case simpleText:
			if id.CaseSensitive {
				csTextMap[string(id.text)] = append(csTextMap[string(id.text)], weightedMIME{id.W, id.MIMEType})
			} else {
				textMap[string(id.text)] = append(textMap[string(id.text)], weightedMIME{id.W, id.MIMEType})
			}
		default:
			patternSlice = append(patternSlice, id)
		}
	}
	return
}

type weightedMIME struct {
	weight, mimeType int
}

func (i weightedMIME) GoString() string {
	return fmt.Sprintf("{%d, %d}", i.weight, i.mimeType)
}

type weightedMIMESlice []weightedMIME

func (i weightedMIMESlice) GoString() string {
	return fmt.Sprintf("%#v", []weightedMIME(i))[len("[]main.weightedMIME"):]
}

type prefix string
type suffix string
type text string

type simplePrefix struct {
	prefix
	CaseSensitive bool
	MIMEType, W   int
}
type simpleSuffix struct {
	suffix
	CaseSensitive bool
	MIMEType, W   int
}
type simpleText struct {
	text
	CaseSensitive bool
	MIMEType, W   int
}

type anyOf []matcher

type byteRange struct {
	Start, End byte
}

type list string

type pattern struct {
	Pattern                       []matcher
	CaseSensitive, Prefix, Suffix bool
	MIMEType, W                   int
}

func (p prefix) String() string {
	return fmt.Sprintf("Prefix(%q)", string(p))
}

func (p suffix) String() string {
	return fmt.Sprintf("Suffix(%q)", string(p))
}

func (p text) String() string {
	return fmt.Sprintf("value(%q)", string(p))
}

func (p list) String() string {
	return fmt.Sprintf("list(%q)", string(p))
}

func (p byteRange) String() string {
	return fmt.Sprintf("byteRange{%q, %q}", p.Start, p.End)
}

func (p simplePrefix) String() string {
	return fmt.Sprintf("prefix{%q, %t, %d}", string(p.prefix), p.CaseSensitive, p.MIMEType)
}

func (p simpleSuffix) String() string {
	return fmt.Sprintf("suffix{%q, %t, %d}", string(p.suffix), p.CaseSensitive, p.MIMEType)
}

func (p simpleText) String() string {
	return fmt.Sprintf("text{%q, %t, %d}", string(p.text), p.CaseSensitive, p.MIMEType)

}

func (p pattern) String() string {
	s := make([]string, len(p.Pattern))
	n := 0
	for i := range p.Pattern {
		s[i] = p.Pattern[i].String()
		n += p.Pattern[i].length()
	}
	t := "textPattern"
	if p.Suffix {
		t = "suffixPattern"
	} else if p.Prefix {
		t = "prefixPattern"
	}
	return fmt.Sprintf("%s{pattern{[]matcher{%s}, %d}, %t, %d, %d}", t, strings.Join(s, ", "), n, p.CaseSensitive, p.MIMEType, p.W)
}

func (p anyOf) String() string {
	s := make([]string, len(p))
	for i := range p {
		s[i] = p[i].String()
	}
	return fmt.Sprintf("any{%s}", strings.Join(s, ", "))
}

type matcher interface {
	length() int
	String() string
	toText() string
}

type identifier interface {
	matcher
	weight() int
	isUpperCase() bool
	mimeType() int
}

func (p simplePrefix) weight() int { return p.W }
func (p simpleText) weight() int   { return p.W }
func (p simpleSuffix) weight() int { return p.W }
func (p pattern) weight() int      { return p.W }

func (p simplePrefix) mimeType() int { return p.MIMEType }
func (p simpleText) mimeType() int   { return p.MIMEType }
func (p simpleSuffix) mimeType() int { return p.MIMEType }
func (p pattern) mimeType() int      { return p.MIMEType }

func (p simplePrefix) isUpperCase() bool { return p.CaseSensitive }
func (p simpleText) isUpperCase() bool   { return p.CaseSensitive }
func (p simpleSuffix) isUpperCase() bool { return p.CaseSensitive }
func (p pattern) isUpperCase() bool      { return p.CaseSensitive }

func (p prefix) length() int  { return len(p) }
func (p text) length() int    { return len(p) }
func (p suffix) length() int  { return len(p) }
func (list) length() int      { return 1 }
func (byteRange) length() int { return 1 }
func (anyOf) length() int     { return 1 }
func (p pattern) length() int {
	n := 0
	for _, p := range p.Pattern {
		n += p.length()
	}
	return n
}

func (p prefix) toText() string { return strings.ToLower(string(p)) }
func (p text) toText() string   { return strings.ToLower(string(p)) }
func (p suffix) toText() string { return strings.ToLower(string(p)) }
func (p list) toText() string {
	var r rune = math.MaxInt32
	for _, c := range string(p) {
		if c < r {
			r = c
		}
	}
	return strings.ToLower(fmt.Sprintf("%c", r))
}
func (p byteRange) toText() string {
	return strings.ToLower(fmt.Sprintf("%c", p.Start))
}
func (p anyOf) toText() string {
	s := p[0].toText()
	for _, m := range p {
		if m.toText() < s {
			s = m.toText()
		}
	}
	return s
}
func (p pattern) toText() string {
	str := ""
	for _, m := range p.Pattern {
		str += m.toText()
	}
	return str
}
