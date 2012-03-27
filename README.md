## bindata

This tool converts any file into managable Go source code. Useful for embedding
binary data into a go program. The file data is gzip compressed before being
converted to a raw byte slice.

### Usage

The simplest invocation is to pass it only the input file name.
The output file and code settings are inferred from this automatically.

    $ bindata -i testdata/gophercolor.png
    [w] No output file specified. Using 'testdata/gophercolor.png.go'.
    [w] No package name specified. Using 'main'.
    [w] No function name specified. Using 'gophercolor_png'.
    [i] Done.

This creates the "testdata/gophercolor.png.go" file which has a package
declaration with name 'main' and one function named 'gophercolor_png'.
It looks like this:

     // gophercolor_png returns the decompressed binary data.
     // It panics if an error occurred.
     func gophercolor_png() []byte {
	      gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
              ...
          }))
          
          if err != nil {
              panic("Decompression failed: " + err.Error())
          }

          var b bytes.Buffer
          io.Copy(&b, gz)
          gz.Close()
          
          return b.Bytes()
     }

You can now simply include the new .go file in your program and call
gophercolor_png() to get the uncompressed image data. The function panics
if something went wrong during decompression. This makes any faults appearant
during initialization of your program, so it can quickly be fixed. Additionally,
this approach allows us to assign the decompressed file data to global
variables where necessary.

See the testdata directory for example input and output.

Aternatively, you can pipe the input file data into stdin. bindata will then
spit out the generated Go code to stdout. This does require explicitly naming
the desired function name, as it can not be inferred from the input data.
The package name will still default to 'main'.

     $ cat testdata/gophercolor.png | ./bindata -f gophercolor_png | gofmt

Invoke the program with the -h flag for more options.


### Optional compression

When the `-u` flag is given, the supplied resource is *not* GZIP compressed
before being turned into Go code. This also alters the generated output in that
we no longer need a function that decompresses the data. The resource's raw
byte data is simply assigned to a global variable of the same name as the 
function would otherwise get. This feature is useful if you do not care for
compression, or the supplied resource is already compressed. Doing it again
would not add any value and may even increase the size of the data.

The default behaviour of the program is to use compression.

