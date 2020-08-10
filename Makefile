PROJECT_NAME:=korgi
GOFILES:=$(shell find . -name '*.go' | grep -v -E '(./vendor)')
OS:="$(shell uname | tr '[:upper:]' '[:lower:]')"
ARCH:=amd64

run:
	go run main.go

all: clean check bin

bin: $(GOFILES)
	mkdir -p  bin/${OS}/
	GOOS=${OS} GOARCH=${ARCH} go build -a -tags musl -ldflags "$(LDFLAGS)" -o bin/${OS}/korgi main.go

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


.PHONY: all
