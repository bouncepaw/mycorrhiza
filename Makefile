run: build
	./mycorrhiza metarrhiza

config_run: build
	./mycorrhiza -config-path "assets/config.ini" metarrhiza

devconfig_run: build
	./mycorrhiza -config-path "assets/devconfig.ini" metarrhiza

build:
	go generate
	go build .

test:
	go test .

help:
	echo "Read the Makefile to see what it can do. It is simple."

