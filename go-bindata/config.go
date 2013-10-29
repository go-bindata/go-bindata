// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config defines command line options.
type Config struct {
	Input  string   // Input directory with assets.
	Output string   // Output directory for generated code.
	Tags   []string // Build tags to include in output files.
}

// NewConfig create s anew, filled configuration instance
// by reading and parsing command line options.
//
// This function exits the program with an error, if
// any of the command line options are incorrect.
func NewConfig() *Config {
	var version bool
	var tagstr string

	c := new(Config)

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input> [<output>]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&tagstr, "tags", "", "Comma-separated list of build tags to include.")
	flag.BoolVar(&version, "version", false, "Displays version information.")
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

	return c
}
