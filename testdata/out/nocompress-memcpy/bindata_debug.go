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

func in_b_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/b/test.asset",
		"in/b/test.asset",
	)
}

func in_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/test.asset",
		"in/test.asset",
	)
}

func in_a_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/a/test.asset",
		"in/a/test.asset",
	)
}

func in_c_test_asset() []byte {
	return bindata_read(
		"/a/code/go/src/github.com/jteeuwen/go-bindata/testdata/in/c/test.asset",
		"in/c/test.asset",
	)
}
