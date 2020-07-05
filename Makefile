run: build
	./mycorrhiza wiki

build:
	go build .

test:
	go test ./util

help:
	echo "Read the Makefile to see what it can do. It is simple."

