package main

import (
	"fmt"

	"github.com/containerum/containerum/embark/pkg/builder"
)

func main() {
	var client = builder.NewCLient(":1245")
	fmt.Println(client.Install())
}
