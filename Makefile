ARCH=$(shell ls $(GOPATH)/pkg | head -n 1)
PKG=$(GOPATH)/pkg/$(ARCH)/github.com/goanywhere

all: test
	@echo $HOST

clean:
	@find $(PKG) -name 'rex.a' -delete
	@find $(PKG) -name 'rex' -type d -print0|xargs -0 rm -r

build: clean
	@go get -v ./...

test: build
	@go test -v ./...
