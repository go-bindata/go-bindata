## bindata

This tool converts any file into managable Go source code. Useful for embedding
binary data into a go program. The file data is optionally gzip compressed
before being converted to a raw byte slice.


### Usage

TODO

### Lower memory footprint

Using the `-nomemcopy` flag, will alter the way the output file is generated.
It will employ a hack that allows us to read the file data directly from
the compiled program's `.rodata` section. This ensures that when we call
call our generated function, we omit unnecessary memcopies.

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

For instance, consider the following two examples:

This would be the default mode, using an extra memcopy but gives a safe
implementation without dependencies on `reflect` and `unsafe`:

```go
func myfile() []byte {
    return []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a}
}
```

Here is the same functionality, but uses the `.rodata` hack.
The byte slice returned from this example can not be written to without
generating a runtime error.

```go
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
```


### Optional compression

When the `-uncompressed` flag is given, the supplied resource is *not* GZIP compressed
before being turned into Go code. The data should still be accessed through
a function call, so nothing changes in the usage of the generated file.

This feature is useful if you do not care for compression, or the supplied
resource is already compressed. Doing it again would not add any value and may
even increase the size of the data.

The default behaviour of the program is to use compression.


#### Table of Contents keys

The keys used in the `go_bindata` map, are the same as the input file name passed to `go-bindata`.
This includes the fully qualified (absolute) path. In most cases, this is not desireable, as it
puts potentially sensitive information in your code base. For this purpose, the tool supplies
another command line flag `-prefix`. This accepts a portion of a path name, which should be
stripped off from the map keys and function names.

For example, running without the `-prefix` flag, we get:

	$ go-bindata /path/to/templates/foo.html
    
	go_bindata["/path/to/templates/foo.html"] = path_to_templates_foo_html

Running with the `-prefix` flag, we get:

	$ go-bindata -prefix "/path/to/" /path/to/templates/foo.html
    
	go_bindata["templates/foo.html"] = templates_foo_html


#### Build tags

With the optional -tags flag, you can specify any go build tags that
must be fulfilled for the output file to be included in a build. This
is useful for including binary data in multiple formats, where the desired
format is specified at build time with the appropriate tag(s).

The tags are appended to a `// +build` line in the beginning of the output file
and must follow the build tags syntax specified by the go tool.
