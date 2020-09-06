GO = env GO111MODULE=on go
GOBUILD = $(GO) build
GOTEST = $(GO) test
PROTOC = protoc
BIN_DIR = bin

LDFLAGS := -ldflags "-X main.version=$(shell git describe --tags)"

TARGETS := protoc-gen-sol

all: build test

build: $(TARGETS)

$(TARGETS):
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -v $(LDFLAGS) -o $(BIN_DIR)/ ./cmd/$@

test-go: $(TARGETS)
	$(GOTEST) -mod=readonly ./...

test-protoc: test-protoc-check test-protoc-pass test-protoc-fail

test-protoc-check:
	$(PROTOC) --version > /dev/null

test-protoc-pass:

test-protoc-fail:
