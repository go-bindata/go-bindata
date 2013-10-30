// +build !release

package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

// bindata_read reads the given file from disk.
// It panics if anything went wrong.
func bindata_read(path, name string) []byte {
	fd, err := os.Open(path)
	if err != nil {
		log.Fatalf("Read %s: %v", name, err)
	}

	defer fd.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, fd)
	if err != nil {
		log.Fatalf("Read %s: %v", name, err)
	}

	return buf.Bytes()
}

func testdata_in_b_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/b/test.asset",
		"../testdata/in/b/test.asset",
	)
}

func testdata_in_a_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/a/test.asset",
		"../testdata/in/a/test.asset",
	)
}

func testdata_in_c_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/c/test.asset",
		"../testdata/in/c/test.asset",
	)
}

