
_EXT = $(shell uname -m)

include examples/Textured/Makefile.mk

.PHONY: examples
examples: Textured
