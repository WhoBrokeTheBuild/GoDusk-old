
_SOURCES = $(shell find . -name '*.go' | grep -v '.gen')

.PHONY: all
all: gofmt goimports dusk examples

.PHONY: dusk
dusk:
	cd dusk && go build

.PHONY: gofmt
gofmt:
	gofmt -s -w $(_SOURCES)

.PHONY: goimports
goimports:
	goimports -w $(_SOURCES)

include examples/Makefile.mk
