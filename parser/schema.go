package main

import (
	"fmt"
	"math"
	"strings"
)

type MIMEInfo struct {
	XMLName  struct{}    `xml:"http://www.freedesktop.org/standards/shared-mime-info mime-info"`
	MIMEType []*MIMEType `xml:"mime-type"`
}

type MIMEType struct {
	Type            string           `xml:"type,attr"`
	Comment         []*Comment       `xml:"comment"`
	Acronym         *Acronym         `xml:"acronym,omitempty"`
	ExpandedAcronym *ExpandedAcronym `xml:"expanded-acronym,omitempty"`
	Icon            *Icon            `xml:"icon,omitempty"`
	GenericIcon     *GenericIcon     `xml:"generic-icon,omitempty"`
	Glob            []*Glob          `xml:"glob,omitempty"`
	Magic           []*Magic         `xml:"magic,omitempty"`
	TreeMagic       []*TreeMagic     `xml:"treemagic,omitempty"`
	RootXML         []*RootXML       `xml:"root-XML,omitempty"`
	Alias           []*Alias         `xml:"alias,omitempty"`
	SubClassOf      []*SubClassOf    `xml:"sub-class-of,omitempty"`
}

type Comment struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

type Acronym struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}
type ExpandedAcronym struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

type Icon struct {
	Name string `xml:"name,attr"`
}
type GenericIcon struct {
	Name string `xml:"name,attr"`
}
type Glob struct {
	Pattern       string `xml:"pattern,attr"`
	Weight        *int   `xml:"weight,attr,omitempty"`
	CaseSensitive bool   `xml:"case-sensitive,attr,omitempty"`
}

type Magic struct {
	Match    []*Match `xml:"match"`
	Priority *int     `xml:"priority,attr,omitempty"`
}

type Match struct {
	Match  []*Match `xml:"match,omitempty"`
	Offset string   `xml:"offset,attr"`
	Type   string   `xml:"type,attr"`
	Value  string   `xml:"value,attr"`
	Mask   string   `xml:"mask,attr,omitempty"`
}

type TreeMagic struct {
	TreeMatch []*TreeMatch `xml:"treematch"`
	Priority  *int         `xml:"priority,attr,omitempty"`
}

type TreeMatch struct {
	TreeMatch  []*TreeMatch `xml:"treematch,omitempty"`
	Path       string       `xml:"path,attr"`
	Type       string       `xml:"type,attr,omitempty"`
	MatchCase  bool         `xml:"match-case,attr,omitempty"`
	Executable bool         `xml:"executable,attr,omitempty"`
	NonEmpty   bool         `xml:"non-empty,attr,omitempty"`
	MIMEType   string       `xml:"mimetype,attr,omitempty"`
}

type RootXML struct {
	NamespaceURI string `xml:"namespaceURI,attr"`
	LocalName    string `xml:"localName,attr"`
}

type Alias struct {
	Type string `xml:"type,attr"`
}
type SubClassOf struct {
	Type string `xml:"type,attr"`
}

type ParsedMIMEInfo []*ParsedMIMEType

type ParsedMIMEType struct {
	Media, Subtype, Comment, Acronym, ExpandedAcronym, Icon, GenericIcon string
	Alias, SubClassOf, Extension                                         []string
	Glob                                                                 []*ParsedGlob
	Magic                                                                []*ParsedMagic
	TreeMagic                                                            []*ParsedTreeMagic
	RootXML                                                              []*ParsedRootXML
	Lexicographic                                                        int
}

