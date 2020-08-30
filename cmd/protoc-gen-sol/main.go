package main

import (
	"io/ioutil"
	"os"

	"github.com/lazyledger/protobuf3-solidity/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	request := &pluginpb.CodeGeneratorRequest{}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = proto.Unmarshal(data, request)
	if err != nil {
		panic(err)
	}

	g := generator.New(request)
	g.Run()
}
