test := go test -v -cover -parallel 4

.PHONY: build test format

build: format
	@go build

test:
	@$(test)

format:
	@goimports -w .
