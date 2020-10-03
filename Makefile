GO := env GO111MODULE=on go
GOBUILD := $(GO) build
GOTEST := $(GO) test
PROTOC := protoc
BIN_DIR := bin

LDFLAGS := -ldflags "-X main.version=$(shell git describe --tags)"

TARGET_GEN_SOL := protoc-gen-sol
TARGETS := $(TARGET_GEN_SOL)

TESTS_PASSING := $(sort $(wildcard test/pass/*))
TESTS_FAILING := $(sort $(wildcard test/fail/*))

all: build test

build: $(TARGETS)

$(TARGETS):
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -v $(LDFLAGS) -o $(BIN_DIR)/ ./cmd/$@

test-go: $(TARGETS)
	$(GOTEST) -mod=readonly ./...

test-protoc: test-protoc-check $(TESTS_PASSING) $(TESTS_FAILING)

test-protoc-check:
	$(PROTOC) --version > /dev/null

$(TESTS_PASSING): build
	$(PROTOC) --plugin $(BIN_DIR)/$(TARGET_GEN_SOL) --sol_out license=Apache-2.0,generate=decoder:$@ -I $@ $@/*.proto;

$(TESTS_FAILING): build
	! $(PROTOC) --plugin $(BIN_DIR)/$(TARGET_GEN_SOL) --sol_out $@ -I $@ $@/*.proto;
