// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"compress/gzip"
	"fmt"
	"io"
)

// Translate the input file.
// input -> gzip -> gowriter -> output.
func translate(input io.Reader, output io.Writer, pkgname, funcname string) (err error) {
	fmt.Fprintf(output, `package %s

import (
	"bytes"
	"compress/gzip"
	"io"
)

func %s() ([]byte, error) {
	var gz *gzip.Decompressor
	var err error
	if gz, err = gzip.NewReader(bytes.NewBuffer([]byte{`, pkgname, funcname)

	gz := gzip.NewWriter(&GoWriter{Writer: output})
	io.Copy(gz, input)
	gz.Close()

	fmt.Fprint(output, `
	})); err != nil {
		return nil, err
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes(), nil
}`)
	return
}
