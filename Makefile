GO = env GO111MODULE=on go
GOBUILD = $(GO) build
GOTEST = $(GO) test
BIN_DIR = bin

TARGETS := protoc-gen-sol

all: build

build: $(TARGETS)

$(TARGETS):
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -v -o $(BIN_DIR)/ ./cmd/$@
