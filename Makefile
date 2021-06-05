WIKI=~/src/example-wiki

run: build
	./mycorrhiza ${WIKI}

config_run: build
	./mycorrhiza -config-path "assets/config.ini" ${WIKI}

devconfig_run: build
	./mycorrhiza -config-path "assets/devconfig.ini" ${WIKI}

build:
	go generate
	go build .

test:
	go test .

help:
	echo "Read the Makefile to see what it can do. It is simple."

