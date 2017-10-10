
_SOURCES = $(wildcard *.go)

.PHONY: all
all: Textured

Textured: $(_SOURCES)
	gofmt -s -w $(_SOURCES)
	goimports -w $(_SOURCES)
	go build .
