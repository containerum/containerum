package main

import (
	"os"
	"os/exec"

	"github.com/containerum/containerum/embark/pkg/utils/fer"
)

//go:generate fileb0x b0x.toml

func main() {
	var fileb0x = exec.Command("fileb0x", "b0x.toml")
	fileb0x.Stdout = os.Stdout
	fileb0x.Stderr = os.Stderr
	var getWDErr error
	fileb0x.Dir, getWDErr = os.Getwd()
	if getWDErr != nil {
		panic(getWDErr)
	}
	if err := fileb0x.Run(); err != nil {
		fer.Fatal("unable to generate static embedded filesystem:\n%v\n", err)
	}
}
