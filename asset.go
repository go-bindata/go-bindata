// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

// File is an asset entry for the table of contents.
type Asset struct {
	Path string // Full file path.
	Name string // Key used in TOC -- name by which asset is referenced.
	Func string // Function name for the procedure returning the asset contents.
}
