
_TEXTURED_DIR = examples/Textured
_TEXTURED_OUT = Textured.$(_EXT)
_TEXTURED_ASSETS = $(shell find $(_TEXTURED_DIR)/assets/ -type f)
_TEXTURED_BINDATA = $(_TEXTURED_DIR)/assets.gen.go

.PHONY: Textured
Textured: $(_TEXTURED_OUT)

.PHONY: run-Textured
run-Textured: $(_TEXTURED_OUT)
	cd $(_TEXTURED_DIR) && ./$(_TEXTURED_OUT)

$(_TEXTURED_BINDATA): $(_TEXTURED_ASSETS)
	go-bindata -o $(_TEXTURED_BINDATA) -prefix $(_TEXTURED_DIR) $(_TEXTURED_DIR)/assets/...

$(_TEXTURED_OUT): $(_TEXTURED_BINDATA)
	cd $(_TEXTURED_DIR) && go build -o $(_TEXTURED_OUT)
