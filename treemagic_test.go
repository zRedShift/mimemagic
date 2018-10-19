package mimemagic

import (
	"os"
	"path/filepath"
	"testing"
)

var treeMagicTests = []struct {
	path    string
	want    string
	wantErr bool
}{
	{"image-dcf", "x-content/image-dcf", false},
	{"image-picturecd", "x-content/image-picturecd", false},
	{"software", "x-content/unix-software", false},
	{"video-bluray", "x-content/video-bluray", false},
	{"video-dvd", "x-content/video-dvd", false},
	{"video-dvd-2", "x-content/video-dvd", false},
	{"video-dvd-3", "x-content/video-dvd", false},
	{"video-hddvd", "x-content/video-hddvd", false},
	{"video-svcd", "x-content/video-svcd", false},
	{"video-vcd", "x-content/video-vcd", false},
	{"dir", "inode/directory", false},
	{"test.random", "application/octet-stream", false},
	{".", "inode/directory", false},
}

func TestMatchTreeMagic(t *testing.T) {
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
	for _, test := range treeMagicTests {
		t.Run(test.path, func(t *testing.T) {
			got, err := MatchTreeMagic(filepath.Join(path, test.path))
			if (err != nil) != test.wantErr {
				t.Errorf("MatchTreeMagic() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got.MediaType() != test.want {
				t.Errorf("MatchTreeMagic() = %v, want %v", got.MediaType(), test.want)
			}
		})
	}
}

func benchmarkMatchTreeMagic(path string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		MatchTreeMagic(path)
	}
}

func BenchmarkMatchTreeMagic(b *testing.B) {
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
	for _, dir := range treeMagicTests {
		b.Run(dir.path, func(b *testing.B) { benchmarkMatchTreeMagic(filepath.Join(path, dir.path), b) })
	}
}
