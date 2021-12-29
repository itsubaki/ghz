SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GOVERSION}'

install:
	-rm ${GOPATH}/bin/ghz
	go mod tidy
	go install -ldflags "${LDFLAGS}"

.PHONY: test
test:
	GOOGLE_APPLICATION_CREDENTIALS=../credentials.json go test ./appengine --godog.format=pretty -v -coverprofile=coverage.out -covermode=atomic -coverpkg ./...

testpkg:
	go test -v -cover $(shell go list ./... | grep pkg) -coverprofile=coverage-pkg.out -covermode=atomic

run:
	GOOGLE_APPLICATION_CREDENTIALS=./credentials.json go run appengine/main.go

merge:
	echo "" > coverage.txt
	cat coverage.out     >> coverage.txt
	cat coverage-pkg.out >> coverage.txt

deploy:
	gcloud beta app deploy app.yaml cron.yaml

browse:
	gcloud app browse
