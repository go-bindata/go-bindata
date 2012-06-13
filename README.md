## bindata

This tool converts any file into managable Go source code. Useful for embedding
binary data into a go program. The file data is optionally gzip compressed
before being converted to a raw byte slice.

### Usage

The simplest invocation is to pass it only the input file name.
The output file and code settings are inferred from this automatically.

    $ go-bindata -i testdata/gophercolor.png
    [w] No output file specified. Using 'testdata/gophercolor.png.go'.
    [w] No package name specified. Using 'main'.
    [w] No function name specified. Using 'gophercolor_png'.
    [i] Done.

This creates the `testdata/gophercolor.png.go` file which has a package
declaration with name `main` a variable holding the file data in a read-only
string and one function named `gophercolor_png` with the following signature:

    func gophercolor_png() []byte

You can now simply include the new .go file in your program and call
`gophercolor_png()` to get the uncompressed image data. The function panics
if something went wrong during decompression. See the testdata directory for
example input and output.

Aternatively, you can pipe the input file data into stdin. bindata will then
spit out the generated Go code to stdout. This does require explicitly naming
the desired function name, as it can not be inferred from the input data.
The package name will still default to 'main'.

     $ cat testdata/gophercolor.png | go-bindata -f gophercolor_png | gofmt

Invoke the program with the -h flag for more options.


### Optional compression

When the `-u` flag is given, the supplied resource is *not* GZIP compressed
before being turned into Go code. The data should still be accessed through
a function call, so nothing changes in the usage of the generated file.

This feature is useful if you do not care for compression, or the supplied
resource is already compressed. Doing it again would not add any value and may
even increase the size of the data.

The default behaviour of the program is to use compression.

