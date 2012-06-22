// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import "io"

// translate translates the input file to go source code.
func translate(input io.Reader, output io.Writer, pkgname, funcname string, uncompressed, nomemcpy bool) {
	if nomemcpy {
		if uncompressed {
			translate_nomemcpy_uncomp(input, output, pkgname, funcname)
		} else {
			translate_nomemcpy_comp(input, output, pkgname, funcname)
		}
	} else {
		if uncompressed {
			translate_memcpy_uncomp(input, output, pkgname, funcname)
		} else {
			translate_memcpy_comp(input, output, pkgname, funcname)
		}
	}
}
