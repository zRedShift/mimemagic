package mimemagic

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var xmlTests = []struct {
	filename string
	want     string
}{
	{"560051.xml", "application/xml"},
	{"en_US.zip.meta4", "application/metalink4+xml"},
	{"feed2", "application/xml"},
	{"feed.atom", "application/atom+xml"},
	{"feed.rss", "application/xml"},
	{"feeds.opml", "application/xml"},
	{"googleearth.kml", "application/vnd.google-earth.kml+xml"},
	{"ISOcyr1.ent", "application/xml"},
	{"ooo-test.fodg", "application/xml"},
	{"ooo-test.fodp", "application/xml"},
	{"ooo-test.fods", "application/xml"},
	{"ooo-test.fodt", "application/xml"},
	{"pom.xml", "application/xml"},
	{"settings.xml", "application/xml"},
	{"Stallman_Richard_-_The_GNU_Manifesto.fb2", "application/x-fictionbook+xml"},
	{"test10.gpx", "application/gpx+xml"},
	{"test.gpx", "application/gpx+xml"},
	{"test.metalink", "application/metalink+xml"},
	{"test.mml", "application/mathml+xml"},
	{"test.owx", "application/owl+xml"},
	{"test.xht", "application/xhtml+xml"},
	{"test.xhtml", "application/xhtml+xml"},
	{"test.xml.in", "application/xml"},
	{"test.xsl", "application/xslt+xml"},
	{"xml-in-mp3.mp3", "application/octet-stream"},
}

func TestMatchXML(t *testing.T) {
	path, err := unpackFixtures()
	if err != nil {
		t.Fatalf("couldn't unpack archive: %v", err)
	}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}()
	for _, test := range xmlTests {
		t.Run(test.filename, func(t *testing.T) {
			data, err := ioutil.ReadFile(filepath.Join(path, test.filename))
			if err != nil {
				t.Fatalf("couldn't read file %s: %v", test.filename, err)
			}
			if got := MatchXML(data).MediaType(); got != test.want {
				t.Errorf("MatchXML() = %v, want %v", got, test.want)
			}
		})
	}
	for _, test := range xmlTests {
		t.Run(test.filename, func(t *testing.T) {
			f, err := os.Open(filepath.Join(path, test.filename))
			if err != nil {
				t.Fatalf("couldn't open file %s: %v", test.filename, err)
			}
			if got := MatchXMLReader(f, -1).MediaType(); got != test.want {
				t.Errorf("MatchXML() = %v, want %v", got, test.want)
			}
			f.Close()
		})
	}
}

func benchmarkMatchXML(filename string, b *testing.B) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatalf("couldn't read file %s: %v", filename, err)
	}
	for n := 0; n < b.N; n++ {
		MatchXML(data)
	}
}

func BenchmarkMatchXML(b *testing.B) {
	path, err := unpackFixtures()
	if err != nil {
		b.Fatalf("couldn't unpack archive: %v", err)
	}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}()
	for _, f := range xmlTests {
		b.Run(f.filename, func(b *testing.B) {
			benchmarkMatchXML(filepath.Join(path, f.filename), b)
		})
	}
}
