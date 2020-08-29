package main

import (
	"io/ioutil"
	"os"

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

	println(request.String())
}
