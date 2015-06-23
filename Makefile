ARCH=$(shell ls $(GOPATH)/pkg | head -n 1)
PKG=$(GOPATH)/pkg/$(ARCH)/github.com/goanywhere

all: test

clean:
	@find $(PKG) -name 'rex.a' -delete
	@find $(PKG) -name 'rex' -type d -print0|xargs -0 rm -r

build:
	@go get -v ./...

test:
	@go test -v ./...
