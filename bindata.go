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
func translate(input, output, pkgname, funcname string) (err os.Error) {
	var fs, fd *os.File
	var gz *gzip.Compressor

	if fs, err = os.Open(input); err != nil {
		return
	}

	defer fs.Close()

	if fd, err = os.Create(output); err != nil {
		return
	}

	defer fd.Close()

	fmt.Fprintf(fd, `// auto generated from '%s'.

package %s
import ( "io"; "os"; "bytes"; "compress/gzip" )

func %s() ([]byte, os.Error) {
var gz *gzip.Decompressor
var err os.Error
if gz, err = gzip.NewReader(bytes.NewBuffer([]byte{`, input, pkgname, funcname)

	if gz, err = gzip.NewWriter(&GoWriter{Writer: fd}); err != nil {
		return
	}

	io.Copy(gz, fs)
	gz.Close()

	fmt.Fprint(fd, `
})); err != nil {
	return nil, err
}

var b bytes.Buffer
io.Copy(&b, gz)
gz.Close()
return b.Bytes(), nil}`)
	return
}
