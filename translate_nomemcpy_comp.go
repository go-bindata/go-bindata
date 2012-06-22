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
func translate_nomemcpy_comp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

import (
	"bytes"
	"compress/gzip"
	"io"
	"reflect"
	"unsafe"
)

var _%s = "`, pkgname, funcname)

	gz := gzip.NewWriter(&StringWriter{Writer: output})
	io.Copy(gz, input)
	gz.Close()

	fmt.Fprintf(output, `"

// %s returns the raw, uncompressed file data data.
func %s() []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&_%s))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(_%s)
	bx.Cap = bx.Len

	gz, err := gzip.NewReader(bytes.NewBuffer(b))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var buf bytes.Buffer
	io.Copy(&buf, gz)
	gz.Close()

	return buf.Bytes()
}
`, funcname, funcname, funcname, funcname)
}
