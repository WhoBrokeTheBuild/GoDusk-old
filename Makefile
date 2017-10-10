.PHONY: all
all: 
	@gofmt -s -w dusk/*.go
	@goimports -w dusk/*.go
	@cd dusk && go build -o libdusk.a .

