test := go test -v -cover -parallel 4

.PHONY: install build test format

install: format
	@go install

build:
	@go build

test:
	@$(test)

format:
	@goimports -w .
