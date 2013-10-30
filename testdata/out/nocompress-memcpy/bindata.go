package main

// Asset loads and returns the asset for the given name.
// This returns nil of the asset could not be found.
func Asset(name string) []byte {
	if f, ok := _bindata[name]; ok {
		return f()
	}
	return nil
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() []byte{
	"in/b/test.asset": in_b_test_asset,
	"in/test.asset":   in_test_asset,
	"in/a/test.asset": in_a_test_asset,
	"in/c/test.asset": in_c_test_asset,
}
