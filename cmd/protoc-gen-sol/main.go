package main

import (
	"io/ioutil"
	"os"

	"github.com/lazyledger/protobuf3-solidity/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var version = "unspecified version"

func main() {
	// Read marshaled request from stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		err = responseError(err)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	// Initialize request object
	request := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(data, request)
	if err != nil {
		err = responseError(err)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	// Initialize generator with request
	g := generator.New(request, version)

	// Parse any command-line parameters
	err = g.ParseParameters()
	if err != nil {
		err = responseError(err)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	// Generate response
	response, err := g.Generate()
	if err != nil {
		err = responseError(err)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	err = responsePrint(response)
	if err != nil {
		panic(err)
	}
}

func responseError(err error) error {
	response := &pluginpb.CodeGeneratorResponse{}
	response.Error = proto.String(err.Error())

	err = responsePrint(response)
	if err != nil {
		return err
	}

	return nil
}

func responsePrint(response *pluginpb.CodeGeneratorResponse) error {
	// Marshal response for output
	data, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	// Write response to stdout
	_, err = os.Stdout.Write(data)
	if err != nil {
		return err
	}

	return nil
}
