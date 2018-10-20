mimemagic
=========
[![GoDoc](https://godoc.org/github.com/zRedShift/mimemagic?status.svg)](https://godoc.org/github.com/zRedShift/mimemagic)
[![Build Status](https://travis-ci.org/zRedShift/mimemagic.svg?branch=master)](https://travis-ci.org/zRedShift/mimemagic)
[![Codecov](https://codecov.io/gh/zRedShift/mimemagic/branch/master/graph/badge.svg)](https://codecov.io/gh/zRedShift/mimemagic/)
[![Go Report Card](https://goreportcard.com/badge/github.com/zRedShift/mimemagic)](https://goreportcard.com/report/github.com/zRedShift/mimemagic)

Powerful and versatile MIME sniffing package using pre-compiled glob patterns, magic number signatures, xml document
namespaces, and tree magic for mounted volumes, generated from the XDG shared-mime-info database.

## Features

- All in native go, no outside dependencies/C library bindings
- 1003 MIME types, with a description, an acronym (where available), common aliases, extensions, icons, and 
subclasses
- 493 magic signature tests (comprising of 1147 individual patterns), featuring range searches and bit masks, as per
the xdg specification
- 1099 glob patterns, for filename-based matching
- 11 Tree Magic signatures and 28 XML namespace/local name pairs, offered for completeness' sake.
- Included is the xml file parser to generate your own MIME definitions
- Also included is a CLI based on this library that is fully featured and blazing-fast, beating the native 'file'
and KDE's 'kmimetypefinder' in performance
- Cross-platform support

## Installation

The library:
```bash
go get github.com/zRedShift/mimemagic
```
The CLI:
```bash
go get github.com/zRedShift/mimemagic/cmd/mimemagic
```

## API

See the [Godoc](https://godoc.org/github.com/zRedShift/mimemagic) reference, and cmd/mimemagic for an example
implementation.

## Usage

The library:
```go
package main

import (
	"fmt"
	"github.com/zRedShift/mimemagic"
	"strings"
)

func main() {
	// Ignoring Read errors that might arise
	mimeType, _ := mimemagic.MatchFilePath("sample.svgz", -1)

	// image/svg+xml-compressed
	fmt.Println(mimeType.MediaType())

	// compressed SVG image
	fmt.Println(mimeType.Comment)

	// SVG (Scalable Vector Graphics)
	fmt.Printf("%s (%s)\n", mimeType.Acronym, mimeType.ExpandedAcronym)

	// application/gzip
	fmt.Println(strings.Join(mimeType.SubClassOf, ", "))

	// .svgz
	fmt.Println(strings.Join(mimeType.Extensions, ", "))

	// This is an image.
	switch mimeType.Media {
	case "image":
		fmt.Println("This is an image.")
	case "video":
		fmt.Println("This is a video file.")
	case "audio":
		fmt.Println("This is an audio file.")
	case "application":
		fmt.Println("This is an application.")
	default:
		fmt.Printf("This is a(n) %s.", mimeType.Media)
	}

	// true
	fmt.Println(mimeType.IsExtension(".svgz"))
}
```
The CLI:
```
Usage: mimemagic [options] <file> ...
Determines the MIME type of the given file(s).

Options:
  -c    Determine the MIME type of the file(s) using only its content.
  -f    Determine the MIME type of the file(s) using only the file name. Does
        not check for the file's existence. The -c
         flag takes precedence.
  -i    Output the MIME type in a human readable format.
  -l int
        The number of bytes from the beginning of the file mimemagic will
        examine. Reads the entire file if set to a negative value. By default
        mimemagic will only read the first 512 from stdin, however setting this
        flag to a non-default negative value will override this. (default -1)
  -t    Determine the MIME type of the directory/mounted volume using tree
        magic. Can't be used in conjunction with with -c, -f or -x.
  -x    Determine the MIME type of the xml file(s) using the local names and
        namespaces within. Can't be used in conjunction with -c, -f or -t.

Arguments:
  file
        The file(s) to test. '-' to read from stdin. If '-' is set, all other
        inputs will be ignored.

Examples:
  $ mimemagic -c sample.svgz
    	application/gzip
  $ mimemagic *.svg*
    	Olympic_rings_with_transparent_rims.svg: image/svg+xml
    	Piano.svg.png: image/png
    	RAID_5.svg: image/svg+xml
    	sample.svgz: image/svg+xml-compressed
  $ cat /dev/urandom | mimemagic -
    	application/octet-stream
  $ ls software; mimemagic -i -t software/
    	autorun
    	UNIX software
```

## Benchmarks

See [Benchmarks](https://github.com/zRedShift/mimemagic/blob/master/benchmarks.txt). For Match(), the average across 
over 400 completely different files (representing a unique MIME type each) is 13 ± 7 μs/op. For MatchGlob() it's 900
± 200 ns/op, and for 12 ± 7 μs/op MatchMagic().