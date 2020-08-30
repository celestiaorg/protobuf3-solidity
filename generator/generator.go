package generator

import "google.golang.org/protobuf/types/pluginpb"

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

// Run ...
func (*Generator) Run() {

}
