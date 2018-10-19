/*
Package mimemagic implements MIME sniffing using pre-compiled glob
patterns, magic number signatures, xml document namespaces, and tree magic for
mounted volumes, generated from the XDG shared-mime-info database.

To generate your own database simply remove the leading space, point to the
directory with freedesktop.org package files (freedesktop.org.xml, if it
exists, is always processed first and Override.xml is always processed last),
and run go generate:
  go:generate go run -ldflags "-X main.dir=/usr/share/mime/packages/" github.com/zRedShift/mimemagic/parser/

To use the default freedesktop.org.xml file provided in this package:
  go:generate go run github.com/zRedShift/mimemagic/parser/

globs.go is generated unformatted so it's a good idea to run this for your OCD
  go:generate go fmt globs.go
*/
package mimemagic
