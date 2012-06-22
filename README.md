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
example input and output files for various modes.

Aternatively, you can pipe the input file data into stdin. bindata will then
spit out the generated Go code to stdout. This does require explicitly naming
the desired function name, as it can not be inferred from the input data.
The package name will still default to 'main'.

     $ cat testdata/gophercolor.png | go-bindata -f gophercolor_png | gofmt

Invoke the program with the -h flag for more options.


### Lower memory footprint

Using the `-m` flag, will alter the way the output file is generated.
It will employ a hack that allows us to read the file data directly from
the compiled program's `.rodata` section. This ensures that when we call
call our generate function, we omit unnecessary memcopies.

The downside of this, is that it requires dependencies on the `reflect` and
`unsafe` packages. These may be restricted on platforms like AppEngine and
thus prevent you from using this mode.

Another disadvantage is that the byte slice we create, is strictly read-only.
For most use-cases this is not a problem, but if you ever try to alter the
returned byte slice, a runtime panic is thrown. Use this mode only on target
platforms where memory constraints are an issue.

The default behaviour is to use the old code generation method. This
prevents the two previously mentioned issues, but will employ at least one
extra memcopy and thus increase memory requirements.

For instance, consider the following two examples...

This would be the default mode, using an extra memcopy but gives a safe
implementation without dependencies on `reflect` and `unsafe`:

    func myfile() []byte {
        return []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a}
    }

Here is the same functionality, but uses the .rodata hack.
The byte slice returned from this example can not be written to without
generating a runtime error.

    var _myfile = "\x89\x50\x4e\x47\x0d\x0a\x1a"
    
    func myfile() []byte {
        var empty [0]byte
	    sx := (*reflect.StringHeader)(unsafe.Pointer(&_myfile))
	    b := empty[:]
	    bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	    bx.Data = sx.Data
	    bx.Len = len(_myfile)
	    bx.Cap = bx.Len
	    return b
    }


### Optional compression

When the `-u` flag is given, the supplied resource is *not* GZIP compressed
before being turned into Go code. The data should still be accessed through
a function call, so nothing changes in the usage of the generated file.

This feature is useful if you do not care for compression, or the supplied
resource is already compressed. Doing it again would not add any value and may
even increase the size of the data.

The default behaviour of the program is to use compression.

