.PHONY: all clean deps fmt vet test docker

COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: deps build test

clean:
	if [ -f api ] ; then rm -f api ; fi
	if [ -f crawler ] ; then rm -f crawler ; fi

deps:
	go get ./...

lint:
	go fmt $(PACKAGES)
	go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

build: api crawler

api: $(wildcard *.go)
	CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)' -o api cmd/api/api.go

crawler: $(wildcard *.go)
	CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)' -o crawler cmd/crawler/crawler.go
