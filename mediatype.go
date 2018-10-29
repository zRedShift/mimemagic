package mimemagic

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// MediaType stores all the parsed values of a MIME type within a
// shared-mime-info package.
type MediaType struct {
	Media, Subtype, Comment, Acronym, ExpandedAcronym, Icon, GenericIcon string
	Alias, SubClassOf, Extensions                                        []string
	subClassOf                                                           []int
}

const (
	// Default behaviour relies on MatchGlob if it returns a sole
	// match, else it defers to the first magic match.
	Default = iota
	// Magic prefers MatchMagic in case of a contention.
	Magic
	// Glob prefers MatchGlob in case of a contention.
	Glob
)

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
func MatchFilePath(path string, limAndPref ...int) (m MediaType, err error) {
	f, err := os.Open(path)
	if err != nil {
		return mediaTypes[unknownType], err
	}
	defer f.Close()
	return MatchFile(f, limAndPref...)
}

// MatchFile is an *os.File convenience wrapper for MatchReader.
func MatchFile(f *os.File, limAndPref ...int) (MediaType, error) {
	return MatchReader(f, filepath.Base(f.Name()), limAndPref...)
}

// MatchReader is an io.Reader wrapper for Match that can be
// supplied with a filename, a limit on the data to read and
// whether to prefer any of the matching methods in case of
// a contention.
// Negative or non-existent values of limit will read the
// file up until the longest magic signature in the database.
func MatchReader(r io.Reader, filename string, limAndPref ...int) (MediaType, error) {
	limit := magicMaxLen
	preference := Default
	if len(limAndPref) > 0 && limAndPref[0] >= 0 && limAndPref[0] < magicMaxLen {
		limit = limAndPref[0]
	}
	if len(limAndPref) > 1 && limAndPref[1] <= Glob {
		preference = limAndPref[1]
	}
	data := make([]byte, limit)
	//io.EOF check for zero-size files
	if n, err := io.ReadAtLeast(r, data, limit); err == io.ErrUnexpectedEOF || err == io.EOF {
		data = data[:n]
	} else if pErr, ok := err.(*os.PathError); ok && pErr.Err == syscall.EISDIR {
		return mediaTypes[unknownDirectory], nil
	} else if err != nil {
		return mediaTypes[unknownType], err
	}
	if filename == "" {
		return MatchMagic(data), nil
	}
	return Match(data, filename, preference), nil
}

// Match determines the MIME type of the file in a byte slice
// form with a given filename. Anonymous buffers should use
// MatchMagic.
// Preference is an optional value that allows to prioritize
// glob/magic matching in case of a contention. Contention is
// when both magic and glob matches are found, but they can't
// be reconciled via aliases or subclasses.
func Match(data []byte, filename string, preference ...int) MediaType {
	if len(preference) == 0 {
		return mediaTypes[match(data, filename, Default)]
	}
	return mediaTypes[match(data, filename, preference[0])]
}

func match(data []byte, filename string, preference int) int {
	globMatches := matchGlobAll(filename)
	if globMatches[0] == unknownType {
		return matchMagic(data)
	}
	if len(data) == 0 {
		if preference == Magic {
			return emptyDocument
		}
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
	if match == unknownType || preference == Glob || (preference != Magic && len(globMatches) == 1) {
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
