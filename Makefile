GO = env GO111MODULE=on go
GOBUILD = $(GO) build
GOTEST = $(GO) test
BIN_DIR = bin

LDFLAGS = -ldflags "-X main.version=`git describe --tags`"

TARGETS := protoc-gen-sol

all: build

build: $(TARGETS)

$(TARGETS):
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -v $(LDFLAGS) -o $(BIN_DIR)/ ./cmd/$@
