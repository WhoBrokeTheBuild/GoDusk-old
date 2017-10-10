
_SOURCES = $(wildcard *.go)
_OUT = Textured.$(shell uname -m)

.PHONY: all
all: $(_OUT) 

$(_OUT): $(_SOURCES) ../../dusk/libdusk.a
	goimports -w $(_SOURCES)
	gofmt -s -w $(_SOURCES)
	go build -o $@ .
