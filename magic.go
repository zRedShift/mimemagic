package mimemagic

import (
	"bytes"
)

var utf16beBOM, utf16leBOM, utf8BOM = []byte{0xfe, 0xff}, []byte{0xff, 0xfe}, []byte{0xef, 0xbb, 0xbf}

type magic struct {
	mediaType int
	matchers  []*magicMatch
}

type magicMatch struct {
	start, length int
	pattern, mask []byte
	next          []*magicMatch
}

// MatchMagic determines the MIME type of the file in byte slice
// form. For an io.Reader wrapper see MatchReader (blank filename).
func MatchMagic(data []byte) MediaType {
	return mediaTypes[matchMagic(data)]
}

func isTextFile(data []byte) bool {
	if len(data) > 128 {
		data = data[:128]
	}
	if bytes.HasPrefix(data, utf16beBOM) || bytes.HasPrefix(data, utf16leBOM) || bytes.HasPrefix(data, utf8BOM) {
		return true
	}
	for _, b := range data {
		if b < ' ' && b != '\n' && b != '\r' && b != '\t' {
			return false
		}
	}
	return true
}

func matchMagic(data []byte) int {
	if len(data) == 0 {
		return emptyDocument
	}
	for _, m := range magicSignatures {
		if m.match(data) {
			return m.mediaType
		}
	}
	if isTextFile(data) {
		return plainText
	}
	return unknownType
}

func (m *magic) match(data []byte) bool {
	for _, mm := range m.matchers {
		if mm.match(data) {
			return true
		}
	}
	return false
}

func (m *magicMatch) match(data []byte) bool {
	if m.search(data) {
		if m.next == nil {
			return true
		}
		for _, mm := range m.next {
			if mm.match(data) {
				return true
			}
		}
	}
	return false
}

func (m *magicMatch) search(data []byte) bool {
	dataLen := len(data)
	patternLen := len(m.pattern)
	if dataLen < m.start+patternLen {
		return false
	}
	if m.mask == nil {
		if m.length == 0 {
			return bytes.Equal(data[m.start:m.start+patternLen], m.pattern)
		}
		return bytes.Contains(data[m.start:min(m.start+m.length+patternLen, dataLen)], m.pattern)
	}
	searchLen := min(m.start+m.length, dataLen-patternLen)
outer:
	for i := m.start; i <= searchLen; i++ {
		for k := 0; k < patternLen; k++ {
			if m.pattern[k] != data[i+k]&m.mask[k] {
				continue outer
			}
		}
		return true
	}
	return false
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
