// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/go-bindata"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cfg, status := parseArgs()
	err := bindata.Translate(cfg, status)

	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata: %v", err)
	}
}

// parseArgs create s a new, filled configuration instance
// by reading and parsing command line options.
//
// This function exits the program with an error, if
// any of the command line options are incorrect.
func parseArgs() (*bindata.Config, bindata.ProgressFunc) {
	var version, quiet bool
	var tagstr string

	c := new(bindata.Config)

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input> [<output>]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&tagstr, "tags", "", "Comma-separated list of build tags to include.")
	flag.StringVar(&c.Prefix, "prefix", "", "Optional path prefix to strip off map keys and function names.")
	flag.BoolVar(&c.NoMemCopy, "nomemcopy", false, "Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	flag.BoolVar(&c.Compress, "compress", false, "Assets will be GZIP compressed when this flag is specified.")
	flag.BoolVar(&version, "version", false, "Displays version information.")
	flag.BoolVar(&quiet, "quiet", false, "Do not print conversion status.")
	flag.Parse()

	if version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	// Make sure we have in/output paths.
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Missing asset directory.\n")
		os.Exit(1)
	}

	// Test validity of input path.
	c.Input = filepath.Clean(flag.Arg(0))

	stat, err := os.Lstat(c.Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Input path: %v.\n", err)
		os.Exit(1)
	}

	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Input path is not a directory.\n")
		os.Exit(1)
	}

	// Find and test validity of output path.
	if flag.NArg() > 1 {
		c.Output = filepath.Clean(flag.Arg(1))

		stat, err := os.Lstat(c.Output)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Output path: %v.\n", err)
				os.Exit(1)
			}

			// Create output directory
			err = os.MkdirAll(c.Output, 0744)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Create Output directory: %v.\n", err)
				os.Exit(1)
			}

		} else if !stat.IsDir() {
			fmt.Fprintf(os.Stderr, "Output path is not a directory.\n")
			os.Exit(1)
		}
	} else {
		// If no output path is specified, use the current directory.
		c.Output, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to determine current working directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Process build tags.
	if len(tagstr) > 0 {
		c.Tags = strings.Split(tagstr, ",")
	} else {
		c.Tags = append(c.Tags, "debug")
	}

	if quiet {
		return c, nil
	}

	return c, func(file string, current, total int) bool {
		fmt.Printf("[%d/%d] %s\n", current, total, file)
		return false
	}
}
