package main

import (
	"io/ioutil"
	"os"

	"github.com/lazyledger/protobuf3-solidity/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Read marshaled request from stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Initialize request object
	request := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, request)
	if err != nil {
		panic(err)
	}

	// Initialize generator with request
	g := generator.New(request)

	// Generate responses
	responses, err := g.Generate()
	if err != nil {
		panic(err)
	}

	// For each response
	for i := 0; i < len(responses); i++ {
		// Marshal response for output
		data, err = proto.Marshal(responses[i])
		if err != nil {
			panic(err)
		}

		// Write response to stdout
		_, err = os.Stdout.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
