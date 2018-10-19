package mimemagic

import (
	"sort"
	"strings"
)

type byteMatcher interface {
	matchByte(byte) bool
}
type matcher interface {
	len() int
	match(string) bool
}
type glob interface {
	matcher
	isCaseSensitive() bool
	mediaType() simpleGlob
}

type value string
type list string
type byteRange struct{ min, max byte }
type any []byteMatcher
type pattern struct {
	matchers []matcher
	length   int
}
type textPattern struct {
	pattern
	caseSensitive    bool
	mimeType, weight int
}
type suffixPattern struct {
	pattern
	caseSensitive    bool
	mimeType, weight int
}
type prefixPattern struct {
	pattern
	caseSensitive    bool
	mimeType, weight int
}

func (v value) len() int   { return len(v) }
func (list) len() int      { return 1 }
func (byteRange) len() int { return 1 }
func (any) len() int       { return 1 }
func (p pattern) len() int { return p.length }

func (l list) matchByte(b byte) bool      { return strings.IndexByte(string(l), b) >= 0 }
func (r byteRange) matchByte(b byte) bool { return r.min <= b && b <= r.max }
func (a any) matchByte(b byte) bool {
	for _, m := range a {
		if m.matchByte(b) {
			return true
		}
	}
	return false
}

func (v value) match(s string) bool {
	return s == string(v)
}
func (l list) match(s string) bool {
	return l.matchByte(s[0])
}
func (r byteRange) match(s string) bool {
	return r.matchByte(s[0])
}
func (a any) match(s string) bool {
	return a.matchByte(s[0])
}
func (p pattern) match(s string) bool {
	for _, ml := range p.matchers {
		if !ml.match(s[:ml.len()]) {
			return false
		}
		s = s[ml.len():]
	}
	return true
}
func (t textPattern) match(s string) bool { return len(s) == t.len() && t.pattern.match(s) }
func (t suffixPattern) match(s string) bool {
	return len(s) >= t.len() && t.pattern.match(s[len(s)-t.len():])
}
func (t prefixPattern) match(s string) bool { return len(s) >= t.len() && t.pattern.match(s[:t.len()]) }

func (t textPattern) isCaseSensitive() bool   { return t.caseSensitive }
func (t suffixPattern) isCaseSensitive() bool { return t.caseSensitive }
func (t prefixPattern) isCaseSensitive() bool { return t.caseSensitive }

func (t textPattern) mediaType() simpleGlob   { return simpleGlob{t.weight, t.mimeType} }
func (t suffixPattern) mediaType() simpleGlob { return simpleGlob{t.weight, t.mimeType} }
func (t prefixPattern) mediaType() simpleGlob { return simpleGlob{t.weight, t.mimeType} }

type simpleGlob struct {
	weight, mimeType int
}

// MatchGlob determines the MIME type of the file using
// exclusively its filename.
func MatchGlob(filename string) MediaType {
	return mediaTypes[matchGlob(filename)]
}

func matchGlob(filename string) int {
	return matchGlobAll(filename)[0]
}

func matchGlobAll(filename string) []int {
	var globResults []simpleGlob
	lowerCase := strings.ToLower(filename)
	if t, ok := textCS[filename]; ok {
		globResults = append(globResults, t...)
	}
	if t, ok := text[lowerCase]; ok {
		globResults = append(globResults, t...)
	}
	fnLen := len(filename)
	for l := min(len(filename), globMaxLen); l > 0; l-- {
		if t, ok := suffixesCS[filename[fnLen-l:]]; ok {
			globResults = append(globResults, t...)
		}
		if t, ok := prefixesCS[filename[:l]]; ok {
			globResults = append(globResults, t...)
		}
		if t, ok := suffixes[lowerCase[fnLen-l:]]; ok {
			globResults = append(globResults, t...)
		}
		if t, ok := prefixes[lowerCase[:l]]; ok {
			globResults = append(globResults, t...)
		}
	}
	for _, g := range globs {
		if (g.isCaseSensitive() || g.match(lowerCase)) && (!g.isCaseSensitive() || g.match(filename)) {
			globResults = append(globResults, g.mediaType())
		}
	}
	if globResults == nil {
		return []int{unknownType}
	}
	sort.Slice(globResults, func(i, j int) bool { return globResults[i].weight > globResults[j].weight })
	results := make([]int, len(globResults))
	for i := range globResults {
		results[i] = globResults[i].mimeType
	}
	return results
}
