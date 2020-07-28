PROJECT_NAME :=korgi
GOFILES:=$(shell find . -name '*.go' | grep -v -E '(./vendor)')
 
run:
	go run main.go

all: clean check bin


bin: $(GOFILES)
	mkdir -p  bin/linux/
	GOOS=linux GOARCH=amd64 go build -a -tags musl -ldflags "$(LDFLAGS)" -o bin/linux/korgi main.go

gofmt:
	gofmt -w -s pkg/
	gofmt -w -s cmd/

test:
	 GOOS=linux GOARCH=amd64 go test ./...

check:
	@find . -name vendor -prune -o -name '*.go' -exec gofmt -s -d {} +
	@go vet $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

dep:
	go get -d -v


clean:
		rm -rf bin


.PHONY: all
