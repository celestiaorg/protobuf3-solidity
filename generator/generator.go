package generator

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// SolidityVersionString is the Solidity version specifier.
const SolidityVersionString = ">=0.6.0 <8.0.0"

// Generator generates Solidity code from .proto files.
type Generator struct {
	request *pluginpb.CodeGeneratorRequest
}

// New initializes a new Generator.
func New(request *pluginpb.CodeGeneratorRequest) *Generator {
	g := new(Generator)
	g.request = request
	return g
}

// Generate generates Solidity code from the requested .proto files.
func (g *Generator) Generate() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	protoFiles := g.request.GetProtoFile()
	responseFiles := make([]*pluginpb.CodeGeneratorResponse_File, len(protoFiles))

	for i := 0; i < len(protoFiles); i++ {
		responseFile, err := generateFile(protoFiles[i])
		if err != nil {
			return nil, err
		}

		responseFiles[i] = responseFile
	}

	return responseFiles, nil
}

// generateFile generates Solidity code from a single .proto file.
func generateFile(protoFile *descriptorpb.FileDescriptorProto) (*pluginpb.CodeGeneratorResponse_File, error) {
	err := checkSyntaxVersion(protoFile.GetSyntax())
	if err != nil {
		return nil, err
	}

	responseFile := &pluginpb.CodeGeneratorResponse_File{}

	b := &WriteableBuffer{}

	// TODO option for license
	b.P(fmt.Sprintf("// SPDX-License-Identifier: %s", "CC0"))
	b.P(SolidityVersionString)
	b.P()
	b.P("import \"@lazyledger/protobuf3-solidity-lib/contracts/ProtobufLib.sol\";")
	b.P()

	for i := 0; i < len(protoFile.GetMessageType()); i++ {
		err := generateMessage(protoFile.GetMessageType()[i], b)
		if err != nil {
			return nil, err
		}
	}

	// TODO add b to response
	println(b.String())

	return responseFile, nil
}

func generateMessage(descriptor *descriptorpb.DescriptorProto, b *WriteableBuffer) error {
	println(descriptor.GetName())
	// println(descriptor.GetField()[0].String())

	// TODO generate header, imports, and struct definition
	structName := descriptor.GetName()
	err := checkKeyword(structName)
	if err != nil {
		return err
	}

	b.P(fmt.Sprintf("struct %s {", structName))

	// TODO loop over fields

	b.P("}")
	b.P()

	return nil
}

func checkSyntaxVersion(v string) error {
	if v == "proto3" {
		return nil
	}

	return errors.New("Must use syntax = \"proto3\";")
}

func checkKeyword(w string) error {
	switch w {
	case "after",
		"alias",
		"apply",
		"auto",
		"case",
		"copyof",
		"default",
		"define",
		"final",
		"immutable",
		"implements",
		"in",
		"inline",
		"let",
		"macro",
		"match",
		"mutable",
		"null",
		"of",
		"partial",
		"promise",
		"reference",
		"relocatable",
		"sealed",
		"sizeof",
		"static",
		"supports",
		"switch",
		"typedef",
		"typeof",
		"unchecked":
		return errors.New("Using Solidity keyword " + w)
	}

	return nil
}
