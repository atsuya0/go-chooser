.PHONY: build

build:
	@goimports -w .
	@go build
