COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: all
all: build

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: build
build: build-api build-scraper

.PHONY: build-scraper
build-scraper:
	CGO_ENABLED=0 go build -ldflags '-w $(LDFLAGS)' -o krautreporter-scraper cmd/scraper/scraper.go

.PHONY: build-api
build-api:
	CGO_ENABLED=0 go build -ldflags '-w $(LDFLAGS)' -o krautreporter-api cmd/api/api.go

.PHONY: test
test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

.PHONY: fmt
fmt:
	go fmt $(PACKAGES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

# postgres helpers

.PHONY: dump
dump:
	pg_dump -h localhost -p 54321 -U postgres postgres > krautreporter.sql

.PHONY: import
restore:
	psql -h localhost -p 5432 -U postgres postgres < krautreporter.sql
