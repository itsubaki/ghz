SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

update:
	go get -u
	go mod tidy

install:
	-rm ${GOPATH}/bin/ghz
	go install -ldflags "${LDFLAGS}"

test:
	go test -v -cover $(shell go list ./...) -coverprofile=coverage-pkg.out -covermode=atomic
