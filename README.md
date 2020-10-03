# protobuf3-solidity

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lazyledger/protobuf3-solidity)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/lazyledger/protobuf3-solidity)](https://github.com/lazyledger/protobuf3-solidity/releases)
[![Go and Node.js Tests](https://github.com/lazyledger/protobuf3-solidity/workflows/Go%20and%20Node.js%20Tests/badge.svg)](https://github.com/lazyledger/protobuf3-solidity/actions?query=workflow%3A%22Go+and+Node.js+Tests%22)
[![GitHub](https://img.shields.io/github/license/lazyledger/protobuf3-solidity)](https://github.com/lazyledger/protobuf3-solidity/blob/master/LICENSE)

A [protobuf3](https://developers.google.com/protocol-buffers) code generator plugin for [Solidity](https://github.com/ethereum/solidity). Decode and encode protobuf messages in your Solidity contract! Leverages the [protobuf3-solidity-lib](https://github.com/lazyledger/protobuf3-solidity-lib) codec library.

Serialization rules are stricter than default protobuf3 rules, and are specified in [ADR-027](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-027-deterministic-protobuf-serialization.md). The resulting serialization is bijective (one-to-one), rather than the usual non-deterministic and malleable serialization used in most protobuf parsers. This makes it suitable for canonical serialization in blockchain applications.

## Usage

Use as a `protoc` plugin:
```sh
protoc \
--plugin protoc-gen-sol \
--sol_out [license=<license string>,compile=<link,inline>,generate=<all,decoder,encoder>:]<output directory> \
<proto files>
```

Examples:
```sh
# Output foo.proto.sol in current directory
protoc --plugin protoc-gen-sol --sol_out . foo.proto

# Generate Solidity file with Apache-2.0 license identifier
protoc --plugin protoc-gen-sol --sol_out license=Apache-2.0:. foo.proto
```

### Parameters

- `license`: default `CC0`
  - any string is accepted, and the generated license comment will use the string as-is
- `compile`: default `inline`
  - `inline`: the generated library's functions will be inlined (`JUMP`)
  - `link`: the generated library's functions will be linked (`DELEGATECALL`)
- `generate`: default `decoder`
  - `all`: both decoder and encoder will be generated
  - `decoder`: only decoder will be generated
  - `encoder`: only encoder will be generated (experimental!)

### Feature support

The below protobuf file shows all supported features of this plugin.
```protobuf
syntax = "proto3";

// import is supported but not shown here

enum OtherEnum {
  UNSPECIFIED = 0;
  ONE = 1;
  TWO = 2;
};

message OtherMessage {
  uint64 other_field = 1;
}

message Message {
  int32 optional_int32 = 1;
  int64 optional_int64 = 2;
  uint32 optional_uint32 = 3;
  uint64 optional_uint64 = 4;
  sint32 optional_sint32 = 5;
  sint64 optional_sint64 = 6;
  fixed32 optional_fixed32 = 7;
  fixed64 optional_fixed64 = 8;
  sfixed32 optional_sfixed32 = 9;
  sfixed64 optional_sfixed64 = 10;
  bool optional_bool = 11;
  string optional_string = 12;
  bytes optional_bytes = 13;
  OtherEnum optional_enum = 14;
  OtherMessage optional_message = 15;

  repeated int32 repeated_int32 = 16 [packed = true];
  repeated int64 repeated_int64 = 17 [packed = true];
  repeated uint32 repeated_uint32 = 18 [packed = true];
  repeated uint64 repeated_uint64 = 19 [packed = true];
  repeated sint32 repeated_sint32 = 20 [packed = true];
  repeated sint64 repeated_sint64 = 21 [packed = true];
  repeated fixed32 repeated_fixed32 = 22 [packed = true];
  repeated fixed64 repeated_fixed64 = 23 [packed = true];
  repeated sfixed32 repeated_sfixed32 = 24 [packed = true];
  repeated sfixed64 repeated_sfixed64 = 25 [packed = true];
  repeated bool repeated_bool = 26 [packed = true];
  repeated OtherEnum repeated_enum = 27 [packed = true];
  repeated OtherMessage repeated_message = 28;
}
```

**Rules to keep in mind:**
1. Enum values must start at `0` and increment by `1`.
1. Field numbers must start at `1` and increment by `1`.
1. Repeated numeric types must explicitly specify `[packed = true]`.

**Currently unsupported features**:
1. nested `enum` or `message` definitions - All `enum` and `message` definitions must be top-level.
1. `package` - Scoping currently unsupported, including in `import`s.

**Unsupported features**:
1. repeated `string` and `bytes` - Solidity does not support arrays of `string` or `bytes`. Workaround: wrap the field in a `message`.
1. `float` and `double` - Solidity does not support floating-point numbers.
1. `oneof` - Solidity does not support unions.
1. `map` - Maps are forbidden as per [ADR-027](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-027-deterministic-protobuf-serialization.md).

## Building from source

Requires [Go](https://golang.org/) `>= 1.14`.

Build:
```sh
make
```

Test (requires a [`protoc`](https://github.com/protocolbuffers/protobuf) binary in `PATH`):
```sh
make test-protoc
```
