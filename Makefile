WIKI=~/src/example-wiki

run: build
	./mycorrhiza ${WIKI}

config_run: build
	./mycorrhiza ${WIKI}

devconfig_run: build
	./mycorrhiza ${WIKI}

build:
	go generate
	go build .

test:
	go test .
