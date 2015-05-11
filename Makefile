all: regen check

regen:
	go install ./...
	make -C testdata regen

.PHONY: check
check: errcheck go-nyet golint
	errcheck testdata/out/compress-memcopy.go
	errcheck testdata/out/compress-nomemcopy.go
	errcheck testdata/out/debug.go
	errcheck testdata/out/nocompress-memcopy.go
	errcheck testdata/out/nocompress-nomemcopy.go
	go-nyet testdata/out/compress-memcopy.go
	go-nyet testdata/out/compress-nomemcopy.go
	go-nyet testdata/out/debug.go
	go-nyet testdata/out/nocompress-memcopy.go
	go-nyet testdata/out/nocompress-nomemcopy.go
	golint testdata/out/compress-memcopy.go
	golint testdata/out/compress-nomemcopy.go
	golint testdata/out/debug.go
	golint testdata/out/nocompress-memcopy.go
	golint testdata/out/nocompress-nomemcopy.go

errcheck:
	go get github.com/kisielk/errcheck

go-nyet:
	go get github.com/barakmich/go-nyet

golint:
	go get github.com/golang/lint/golint
