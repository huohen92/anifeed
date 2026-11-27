GO ?= go
BIN ?= bin/anifeed
VERSION ?= 0.03.9
LDFLAGS := -s -w -X github.com/huohen92/anifeed/version.Version=$(VERSION)

.PHONY: all build web web-version run test check fmt vet docker-build docker-up clean

all: build

web:
	./scripts/build-web.sh

web-version:
	@test -n "$(ASSET_VERSION)" || (echo "usage: make web-version ASSET_VERSION=0.03.7.1" >&2; exit 1)
	./scripts/set-asset-version.sh $(ASSET_VERSION)

build:
	CGO_ENABLED=1 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) .

run:
	CGO_ENABLED=1 $(GO) run . server

test:
	CGO_ENABLED=1 $(GO) test ./...

check:
	CGO_ENABLED=1 $(GO) vet ./...
	CGO_ENABLED=1 $(GO) test ./...
	./scripts/check-assets.sh

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
