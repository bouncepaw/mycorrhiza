.POSIX:
.SUFFIXES:

PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
MANDIR=$(PREFIX)/share/man
GO=go

all: mycorrhiza

mycorrhiza:
	$(GO) generate $(GOFLAGS)
	$(GO) build $(GOFLAGS) . -o mycorrhiza

install:
	mkdir -m755 -p $(DESTDIR)$(BINDIR) $(DESTDIR)$(MANDIR)/man1
	install -m755 mycorrhiza $(DESTDIR)$(BINDIR)/mycorrhiza
	install -m644 doc/mycorrhiza.1 $(DESTDIR)$(MANDIR)/man1/mycorrhiza.1

.PHONY: all mycorrhiza install
