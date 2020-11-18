PROJECT_NAME  := korgi
SHELL         := /usr/bin/env bash

OS            := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH          := amd64
LDFLAGS       := -w -s
SRC           := $(shell find . -type f -name '*.go' -print)
REPO          := github.com/DataReply/korgi
 
bin: $(SRC)
	mkdir -p  bin/${OS}/
	GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -o bin/${OS}/korgi main.go

gofmt:
	gofmt -w -s pkg/
	gofmt -w -s cmd/

test:
	 GOOS=${OS} GOARCH=${ARCH} go test ./...

check:
	@find . -name vendor -prune -o -name '*.go' -exec gofmt -s -d {} +
	@go vet $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

dep:
	go get -d -v

clean:
		rm -rf bin

run:
	go run main.go

all: clean check bin

.PHONY: all

