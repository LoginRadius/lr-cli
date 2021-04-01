.PHONY: bin/lr
bin/lr: script/build
	@script/build bin/lr

script/build: script/build.go
	go build -o script/build script/build.go

## Site Docs

.PHONY: site-docs
site-docs:
	go run ./app/gen-docs

## Install/uninstall tasks are here for use on *nix platform. On Windows, there is no equivalent.

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin

.PHONY: install
install: bin/lr
	install -d ${DESTDIR}${bindir}
	install -m755 bin/lr ${DESTDIR}${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/lr
