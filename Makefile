# ---------------------------------------------------------------------------
# Makefile for CLI utilities
# 
# Escape '#' and '[' characters with '\', and '$' characters with '$$'
# ---------------------------------------------------------------------------

BUILD_TAG=$(shell git describe --tags 2>/dev/null || echo unreleased)
LDFLAGS=-ldflags=all="-X main.version=${BUILD_TAG} -X main.name=watchenv -s -w"

all: build

build:
	go build -mod vendor ${LDFLAGS}

test:
	go test -mod vendor -v -cover

cover:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

install:
	go install ${LDFLAGS} ./...

update:
	go get -u
	go mod tidy
	@# 'go mod tidy' should update the vendor directory (https://github.com/golang/go/issues/45161)
	go mod vendor

snapshot:
	goreleaser --snapshot --skip-publish --clean

release: 
	@echo releasing ${BUILD_TAG}
	@sed '1,/\#\# \[${BUILD_TAG}/d;/^\#\# /Q' CHANGELOG.md > releaseinfo
	@cat releaseinfo
	goreleaser release --clean --release-notes=releaseinfo
	@rm -f releaseinfo

clean:
	go clean
	rm -rf dist
	rm -f releaseinfo
	rm -f coverage.out

.PHONY: all test clean
