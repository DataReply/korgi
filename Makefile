PROJECT_NAME  := korgi
SHELL         := /usr/bin/env bash

OS            := "$(shell uname | tr '[:upper:]' '[:lower:]')"
ARCH          := amd64
LDFLAGS       := -w -s
SRC           := $(shell find . -type f -name '*.go' -print)
REPO          := github.com/DataReply/korgi

GIT_COMMIT    := $(shell git rev-parse HEAD)
GIT_SHA       := $(shell git rev-parse --short HEAD)
GIT_TAG       := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY     := $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X ${REPO}/internal/version.version=${BINARY_VERSION}
endif

VERSION_METADATA = unreleased
# Clear the "unreleased" string in BuildMetadata if there's a tag
ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X ${REPO}/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X ${REPO}/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X ${REPO}/internal/version.gitTreeState=${GIT_DIRTY}

bin: $(SRC)
	mkdir -p  bin/${OS}/
	GOOS=${OS} GOARCH=${ARCH} go build -ldflags "$(LDFLAGS)" -o bin/${OS}/korgi main.go

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

