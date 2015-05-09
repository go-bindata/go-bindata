all: regen check

regen:
	go install ./...
	make -C testdata regen

.PHONY: check
check: errcheck
	errcheck testdata/out/compress-memcopy.go
	errcheck testdata/out/compress-nomemcopy.go
	errcheck testdata/out/debug.go
	errcheck testdata/out/nocompress-memcopy.go
	errcheck testdata/out/nocompress-nomemcopy.go

errcheck:
	go get github.com/kisielk/errcheck
