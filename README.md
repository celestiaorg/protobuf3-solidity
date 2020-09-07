# protobuf3-solidity

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lazyledger/protobuf3-solidity)

A [protobuf3](https://developers.google.com/protocol-buffers) code generator plugin for [Solidity](https://github.com/ethereum/solidity). Decode and encode protobuf messages in your Solidity contract! Leverages the [protobuf3-solidity-lib](https://github.com/lazyledger/protobuf3-solidity-lib) codec library.

Serialization rules are stricter than default protobuf3 rules, and are specified in [ADR-027](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-027-deterministic-protobuf-serialization.md). The resulting serialization is bijective (one-to-one), rather than the usual non-deterministic and malleable serialization used in most protobuf parsers. This makes it suitable for canonical serialization in blockchain applications.

## Usage

Use as a `protoc` plugin:
```sh
protoc --plugin protoc-gen-sol --sol_out [license=<license string>:]<output directory> <proto files>
```

Examples:
```sh
# Output foo.proto.sol in current directory
protoc --plugin protoc-gen-sol --sol_out . foo.proto

# Generate Solidity file with Apache-2.0 license identifier
protoc --plugin protoc-gen-sol --sol_out license=Apache-2.0:. foo.proto
```

### Feature support



## Building from source

Requires [Go](https://golang.org/) `>= 1.14`.

Build:
```sh
make
```

Test:
```sh
make test-protoc
```
