// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

// writeDebug writes the debug code file.
func writeDebug(w io.Writer, toc []Asset) error {
	err := writeDebugHeader(w)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeDebugAsset(w, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeDebugHeader writes output file headers.
// This targets debug builds.
func writeDebugHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
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
		log.Fatalf("Read %%s: %%v", name, err)
	}

	defer fd.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, fd)
	if err != nil {
		log.Fatalf("Read %%s: %%v", name, err)
	}

	return buf.Bytes()
}

`)
	return err
}

// writeDebugAsset write a debug entry for the given asset.
// A debug entry is simply a function which reads the asset from
// the original file (e.g.: from disk).
func writeDebugAsset(w io.Writer, asset *Asset) error {
	_, err := fmt.Fprintf(w, `
// %s reads file data from disk.
// It panics if something went wrong in the process.
func %s() []byte {
	return bindata_read(
		%q,
		%q,
	)
}
`, asset.Func, asset.Func, asset.Path, asset.Name)
	return err
}
