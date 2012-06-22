// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"compress/gzip"
	"fmt"
	"io"
)

// input -> gzip -> gowriter -> output.
func translate_memcpy_comp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

import (
	"bytes"
	"compress/gzip"
	"io"
)

// %s returns the raw, uncompressed file data data.
func %s() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{`, pkgname, funcname, funcname)

	gz := gzip.NewWriter(&ByteWriter{Writer: output})
	io.Copy(gz, input)
	gz.Close()

	fmt.Fprint(output, `
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}`)
}
