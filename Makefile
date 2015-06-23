.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)
ARCH=$(shell ls $(GOPATH)/pkg | head -n 1)
PKG=$(GOPATH)/pkg/$(ARCH)/github.com/goanywhere

all: test build

clean:
	@find $(PKG) -name 'rex.a' -delete
	@find $(PKG) -name 'rex' -type d -print0|xargs -0 rm -r

build: clean
	@go get -v ./...

test:
	@go test -v ./...

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...