func (p *ParsedMIMEType) String() string {
	alias, subclass, ext, subint := "nil", "nil", "nil", "nil"
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

type ParsedGlob struct {
	Pattern       string
	Weight        int
	CaseSensitive bool
}

type ParsedMagic struct {
	Priority int
	MIMEType int
	Match    []*ParsedMatch
}

func (p *ParsedMagic) MaxLen() int {
	max := 0
	for _, pp := range p.Match {
		if nmax := pp.MaxLen(); nmax > max {
			max = nmax
		}
	}
	return max
}

func (p *ParsedMagic) TestNum() int {
	t := 0
	for _, pp := range p.Match {
		t += pp.TestNum()
	}
	return t
}

func (p *ParsedMagic) MinPatternLen() int {
	min := 0
	for _, pp := range p.Match {
		if min == 0 || min > pp.MinPatternLen() {
			min = pp.MinPatternLen()
		}
	}
	return min
}

func (p *ParsedMagic) String() string {
	s := make([]string, 0, len(p.Match))
	for _, pp := range p.Match {
		s = append(s, fmt.Sprintf("%s", pp))
	}
	pMatch := fmt.Sprintf("[]*magicMatch{%s}", strings.Join(s, ", "))
	return fmt.Sprintf("{%d, %s}", p.MIMEType, pMatch)
}

type MagicSlice []*ParsedMagic

func (p MagicSlice) Len() int      { return len(p) }
func (p MagicSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p MagicSlice) Less(j, i int) bool {
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

type ParsedMatch struct {
	RangeStart, RangeLength int
	Data, Mask              []byte
	Match                   []*ParsedMatch
}

func (p *ParsedMatch) MaxLen() int {
	max := p.RangeStart + p.RangeLength + len(p.Data)
	for _, pp := range p.Match {
		if nmax := pp.MaxLen(); nmax > max {
			max = nmax
		}
	}
	return max
}

func (p *ParsedMatch) TestNum() int {
	t := 1
	for _, pp := range p.Match {
		t += pp.TestNum()
	}
	return t
}

func (p *ParsedMatch) MinPatternLen() int {
	t := len(p.Data)
	min := 0
	for _, pp := range p.Match {
		if min == 0 || min > pp.MinPatternLen() {
			min = pp.MinPatternLen()
		}
	}
	return t + min
}

func (p *ParsedMatch) String() string {
	pMatch := "nil"
	if len(p.Match) > 0 {
		s := make([]string, 0, len(p.Match))
		for _, pp := range p.Match {
			s = append(s, fmt.Sprintf("%s", pp))
		}
		pMatch = fmt.Sprintf("[]*magicMatch{%s}", strings.Join(s, ", "))
	}
	pMask := "nil"
	if len(p.Mask) > 0 {
		pMask = fmt.Sprintf("%#v", p.Mask)
	}
	return fmt.Sprintf("{%d, %d, %#v, %s, %s}", p.RangeStart, p.RangeLength, p.Data, pMask, pMatch)
}

type ParsedTreeMagic struct {
	Priority, MIMEType int
	TreeMatch          []*ParsedTreeMatch
}

func (p *ParsedTreeMagic) String() string {
	s := make([]string, 0, len(p.TreeMatch))
	for _, pp := range p.TreeMatch {
		s = append(s, fmt.Sprintf("%s", pp))
	}
	pMatch := fmt.Sprintf("[]treeMatch{%s}", strings.Join(s, ", "))
	return fmt.Sprintf("{%d, %s}", p.MIMEType, pMatch)
}

func (p *ParsedTreeMagic) TestNum() int {
	t := 0
	for _, pp := range p.TreeMatch {
		t += pp.TestNum()
	}
	return t
}

type TreeMagicSlice []*ParsedTreeMagic

func (p TreeMagicSlice) Len() int      { return len(p) }
func (p TreeMagicSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p TreeMagicSlice) Less(j, i int) bool {
	switch {
	case p[i].Priority < p[j].Priority:
		return true
	case p[i].Priority > p[j].Priority:
		return false
	default:
		return p[i].TestNum() < p[j].TestNum()
	}
}

type ParsedTreeMatch struct {
	Path, MIMEType                  string
	Type                            int
	MatchCase, Executable, NonEmpty bool
	TreeMatch                       []*ParsedTreeMatch
}

func (p *ParsedTreeMatch) TestNum() int {
	t := 1
	for _, pp := range p.TreeMatch {
		t += pp.TestNum()
	}
	return t
}

func (p *ParsedTreeMatch) String() string {
	pMatch := "nil"
	if len(p.TreeMatch) > 0 {
		s := make([]string, 0, len(p.TreeMatch))
		for _, pp := range p.TreeMatch {
			s = append(s, fmt.Sprintf("%s", pp))
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

type ParsedRootXML struct {
	NamespaceURI, LocalName string
	MIMEType                int
}

func (p *ParsedRootXML) String() string {
	return fmt.Sprintf("{%q, %q, %d}", p.NamespaceURI, p.LocalName, p.MIMEType)
}

type RootXMLSLice []*ParsedRootXML

func (p RootXMLSLice) Len() int      { return len(p) }
func (p RootXMLSLice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p RootXMLSLice) Less(i, j int) bool {
	return p[i].NamespaceURI < p[j].NamespaceURI
}

type IdentifierSlice []Identifier

func (p IdentifierSlice) Len() int { return len(p) }

func (p IdentifierSlice) Less(j, i int) bool {
	switch {
	case p[i].Weight() < p[j].Weight():
		return true
	case p[i].Weight() > p[j].Weight():
		return false
	case p[i].Len() < p[j].Len():
		return true
	case p[i].Len() > p[j].Len():
		return false
	case !p[i].Case() && p[j].Case():
		return true
	case p[i].Case() && !p[j].Case():
		return false
	default:
		return p[i].Type() > p[j].Type()
	}
}

func (p IdentifierSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p IdentifierSlice) GenerateMaps() (maxLen int) {
	for _, id := range p {
		if l := id.Len(); l > maxLen {
			maxLen = l
		}
		switch id := id.(type) {
		case SimpleSuffix:
			if id.CaseSensitive {
				csSuffixMap[string(id.Suffix)] = append(csSuffixMap[string(id.Suffix)], weightedMIME{id.W, id.MIMEType})
			} else {
				suffixMap[string(id.Suffix)] = append(suffixMap[string(id.Suffix)], weightedMIME{id.W, id.MIMEType})
			}
		case SimplePrefix:
			if id.CaseSensitive {
				csPrefixMap[string(id.Prefix)] = append(csPrefixMap[string(id.Prefix)], weightedMIME{id.W, id.MIMEType})
			} else {
				prefixMap[string(id.Prefix)] = append(prefixMap[string(id.Prefix)], weightedMIME{id.W, id.MIMEType})
			}
		case SimpleText:
			if id.CaseSensitive {
				csTextMap[string(id.Text)] = append(csTextMap[string(id.Text)], weightedMIME{id.W, id.MIMEType})
			} else {
				textMap[string(id.Text)] = append(textMap[string(id.Text)], weightedMIME{id.W, id.MIMEType})
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

type Prefix string
type Suffix string
type Text string

type SimplePrefix struct {
	Prefix
	CaseSensitive bool
	MIMEType, W   int
}
type SimpleSuffix struct {
	Suffix
	CaseSensitive bool
	MIMEType, W   int
}
type SimpleText struct {
	Text
	CaseSensitive bool
	MIMEType, W   int
}

type Any []Matcher

type Range struct {
	Start, End byte
}

type List string

type Pattern struct {
	Pattern                       []Matcher
	CaseSensitive, Prefix, Suffix bool
	MIMEType, W                   int
}

func (p Prefix) String() string {
	return fmt.Sprintf("Prefix(%q)", string(p))
}

func (p Suffix) String() string {
	return fmt.Sprintf("Suffix(%q)", string(p))
}

func (p Text) String() string {
	return fmt.Sprintf("value(%q)", string(p))
}

func (p List) String() string {
	return fmt.Sprintf("list(%q)", string(p))
}

func (p Range) String() string {
	return fmt.Sprintf("byteRange{%q, %q}", p.Start, p.End)
}

func (p SimplePrefix) String() string {
	return fmt.Sprintf("prefix{%q, %t, %d}", string(p.Prefix), p.CaseSensitive, p.MIMEType)
}

func (p SimpleSuffix) String() string {
	return fmt.Sprintf("suffix{%q, %t, %d}", string(p.Suffix), p.CaseSensitive, p.MIMEType)
}

func (p SimpleText) String() string {
	return fmt.Sprintf("text{%q, %t, %d}", string(p.Text), p.CaseSensitive, p.MIMEType)

}

func (p Pattern) String() string {
	s := make([]string, len(p.Pattern))
	n := 0
	for i := range p.Pattern {
		s[i] = p.Pattern[i].String()
		n += p.Pattern[i].Len()
	}
	t := "textPattern"
	if p.Suffix {
		t = "suffixPattern"
	} else if p.Prefix {
		t = "prefixPattern"
	}
	return fmt.Sprintf("%s{pattern{[]matcher{%s}, %d}, %t, %d, %d}", t, strings.Join(s, ", "), n, p.CaseSensitive, p.MIMEType, p.W)
}

func (p Any) String() string {
	s := make([]string, len(p))
	for i := range p {
		s[i] = p[i].String()
	}
	return fmt.Sprintf("any{%s}", strings.Join(s, ", "))
}

type Matcher interface {
	Len() int
	String() string
	ToText() string
}

type Identifier interface {
	Matcher
	Weight() int
	Case() bool
	Type() int
}

func (p SimplePrefix) Weight() int { return p.W }
func (p SimpleText) Weight() int   { return p.W }
func (p SimpleSuffix) Weight() int { return p.W }
func (p Pattern) Weight() int      { return p.W }

func (p SimplePrefix) Type() int { return p.MIMEType }
func (p SimpleText) Type() int   { return p.MIMEType }
func (p SimpleSuffix) Type() int { return p.MIMEType }
func (p Pattern) Type() int      { return p.MIMEType }

func (p SimplePrefix) Case() bool { return p.CaseSensitive }
func (p SimpleText) Case() bool   { return p.CaseSensitive }
func (p SimpleSuffix) Case() bool { return p.CaseSensitive }
func (p Pattern) Case() bool      { return p.CaseSensitive }

func (p Prefix) Len() int { return len(p) }
func (p Text) Len() int   { return len(p) }
func (p Suffix) Len() int { return len(p) }
func (List) Len() int     { return 1 }
func (Range) Len() int    { return 1 }
func (Any) Len() int      { return 1 }
func (p Pattern) Len() int {
	n := 0
	for _, p := range p.Pattern {
		n += p.Len()
	}
	return n
}

func (p Prefix) ToText() string { return strings.ToLower(string(p)) }
func (p Text) ToText() string   { return strings.ToLower(string(p)) }
func (p Suffix) ToText() string { return strings.ToLower(string(p)) }
func (p List) ToText() string {
	var r rune = math.MaxInt32
	for _, c := range string(p) {
		if c < r {
			r = c
		}
	}
	return strings.ToLower(fmt.Sprintf("%c", r))
}
func (p Range) ToText() string {
	return strings.ToLower(fmt.Sprintf("%c", p.Start))
}
func (p Any) ToText() string {
	s := p[0].ToText()
	for _, m := range p {
		if m.ToText() < s {
			s = m.ToText()
		}
	}
	return s
}
func (p Pattern) ToText() string {
	var str strings.Builder
	for _, m := range p.Pattern {
		str.WriteString(m.ToText())
	}
	return str.String()
}
