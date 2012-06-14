// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"compress/gzip"
	"fmt"
	"io"
)

// Translate the input file with optional GZIP compression.
// input [-> gzip] -> gowriter -> output.
func translate(input io.Reader, output io.Writer, pkgname, funcname string, uncompressed bool) {
	fmt.Fprintf(output, `package %s

import (`, pkgname)

	if uncompressed {
		fmt.Fprint(output, `
		"reflect"
		"unsafe"`)
	} else {
		fmt.Fprint(output, `
		"bytes"
		"compress/gzip"
		"io"
		"reflect"
		"unsafe"`)
	}

	fmt.Fprintf(output, `
)

var _%s = "`, funcname)

	if uncompressed {
		io.Copy(&GoWriter{Writer: output}, input)
	} else {
		gz := gzip.NewWriter(&GoWriter{Writer: output})
		io.Copy(gz, input)
		gz.Close()
	}

	fmt.Fprintf(output, `"

// %s returns the binary data for a given file.
func %s() []byte {`, funcname, funcname)

	if uncompressed {
		fmt.Fprintf(output, `
		// This bit of black magic ensures we do not get
		// unneccesary memcpy's and can read directly from
		// the .rodata section.
		var empty [0]byte
		sx := (*reflect.StringHeader)(unsafe.Pointer(&_%s))
		b := empty[:]
		bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		bx.Data = sx.Data
		bx.Len = len(_%s)
		bx.Cap = bx.Len
		return b`, funcname, funcname)
	} else {
		fmt.Fprintf(output, `
		// This bit of black magic ensures we do not get
		// unneccesary memcpy's and can read directly from
		// the .rodata section.
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

		return buf.Bytes()`, funcname, funcname)
	}

	fmt.Fprintf(output, "\n}")
}
