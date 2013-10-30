// +build release

package main

import (
	"reflect"
	"unsafe"
)

func bindata_read(data, name string) []byte {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b
}

var _in_b_test_asset = "\x2f\x2f\x20\x73\x61\x6d\x70\x6c\x65\x20\x66\x69\x6c\x65\x0a"

func in_b_test_asset() []byte {
	return bindata_read(
		_in_b_test_asset,
		"in/b/test.asset",
	)
}

var _in_test_asset = "\x2f\x2f\x20\x73\x61\x6d\x70\x6c\x65\x20\x66\x69\x6c\x65\x0a"

func in_test_asset() []byte {
	return bindata_read(
		_in_test_asset,
		"in/test.asset",
	)
}

var _in_a_test_asset = "\x2f\x2f\x20\x73\x61\x6d\x70\x6c\x65\x20\x66\x69\x6c\x65\x0a"

func in_a_test_asset() []byte {
	return bindata_read(
		_in_a_test_asset,
		"in/a/test.asset",
	)
}

var _in_c_test_asset = "\x2f\x2f\x20\x73\x61\x6d\x70\x6c\x65\x20\x66\x69\x6c\x65\x0a"

func in_c_test_asset() []byte {
	return bindata_read(
		_in_c_test_asset,
		"in/c/test.asset",
	)
}
