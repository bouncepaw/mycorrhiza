run: build
	./mycorrhiza metarrhiza

run_with_fixed_auth: build
	./mycorrhiza -auth-method fixed metarrhiza

build:
	go generate
	go build .

test:
	go test .

help:
	echo "Read the Makefile to see what it can do. It is simple."

