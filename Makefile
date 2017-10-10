.PHONY: all
all: dusk

.PHONY: dusk
dusk:
	@gofmt -s -w dusk/*.go
	@goimports -w dusk/*.go
	@cd dusk && go build .

