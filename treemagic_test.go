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
	{"/root", "inode/directory", true},
	{"/non/existent", "application/octet-stream", true},
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
			fpath := test.path
			if fpath[0] != '/' {
				fpath = filepath.Join(path, test.path)
			}
			got, err := MatchTreeMagic(fpath)
			if (err != nil) != test.wantErr {
				t.Errorf("MatchTreeMagic() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got.MediaType() != test.want {
				t.Errorf("MatchTreeMagic() = %v, want %v", got.MediaType(), test.want)
			}
		})
	}
	preserve := treeMagicSignatures
	t.Run("special cases", func(t *testing.T) {
		treeMagicSignatures = []treeMagic{{0, []treeMatch{{
			"mpegav", -1, fileType,
			false, false, false, nil,
		}}}}
		testfunc := func(fpath, want string) (err error) {
			fpath = filepath.Join(path, fpath)
			got, err := MatchTreeMagic(fpath)
			if err != nil {
				t.Errorf("MatchTreeMagic() error = %v", err)
				return
			}
			if got.MediaType() != want {
				t.Errorf("MatchTreeMagic() = %v, want %v", got.MediaType(), want)
			}
			return nil
		}
		testfunc("video-vcd", "inode/directory")
		treeMagicSignatures[0].matchers[0].path = "mpegav/AVSEQ01.DAT"
		treeMagicSignatures[0].matchers[0].executable = true
		testfunc("video-vcd", "inode/directory")
		treeMagicSignatures[0].matchers[0].path = "mpegav"
		treeMagicSignatures[0].matchers[0].mediaType = 0
		treeMagicSignatures[0].matchers[0].objectType = directoryType
		treeMagicSignatures[0].matchers[0].executable = false
		testfunc("video-vcd", "inode/directory")
		treeMagicSignatures[0].matchers[0].mediaType = -1
		treeMagicSignatures[0].matchers[0].nonEmpty = true
		treeMagicSignatures[0].matchers[0].path = "dir"
		testfunc(".", "inode/directory")
		treeMagicSignatures[0].matchers[0].path = "core"
		treeMagicSignatures[0].matchers[0].objectType = fileType
		testfunc(".", "inode/directory")
		treeMagicSignatures[0].matchers[0].nonEmpty = false
		treeMagicSignatures[0].matchers[0].next = make([]treeMatch, 1)
		treeMagicSignatures[0].matchers[0].next[0] = treeMagicSignatures[0].matchers[0]
		treeMagicSignatures[0].matchers[0].next[0].next = nil
		testfunc(".", "all/all")
		treeMagicSignatures[0].matchers[0].next[0].nonEmpty = true
		testfunc(".", "inode/directory")
	})
	treeMagicSignatures = preserve
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
