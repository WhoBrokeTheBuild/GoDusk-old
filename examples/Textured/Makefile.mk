
_OUT = Textured.$(_EXT)

.PHONY: Textured
Textured: $(_OUT)

.PHONY: run-Textured
run-Textured: $(_OUT)
	cd examples/Textured && ./$(_OUT)

$(_OUT):
	cd examples/Textured && go build -o $(_OUT)
