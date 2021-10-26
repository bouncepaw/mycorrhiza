.POSIX:
include config.example.mk
include config.mk

run: build
	./mycorrhiza ${WIKIPATH}

config_run: build
	./mycorrhiza ${WIKIPATH}

devconfig_run: build
	./mycorrhiza ${WIKIPATH}

build:
	go generate
	go build .

test:
	go test .
