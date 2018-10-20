package mimemagic

import (
	"os"
	"path/filepath"
	"strings"
)

type treeMagic struct {
	mediaType int
	matchers  []treeMatch
}

type treeMatch struct {
	path                            string
	mediaType                       int
	objectType                      objectType
	matchCase, executable, nonEmpty bool
	next                            []treeMatch
}

type objectType int

const (
	anyType objectType = iota
	fileType
	directoryType
	linkType
)

// MatchTreeMagic determines if the path or the directory
// of the file supplied in the path matches any common mounted
// volume signatures and returns their x-content MIME type.
// Return inode/directory MediaType in the case of a negative
// identification for a directory, and application/octet-stream
// in the case of a file.
func MatchTreeMagic(path string) (MediaType, error) {
	m, err := matchTreeMagic(path)
	return mediaTypes[m], err
}

func matchTreeMagic(path string) (int, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return unknownType, err
	}
	dir := path
	isDir := info.IsDir()
	uType := unknownType
	if !isDir {
		dir = filepath.Dir(dir)
	} else {
		uType = unknownDirectory
	}
	contents, lowercase := make(map[string]os.FileInfo), make(map[string]os.FileInfo)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path, err = filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if path == "." || path == "" {
			return nil
		}
		contents[path], lowercase[strings.ToLower(path)] = info, info
		return nil
	})
	if err != nil {
		return uType, err
	}
	for _, t := range treeMagicSignatures {
		if t.match(contents, lowercase) {
			return t.mediaType, nil
		}
	}
	return uType, nil
}

func (t treeMagic) match(contents, lowercase map[string]os.FileInfo) bool {
	for _, tt := range t.matchers {
		if tt.match(contents, lowercase) {
			return true
		}
	}
	return false
}

func (t treeMatch) match(contents, lowercase map[string]os.FileInfo) bool {
	path := t.path
	var f os.FileInfo
	var ok bool
	if !t.matchCase {
		path = strings.ToLower(path)
		f, ok = lowercase[path]
	} else {
		f, ok = contents[path]
	}
	if !ok {
		return false
	}
	switch {
	case t.objectType == fileType && !f.Mode().IsRegular(),
		t.objectType == linkType && f.Mode()&os.ModeSymlink == 0,
		t.objectType == directoryType && !f.Mode().IsDir():
		return false
	}
	if t.executable && f.Mode()&0111 == 0 {
		return false
	}
	if t.mediaType > -1 {
		if matchGlob(f.Name()) != t.mediaType {
			return false
		}
	}
	if t.nonEmpty {
		if t.objectType == fileType && f.Size() == 0 {
			return false
		} else if t.objectType == directoryType {
			m := contents
			if !t.matchCase {
				m = lowercase
			}
			for ff := range m {
				if rel, err := filepath.Rel(path, ff); err == nil && rel[0] != '.' && rel != "" {
					goto next
				}
			}
			return false
		}
	}
next:
	if t.next == nil {
		return true
	}
	for _, tt := range t.next {
		if !tt.match(contents, lowercase) {
			return false
		}
	}
	return true
}
