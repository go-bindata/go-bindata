// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"io"
)

// input -> gowriter -> output.
func translate_nomemcpy_uncomp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

import (
	"reflect"
	"unsafe"
)

var _%s = "`, pkgname, funcname)

	io.Copy(&StringWriter{Writer: output}, input)

	fmt.Fprintf(output, `"

// %s returns raw file data.
//
// WARNING: The returned byte slice is READ-ONLY.
// Attempting to alter the slice contents will yield a runtime panic.
func %s() []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&_%s))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(_%s)
	bx.Cap = bx.Len
	return b
}
`, funcname, funcname, funcname, funcname)
}
