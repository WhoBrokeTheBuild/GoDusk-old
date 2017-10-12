
_TEXTURED_DIR = examples/Textured
_TEXTURED_OUT = Textured.$(_EXT)

.PHONY: Textured
Textured: $(_TEXTURED_OUT)

.PHONY: run-Textured
run-Textured: $(_TEXTURED_OUT)
	cd $(_TEXTURED_DIR) && ./$(_TEXTURED_OUT)

$(_TEXTURED_OUT):
	go-bindata -o $(_TEXTURED_DIR)/assets.gen.go -prefix $(_TEXTURED_DIR) $(_TEXTURED_DIR)/assets/...
	cd $(_TEXTURED_DIR) && go build -o $(_TEXTURED_OUT)
