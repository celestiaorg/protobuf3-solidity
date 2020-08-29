package main

import (
	"fmt"
	"io/ioutil"
	"os"

	plugin "google.golang.org/protobuf/types/pluginpb"
)

func main() {
	req := &plugin.CodeGeneratorRequest{}
	resp := &plugin.CodeGeneratorResponse{}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
