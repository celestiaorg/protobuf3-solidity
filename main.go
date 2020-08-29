package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// req := &plugin.CodeGeneratorRequest{}
	// resp := &plugin.CodeGeneratorResponse{}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
