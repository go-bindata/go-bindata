// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"path"
	"fmt"
	"exec"
	"io"
	"os"
	"compress/gzip"
)

// If gofmt exists on the system, run it over the target file to 
// fix up the generated code. This is not necessary, just a convenience.
func gofmt(file string) (err os.Error) {
	var prog string
	if prog = os.Getenv("GOBIN"); len(prog) == 0 {
		return
	}

	prog = path.Join(prog, "gofmt")
	cmd := exec.Command(prog, "-w", file)
	return cmd.Run()
}

// Translate the input file.
// input -> gzip -> gowriter -> output.
func translate(input io.Reader, output io.Writer, pkgname, funcname string) (err os.Error) {
	var gz *gzip.Compressor

	fmt.Fprintf(output, `package %s
import ( "io"; "os"; "bytes"; "compress/gzip" )

func %s() ([]byte, os.Error) {
var gz *gzip.Decompressor
var err os.Error
if gz, err = gzip.NewReader(bytes.NewBuffer([]byte{`, pkgname, funcname)

	if gz, err = gzip.NewWriter(&GoWriter{Writer: output}); err != nil {
		return
	}

	io.Copy(gz, input)
	gz.Close()

	fmt.Fprint(output, `
})); err != nil {
	return nil, err
}

var b bytes.Buffer
io.Copy(&b, gz)
gz.Close()
return b.Bytes(), nil}`)
	return
}
