package generator

import "google.golang.org/protobuf/types/pluginpb"

// SolidityVersionString is the Solidity version specifier.
const SolidityVersionString = ">=0.6.0 <8.0.0"

// Generator generates Solidity code from a .proto file.
type Generator struct {
	request *pluginpb.CodeGeneratorRequest
}

// New initializes a new Generator.
func New(request *pluginpb.CodeGeneratorRequest) *Generator {
	g := new(Generator)
	g.request = request
	return g
}

// Generate ...
func (g *Generator) Generate() (*pluginpb.CodeGeneratorResponse, error) {
	println(g.request.String())

	return nil, nil
}
