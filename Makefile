VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build test install clean

build:
	go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" -o bin/vibediff ./cmd/vibediff

test:
	go test -v ./...

install:
	go install -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" ./cmd/vibediff

clean:
	rm -rf bin/

# Build for all platforms
release:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" -o bin/vibediff-linux-amd64 ./cmd/vibediff
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" -o bin/vibediff-darwin-amd64 ./cmd/vibediff
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" -o bin/vibediff-darwin-arm64 ./cmd/vibediff
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'" -o bin/vibediff-windows-amd64.exe ./cmd/vibediff
