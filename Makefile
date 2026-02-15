VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build test test-cover lint clean install

build:
	go build $(LDFLAGS) -o chp ./cmd/chp

test:
	go test ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...

clean:
	rm -f chp chp-go coverage.out coverage.html

install: build
	cp chp $(GOPATH)/bin/chp 2>/dev/null || cp chp /usr/local/bin/chp
