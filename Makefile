# This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
# license. Its contents can be found at:
# http://creativecommons.org/publicdomain/zero/1.0

include $(GOROOT)/src/Make.inc

TARG = bindata
GOFILES = main.go gowriter.go bindata.go

all:
	$(GC) -o $(TARG).6 $(GOFILES)
	$(LD) -s -o $(TARG) $(TARG).6

clean:
	rm -rf *.o *.a *.[568vq] [568vq].out *.cgo1.go *.cgo2.c _cgo_defun.c _cgo_gotypes.go _cgo_export.* *.so *.exe $(TARG)
