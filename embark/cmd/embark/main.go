package main

import (
	"github.com/containerum/containerum/embark/pkg/cli"
	"github.com/containerum/containerum/embark/pkg/utils/fer"
)

func main() {
	if err := cli.Root().Execute(); err != nil {
		fer.Fatal("%s", err)
	}
}
