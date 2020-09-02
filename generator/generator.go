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
func (g *Generator) Generate() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	for i := 0; i < len(g.request.GetProtoFile()); i++ {
		_, err := g.generateFile(i)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// generateFile generates Solidity code from a single .proto file.
func (g *Generator) generateFile(fileIndex int) (*pluginpb.CodeGeneratorResponse_File, error) {
	err := checkSyntaxVersion(g.request.GetProtoFile()[fileIndex].GetSyntax())
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
