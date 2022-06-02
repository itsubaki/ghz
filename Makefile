SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

install:
	-rm ${GOPATH}/bin/ghz
	go get -u
	go mod tidy
	go install -ldflags "${LDFLAGS}"

test:
	go test -v -cover $(shell go list ./... | grep pkg) -coverprofile=coverage-pkg.out -covermode=atomic

merge:
	echo "" > coverage.txt
	cat coverage.out     >> coverage.txt
	cat coverage-pkg.out >> coverage.txt
