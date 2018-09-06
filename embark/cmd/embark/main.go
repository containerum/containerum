package main

import (
	"github.com/containerum/containerum/embark/pkg/cli/install"
	"github.com/containerum/containerum/embark/pkg/utils/fer"
)

func main() {
	if err := install.Install(nil).Execute(); err != nil {
		fer.Fatal("%v", err)
	}
}
