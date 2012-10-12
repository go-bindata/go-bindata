// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"io"
)

// input -> gzip -> gowriter -> output.
func translate_memcpy_uncomp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

// %s returns raw file data.
func %s() []byte {
	return []byte{`, pkgname, funcname, funcname)

	io.Copy(&ByteWriter{Writer: output}, input)

	fmt.Fprint(output, `
	}
}`)
}
