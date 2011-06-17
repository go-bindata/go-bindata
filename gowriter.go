// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"fmt"
	"io"
	"os"
)

type GoWriter struct {
	io.Writer
	c int
}

func (this *GoWriter) Write(p []byte) (n int, err os.Error) {
	if len(p) == 0 {
		return
	}

	for n = range p {
		if this.c%12 == 0 {
			this.Writer.Write([]byte{'\n'})
			this.c = 0
		}

		fmt.Fprintf(this.Writer, "0x%02x,", p[n])
		this.c++
	}

	return
}
