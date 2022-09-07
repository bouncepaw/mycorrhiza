.POSIX:
.SUFFIXES:

PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
MANDIR=$(PREFIX)/share/man
GO=go
LDFLAGS=-X "github.com/bouncepaw/mycorrhiza/version.taggedRelease=$$(git describe --tags --abbrev=0)"

all: mycorrhiza

mycorrhiza:
	$(GO) generate $(GOFLAGS)
	CGO_ENABLED=0 $(GO) build -ldflags="$(LDFLAGS)" $(GOFLAGS) -o mycorrhiza .

install:
	mkdir -m755 -p $(DESTDIR)$(BINDIR) $(DESTDIR)$(MANDIR)/man1
	install -m755 mycorrhiza $(DESTDIR)$(BINDIR)/mycorrhiza
	install -m644 help/mycorrhiza.1 $(DESTDIR)$(MANDIR)/man1/mycorrhiza.1

.PHONY: all mycorrhiza install
