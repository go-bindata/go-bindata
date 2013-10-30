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

var _testdata_in_b_test_asset = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xd7\x57\x28\x4e\xcc\x2d\xc8\x49\x55\x48\xcb\xcc\x49\xe5\x02\x04\x00\x00\xff\xff\x8a\x82\x8c\x85\x0f\x00\x00\x00"

func testdata_in_b_test_asset() []byte {
	return bindata_read(
		_testdata_in_b_test_asset,
		"../testdata/in/b/test.asset",
	)
}

var _testdata_in_a_test_asset = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x2a\x4e\xcc\x2d\xc8\x49\x55\x48\xcb\xcc\x49\xe5\x02\x04\x00\x00\xff\xff\xda\x3d\x49\xdd\x0c\x00\x00\x00"

func testdata_in_a_test_asset() []byte {
	return bindata_read(
		_testdata_in_a_test_asset,
		"../testdata/in/a/test.asset",
	)
}

var _testdata_in_c_test_asset = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xd7\x57\x28\x4e\xcc\x2d\xc8\x49\x55\x48\xcb\xcc\x49\xe5\x02\x04\x00\x00\xff\xff\x8a\x82\x8c\x85\x0f\x00\x00\x00"

func testdata_in_c_test_asset() []byte {
	return bindata_read(
		_testdata_in_c_test_asset,
		"../testdata/in/c/test.asset",
	)
}

