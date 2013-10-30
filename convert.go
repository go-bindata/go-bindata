// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output file specified
// in the given configuration.
func Translate(c *Config) error {
	var toc []Asset

	// Ensure our configuration has sane values.
	err := c.validate()
	if err != nil {
		return err
	}

	// Locate all the assets.
	err = findFiles(c.Input, c.Prefix, &toc)
	if err != nil {
		return err
	}

	// Create output file.
	fd, err := os.Create(c.Output)
	if err != nil {
		return err
	}

	defer fd.Close()

	// Write build tags, if applicable.
	if len(c.Tags) > 0 {
		_, err = fmt.Fprintf(fd, "// +build %s\n\n", c.Tags)
		if err != nil {
			return err
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(fd, "package %s\n\n", c.Package)
	if err != nil {
		return err
	}

	// Write assets.
	if c.Debug {
		err = writeDebug(fd, toc)
	} else {
		err = writeRelease(fd, c, toc)
	}

	if err != nil {
		return err
	}

	// Write table of contents
	return writeTOC(fd, toc)
}

// fillTOC recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(dir, prefix string, toc *[]Asset) error {
	if len(prefix) > 0 {
		dir, _ = filepath.Abs(dir)
		prefix, _ = filepath.Abs(prefix)
	}

	fd, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer fd.Close()

	list, err := fd.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range list {
		var asset Asset
		asset.Path = filepath.Join(dir, file.Name())
		asset.Name = asset.Path

		if file.IsDir() {
			findFiles(asset.Path, prefix, toc)
			continue
		}

		if strings.HasPrefix(asset.Name, prefix) {
			asset.Name = asset.Name[len(prefix):]
		}

		// If we have a leading slash, get rid of it.
		if len(asset.Name) > 0 && asset.Name[0] == '/' {
			asset.Name = asset.Name[1:]
		}

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		asset.Func = safeFunctionName(asset.Name)
		asset.Path, _ = filepath.Abs(asset.Path)
		*toc = append(*toc, asset)
	}

	return nil
}

var regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// safeFunctionName converts the given name into a name
// which qualifies as a valid function identifier.
func safeFunctionName(name string) string {
	name = strings.ToLower(name)
	name = regFuncName.ReplaceAllString(name, "_")

	// Get rid of "__" instances for niceness.
	for strings.Index(name, "__") > -1 {
		name = strings.Replace(name, "__", "_", -1)
	}

	// Leading underscores are silly (unless they prefix a digit (see below)).
	for len(name) > 1 && name[0] == '_' {
		name = name[1:]
	}

	// Identifier can't start with a digit.
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}

	return name
}
