VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build test test-cover lint clean install

build:
	go build $(LDFLAGS) -o cpk ./cmd/cpk

test:
	go test ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...

clean:
	rm -f cpk cpk-go coverage.out coverage.html

install: build
	cp cpk $(GOPATH)/bin/cpk 2>/dev/null || cp cpk /usr/local/bin/cpk
