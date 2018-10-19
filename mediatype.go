package mimemagic

import (
	"io"
	"os"
	"syscall"
)

// MediaType stores all the parsed values of a MIME type within a
// shared-mime-info package.
type MediaType struct {
	Media, Subtype, Comment, Acronym, ExpandedAcronym, Icon, GenericIcon string
	Alias, SubClassOf, Extensions                                        []string
	subClassOf                                                           []int
}

// MediaType returns the MIME type in the format of the MIME spec.
func (m MediaType) MediaType() string {
	return m.Media + "/" + m.Subtype
}

// IsExtension checks if the extension ext is associated with
// the MIME type. The extension should begin with a leading
// dot, as in ".html".
func (m MediaType) IsExtension(ext string) bool {
	for _, e := range m.Extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// MatchFilePath is a file path convenience wrapper for MatchReader.
func MatchFilePath(path string, limit int) (MediaType, error) {
	f, err := os.Open(path)
	if err != nil {
		return mediaTypes[unknownType], err
	}
	defer f.Close()
	return MatchReader(f, f.Name(), limit)
}

// MatchFile is an *os.File convenience wrapper for MatchReader.
func MatchFile(f *os.File, limit int) (MediaType, error) {
	return MatchReader(f, f.Name(), limit)
}

// MatchReader is an io.Reader wrapper for Match that can be
// supplied with a filename and a limit on the data to read.
// Negative values of limit will read the file up until the
// longest magic signature in the database.
func MatchReader(r io.Reader, filename string, limit int) (MediaType, error) {
	if limit < 0 || limit > magicMaxLen {
		limit = magicMaxLen
	}
	data := make([]byte, limit)
	//io.EOF check for zero-size files
	if n, err := io.ReadAtLeast(r, data, limit); err == io.ErrUnexpectedEOF || err == io.EOF {
		data = data[:n]
	} else if pErr, ok := err.(*os.PathError); ok {
		if pErr.Err == syscall.EISDIR {
			return mediaTypes[unknownDirectory], nil
		}
	} else if err != nil {
		return mediaTypes[unknownType], err
	}
	if filename == "" {
		return MatchMagic(data), nil
	}
	return Match(data, filename), nil
}

// Match determines the MIME type of the file in a byte slice
// form with a given filename. Anonymous buffers should use
// MatchMagic.
func Match(data []byte, filename string) MediaType {
	return mediaTypes[match(data, filename)]
}

func match(data []byte, filename string) int {
	globMatches := matchGlobAll(filename)
	if globMatches[0] == unknownType {
		return matchMagic(data)
	}
	if len(data) == 0 {
		return globMatches[0]
	}
	match := unknownType
	for _, m := range magicSignatures {
		if m.match(data) {
			if t := equalOrSuperClass(globMatches, m.mediaType); t > -1 {
				return globMatches[t]
			}
			if match == unknownType {
				match = m.mediaType
			}
		}
	}
	if match == unknownType && isTextFile(data) {
		if t := equalOrSuperClass(globMatches, plainText); t > -1 {
			return globMatches[t]
		}
		match = plainText
	}
	if match == unknownType || len(globMatches) == 1 {
		return globMatches[0]
	}
	return match
}

func equalOrSuperClass(globMatches []int, magicMatch int) int {
	for i := range globMatches {
		if magicMatch == globMatches[i] || equalOrSuperClass(mediaTypes[globMatches[i]].subClassOf, magicMatch) > -1 {
			return i
		}
	}
	return -1
}
