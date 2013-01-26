// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"unicode"
)

const (
	AppName    = "bindata"
	AppVersion = "1.0.0"
)

var (
	pipe         = false
	in           = flag.String("i", "", "Path to the input file. Alternatively, pipe the file data into stdin.")
	out          = flag.String("o", "", "Optional path to the output file.")
	pkgname      = flag.String("p", "", "Optional name of the package to generate.")
	funcname     = flag.String("f", "", "Optional name of the function/variable to generate.")
	uncompressed = flag.Bool("u", false, "The specified resource will /not/ be GZIP compressed when this flag is specified. This alters the generated output code.")
	nomemcopy    = flag.Bool("m", false, "Use the memcopy hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	version      = flag.Bool("v", false, "Display version information.")
)

func main() {
	parseArgs()

	if pipe {
		translate(os.Stdin, os.Stdout, *pkgname, *funcname, *uncompressed, *nomemcopy)
	} else {
		fs, err := os.Open(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s\n", err)
			return
		}

		defer fs.Close()

		fd, err := os.Create(*out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s\n", err)
			return
		}

		defer fd.Close()

		translate(fs, fd, *pkgname, *funcname, *uncompressed, *nomemcopy)

		fmt.Fprintln(os.Stdout, "[i] Done.")
	}
}

// parseArgs processes and verifies commandline arguments.
func parseArgs() {
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "%s v%s (Go runtime %s)\n",
			AppName, AppVersion, runtime.Version())
		os.Exit(0)
	}

	pipe = len(*in) == 0

	if !pipe && len(*out) == 0 {
		// Ensure we create our own output filename that does not already exist.
		dir, file := path.Split(*in)

		*out = path.Join(dir, file) + ".go"
		if _, err := os.Lstat(*out); err == nil {
			// File already exists. Pad name with a sequential number until we
			// find a name that is available.
			count := 0
			for {
				f := path.Join(dir, fmt.Sprintf("%s.%d.go", file, count))
				if _, err := os.Lstat(f); err != nil {
					*out = f
					break
				}

				count++
			}
		}

		fmt.Fprintf(os.Stderr, "[w] No output file specified. Using '%s'.\n", *out)
	}

	if len(*pkgname) == 0 {
		fmt.Fprintln(os.Stderr, "[w] No package name specified. Using 'main'.")
		*pkgname = "main"
	} else {
		if unicode.IsDigit(rune((*pkgname)[0])) {
			// Identifier can't start with a digit.
			*pkgname = "_" + *pkgname
		}
	}

	if len(*funcname) == 0 {
		if pipe {
			// Can't infer from input file name in this mode.
			fmt.Fprintln(os.Stderr, "[e] No function name specified.")
			os.Exit(1)
		}

		_, file := path.Split(*in)
		file = strings.ToLower(file)
		file = strings.Replace(file, " ", "_", -1)
		file = strings.Replace(file, ".", "_", -1)
		file = strings.Replace(file, "-", "_", -1)

		if unicode.IsDigit(rune(file[0])) {
			// Identifier can't start with a digit.
			file = "_" + file
		}

		fmt.Fprintf(os.Stderr, "[w] No function name specified. Using '%s'.\n", file)
		*funcname = file
	}
}
