run: build
	./mycorrhiza metarrhiza

auth_run: build
	./mycorrhiza -auth-method fixed metarrhiza

gemini_run: build
	./mycorrhiza -gemini-cert-path "." metarrhiza

build:
	go generate
	go build .

test:
	go test .

help:
	echo "Read the Makefile to see what it can do. It is simple."

