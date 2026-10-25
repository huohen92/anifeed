GO ?= go
BIN ?= bin/anifeed
VERSION ?= 0.02
LDFLAGS := -s -w -X github.com/zhifengle/anifeed/version.Version=$(VERSION)

.PHONY: all build run test fmt vet docker-build docker-up clean

all: build

build:
	CGO_ENABLED=1 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) .

run:
	CGO_ENABLED=1 $(GO) run . server

test:
	CGO_ENABLED=1 $(GO) test ./...

fmt:
	gofmt -w cmd config p115 request rsssite store version webui

vet:
	$(GO) vet ./...

docker-build:
	docker build -t anifeed:$(VERSION) .

docker-up:
	docker compose up -d --build

clean:
	rm -rf bin
