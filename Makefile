.POSIX:
include config.example.mk
-include config.mk

mycorrhiza:
	go generate
	go build .

run: mycorrhiza
	./mycorrhiza ${WIKIPATH}

check:
	go test .

.PHONY: mycorrhiza run check
