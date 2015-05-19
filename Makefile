.PHONY: .FORCE
include build.mk

GO=go
GOPATH=$(shell pwd)
GO_BUILD_TAGS=""
LDFLAGS="-X main._BUILDDATE_ '$(BUILD_DATE)' -X main._VERSION_ '$(VERSION)'"

all: build

build:
#GOPATH=$(GOPATH) $(GO) install -ldflags $(LDFLAGS) -tags $(GO_BUILD_TAGS) ./src/cmd/scorpio
	GOPATH=$(GOPATH) $(GO) install ./src/cmd/scorpio

clean:
	rm -rf bin pkg

g:
	./script/gen_protos.sh
