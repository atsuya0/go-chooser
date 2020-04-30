test := go test -v -cover -parallel 4

.PHONY: build test format

build: format
	@go build

test: format
	@$(test)

format:
	@goimports -w .
