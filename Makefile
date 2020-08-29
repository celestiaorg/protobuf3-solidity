GO = env GO111MODULE=on go
GOBUILD = $(GO) build
GOTEST = $(GO) test
BIN_NAME = protoc-gen-sol
BIN_DIR = bin

all: build

build:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BIN_NAME) -v
