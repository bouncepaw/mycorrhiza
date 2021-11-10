.POSIX:
include config.example.mk
-include config.mk

mycorrhiza:
	go build .

generate:
	go generate

run: mycorrhiza
	./mycorrhiza ${WIKIPATH}

dev: generate run

check:
	go test .

.PHONY: mycorrhiza generate run dev check
