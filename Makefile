.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: test build

build:
	@go get -v ./...

test: build
	@go test -v ./...

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

