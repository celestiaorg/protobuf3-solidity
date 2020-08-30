package generator

import (
	"errors"

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
func (g *Generator) Generate() (*pluginpb.CodeGeneratorResponse_File, error) {
	err := checkSyntaxVersion(g.request.GetProtoFile()[0].GetSyntax())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func checkSyntaxVersion(v string) error {
	if v == "proto3" {
		return nil
	}
	return errors.New("Must use syntax = \"proto3\";")
}
