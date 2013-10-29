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

// ProgressFunc is a callback handler which is fired whenever
// bindata begins translating a new file.
//
// It takes the file path and the current and total file counts.
// This can be used to indicate progress of the conversion.
//
// If this handler returns true, the processing is stopped and
// Translate() returns immediately.
type ProgressFunc func(file string, current, total int) bool

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output directory specified
// in the given configuration.
func Translate(c *Config, pf ProgressFunc) error {
	toc := make(map[string]string)
	err := findFiles(c.Input, c.Prefix, toc)

	if err != nil {
		return err
	}

	var current int
	for key, value := range toc {
		if pf != nil {
			current++
			if pf(key, current, len(toc)) {
				return nil
			}
		}

		_ = value
	}

	return nil
}

// fillTOC recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(dir, prefix string, toc map[string]string) error {
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
		key := filepath.Join(dir, file.Name())

		if file.IsDir() {
			findFiles(key, prefix, toc)
		} else {
			if strings.HasPrefix(key, prefix) {
				key = key[len(prefix):]
			}

			// If we have a leading slash, get rid of it.
			if len(key) > 0 && key[0] == '/' {
				key = key[1:]
			}

			// This shouldn't happen.
			if len(key) == 0 {
				return fmt.Errorf("Invalid file: %v", filepath.Join(dir, file.Name()))
			}

			value := safeFunctionName(key)
			toc[key] = value
		}
	}

	return nil
}

var regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// safeFunctionName converts the given name into a name
// which qualifies as a valid function identifier.
func safeFunctionName(name string) string {
	name = strings.ToLower(name)
	name = regFuncName.ReplaceAllString(name, "_")

	// Identifier can't start with a digit.
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}

	// Get rid of "__" instances for niceness.
	for strings.Index(name, "__") > -1 {
		name = strings.Replace(name, "__", "_", -1)
	}

	return name
}
