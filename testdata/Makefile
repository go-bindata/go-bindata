.PHONY: check
check: out/compress-memcopy.go out/compress-nomemcopy.go out/debug.go out/nocompress-memcopy.go out/nocompress-nomemcopy.go errcheck go-nyet golint
	$(foreach file, $(wildcard out/*.go), errcheck $(file);)
	$(foreach file, $(wildcard out/*.go), go-nyet $(file);)
	$(foreach file, $(wildcard out/*.go), golint $(file);)

$(GOPATH)/bin/go-bindata: $(wildcard ../*.go) $(wildcard ../**/*.go)
	go install ../...

out/compress-memcopy.go: 		 $(wildcard in/**/*) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata                        -o $@ in/...

out/compress-nomemcopy.go: 	 $(wildcard in/**/*) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -nomemcopy             -o $@ in/...

out/debug.go: 							 $(wildcard in/**/*) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -debug                 -o $@ in/...

out/nocompress-memcopy.go:   $(wildcard in/**/*) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -nocompress            -o $@ in/...

out/nocompress-nomemcopy.go: $(wildcard in/**/*) $(GOPATH)/bin/go-bindata
	$(GOPATH)/bin/go-bindata -nocompress -nomemcopy -o $@ in/...

errcheck:
	go get github.com/kisielk/errcheck

go-nyet:
	go get github.com/barakmich/go-nyet

golint:
	go get github.com/golang/lint/golint
