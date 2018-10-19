package main

import (
	"flag"
	"fmt"
	"github.com/zRedShift/mimemagic"
	"os"
	"path/filepath"
)

var (
	contentOnly     bool
	filenameOnly    bool
	xmlNamespace    bool
	treeMagic       bool
	humanReadable   bool
	prependFilename bool
	standardInput   bool
	mimeType        mimemagic.MediaType
	input           *os.File
	err             error
	limit           int
	files           []string
)

func init() {
	flag.BoolVar(&contentOnly, "c", false,
		"Determine the MIME type of the file(s) using only its content.")
	flag.BoolVar(&humanReadable, "i", false,
		"Output the MIME type in a human readable format.")
	flag.BoolVar(&filenameOnly, "f", false,
		"Determine the MIME type of the file(s) using only the file name. Does\n"+
			"not check for the file's existence. The -c\n flag takes precedence.")
	flag.IntVar(&limit, "l", -1,
		"The number of bytes from the beginning of the file mimemagic will\n"+
			"examine. Reads the entire file if set to a negative value. By default\n"+
			"mimemagic will only read the first 512 from stdin, however setting this\n"+
			"flag to a non-default negative value will override this.")
	flag.BoolVar(&treeMagic, "t", false,
		"Determine the MIME type of the directory/mounted volume using tree\n"+
			"magic. Can't be used in conjunction with with -c, -f or -x.")
	flag.BoolVar(&xmlNamespace, "x", false,
		"Determine the MIME type of the xml file(s) using the local names and\n"+
			"namespaces within. Can't be used in conjunction with -c, -f or -t.")
}

func main() {
	flag.Usage = usage
	setup()
	if standardInput {
		identify("")
		os.Exit(0)
	}
	for _, filename := range files {
		input, err = os.Open(filename)
		if printError(err) {
			continue
		}
		identify(filename)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, "Usage: mimemagic [options] <file> ...\n"+
		"Determines the MIME type of the given file(s).\n\n"+
		"Options:\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\nArguments:\n"+
		"  file\n"+
		"    \tThe file(s) to test. '-' to read from stdin. If '-' is set, all other\n"+
		"    \tinputs will be ignored.\n")
}

func setup() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprint(os.Stderr, "filename not specified\n")
		flag.Usage()
		os.Exit(2)
	}
	if (treeMagic || xmlNamespace) && (contentOnly || filenameOnly) || (treeMagic && xmlNamespace) {
		fmt.Fprint(os.Stderr, "invalid flag combination\n")
		flag.Usage()
		os.Exit(2)
	}
	files = flag.Args()
	if files[0] == "-" {
		filenameOnly = false
		treeMagic = false
		input = os.Stdin
		if limit == -1 {
			limit = 512
		}
		standardInput = true
		files = files[:1]
		return
	}
	if len(files) > 1 {
		prependFilename = true
	}
}

func printError(err error) bool {
	if err != nil {
		if pErr, ok := err.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "cannot %s '%s': %v\n", pErr.Op, pErr.Path, pErr.Err)
			err = nil
			return true
		}
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(2)
	}
	return false
}

func identify(filename string) {
	switch {
	case contentOnly:
		mimeType, err = mimemagic.MatchReader(input, "", limit)
	case filenameOnly:
		mimeType = mimemagic.MatchGlob(filepath.Base(filename))
	case treeMagic:
		mimeType, err = mimemagic.MatchTreeMagic(filename)
	case xmlNamespace:
		mimeType = mimemagic.MatchXMLReader(input, limit)
	default:
		mimeType, err = mimemagic.MatchFile(input, limit)
	}
	if !printError(err) {
		if prependFilename {
			if !humanReadable {
				fmt.Println(filepath.Base(filename) + ": " + mimeType.MediaType())
			} else {
				fmt.Println(filepath.Base(filename) + ": " + mimeType.Comment)
			}
		} else {
			if !humanReadable {
				fmt.Println(mimeType.MediaType())
			} else {
				fmt.Println(mimeType.Comment)
			}
		}
	}
	input.Close()
}
