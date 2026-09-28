BINARY  := gitbatch
VERSION ?= dev
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all build test race lint fmt vet tidy clean

all: vet lint test build

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/gitbatch

test:
	go test ./...

race:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w cmd internal

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -f $(BINARY) coverage.txt
	rm -rf dist
