package mimemagic

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
)

type namespace struct {
	namespaceURI, localName string
	mediaType               int
}

// MatchXMLReader is an io.Reader wrapper for MatchXML that
// can be supplied with a limit on the data to read.
func MatchXMLReader(r io.Reader, limit int) MediaType {
	if limit < 0 || limit > 1024 {
		limit = 1024
	}
	return mediaTypes[matchXML(io.LimitReader(r, int64(limit)))]
}

// MatchXML determines the MIME type of the xml file in a byte
// slice form. Returns application/octet-stream in case the
// file isn't a valid xml and application/xml if the
// identification comes back negative.
func MatchXML(data []byte) MediaType {
	if len(data) > 1024 {
		data = data[:1024]
	}
	return mediaTypes[matchXML(bytes.NewReader(data))]
}

func matchXML(r io.Reader) int {
	uType := unknownType
	dec := xml.NewDecoder(r)
	dec.Strict = false
	dec.CharsetReader = charset.NewReaderLabel
	for {
		t, err := dec.Token()
		if err != nil {
			break
		}
		switch t := t.(type) {
		case xml.ProcInst, xml.Directive, xml.Comment:
			uType = unknownXML
		case xml.StartElement:
			uType = unknownXML
			var m int
			if m = isLocalName(t.Name.Local); m < 0 {
				continue
			}
			for _, attr := range t.Attr {
				if attr.Name.Local == "xmlns" {
					if m := isNameSpace(attr.Value); m > -1 {
						return m
					}
				}
			}
			return m
		}
	}
	return uType
}

func isLocalName(name string) int {
	for _, n := range namespaces {
		if n.localName == name {
			return n.mediaType
		}
	}
	return -1
}

func isNameSpace(name string) int {
	for _, n := range namespaces {
		if n.namespaceURI == name {
			return n.mediaType
		}
	}
	return -1
}
