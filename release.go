// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// writeReleaseHeader writes output file headers.
// This targets release builds.
func writeReleaseHeader(w io.Writer, c *Config) {
	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		fmt.Fprintf(w, "// +build release %s\n\n", c.Tags)
	} else {
		fmt.Fprintf(w, "// +build release\n\n")
	}

	// Write package declaration
	fmt.Fprintf(w, "package %s\n\n", c.Package)

	if c.Compress {
		if c.NoMemCopy {
			header_compressed_nomemcopy(w)
		} else {
			header_compressed_memcopy(w)
		}
	} else {
		if c.NoMemCopy {
			header_uncompressed_nomemcopy(w)
		} else {
			header_uncompressed_memcopy(w)
		}
	}
}

// writeRelease write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeRelease(w io.Writer, c *Config, asset *Asset) error {
	fd, err := os.Open(asset.Path)
	if err != nil {
		return err
	}

	defer fd.Close()

	if c.Compress {
		if c.NoMemCopy {
			compressed_nomemcopy(w, asset, fd)
		} else {
			compressed_memcopy(w, asset, fd)
		}
	} else {
		if c.NoMemCopy {
			uncompressed_nomemcopy(w, asset, fd)
		} else {
			uncompressed_memcopy(w, asset, fd)
		}
	}

	return nil
}

func header_compressed_nomemcopy(w io.Writer) {
	fmt.Fprintf(w, `
import (
    "bytes"
    "compress/gzip"
    "io"
    "log"
    "reflect"
    "unsafe"
)

func bindata_read(data, name string) []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len

	gz, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		log.Fatalf("Read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		log.Fatalf("Read %%q: %%v", name, err)
	}

	return buf.Bytes()
}

`)
}

func header_compressed_memcopy(w io.Writer) {
	fmt.Fprintf(w, `
import (
    "bytes"
    "compress/gzip"
    "io"
    "log"
)

func bindata_read(data []byte, name string) []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		log.Fatalf("Read %%q: %%v", name, err)
	}

	return buf.Bytes()
}

`)
}

func header_uncompressed_nomemcopy(w io.Writer) {
	fmt.Fprintf(w, `
import (
    "reflect"
    "unsafe"
)

func bindata_read(data, name string) []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b
}

`)
}

func header_uncompressed_memcopy(w io.Writer) {
	// nop -- We require no imports or helper functions.
}

func compressed_nomemcopy(w io.Writer, asset *Asset, r io.Reader) {
	fmt.Fprintf(w, `var _%s = "`, asset.Func)

	gz := gzip.NewWriter(&StringWriter{Writer: w})
	io.Copy(gz, r)
	gz.Close()

	fmt.Fprintf(w, `"

func %s() []byte {
	return bindata_read(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
}

func compressed_memcopy(w io.Writer, asset *Asset, r io.Reader) {
	fmt.Fprintf(w, `func %s() []byte {
	return bindata_read([]byte{`, asset.Func)

	gz := gzip.NewWriter(&ByteWriter{Writer: w})
	io.Copy(gz, r)
	gz.Close()

	fmt.Fprintf(w, `
		},
		%q,
	)
}

`, asset.Name)
}

func uncompressed_nomemcopy(w io.Writer, asset *Asset, r io.Reader) {
	fmt.Fprintf(w, `var _%s = "`, asset.Func)

	io.Copy(&StringWriter{Writer: w}, r)

	fmt.Fprintf(w, `"

func %s() []byte {
	return bindata_read(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name)
}

func uncompressed_memcopy(w io.Writer, asset *Asset, r io.Reader) {
	fmt.Fprintf(w, `func %s() []byte {
	return []byte{`, asset.Func)

	io.Copy(&ByteWriter{Writer: w}, r)

	fmt.Fprintf(w, `
	}
}

`)
}
