// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// writeTOC writes the table of contents file.
func writeTOC(c *Config, toc []Asset) error {
	fd, err := os.Create(filepath.Join(c.Output, "bindata.go"))
	if err != nil {
		return err
	}

	defer fd.Close()

	err = writeTOCHeader(fd, c)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeTOCAsset(fd, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return writeTOCFooter(fd, c)
}

// writeTOCHeader writes the table of contents file header.
func writeTOCHeader(w io.Writer, c *Config) error {
	_, err := fmt.Fprintf(w, `package %s

// Asset loads and returns the asset for the given name.
// This returns nil of the asset could not be found.
func Asset(name string) []byte {
	if f, ok := _bindata[name]; ok {
		return f()
	}
	return nil
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string] func() []byte {
`, c.Package)
	return err
}

// writeTOCAsset write a TOC entry for the given asset.
func writeTOCAsset(w io.Writer, c *Config, asset *Asset) error {
	_, err := fmt.Fprintf(w, "\t%q: %s,\n", asset.Name, asset.Func)
	return err
}

// writeTOCFooter writes the table of contents file footer.
func writeTOCFooter(w io.Writer, c *Config) error {
	_, err := fmt.Fprintf(w, `
}
`)
	return err
}
