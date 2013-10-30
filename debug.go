// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
)

// writeDebugHeader writes output file headers with the given build tags.
// This targets debug builds.
func writeDebugHeader(w io.Writer, c *Config) {
	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		fmt.Fprintf(w, "// +build !release %s\n\n", c.Tags)
	} else {
		fmt.Fprintf(w, "// +build !release\n\n")
	}

	// Write package declaration
	fmt.Fprintf(w, "package %s\n\n", c.Package)

	// Define packages we need to import.
	// And add the asset_read function. This is called
	// from asset-specific functions.
	fmt.Fprintf(w, `import (
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
}

// writeDebug write a debug entry for the given asset.
// A debug entry is simply a function which reads the asset from
// the original file (e.g.: from disk).
func writeDebug(w io.Writer, c *Config, asset *Asset) {
	fmt.Fprintf(w, `func %s() []byte {
	return bindata_read(
		%q,
		%q,
	)
}

`, asset.Func, asset.Path, asset.Name)
}
