// +build IntegrationTests

package main

import (
	"sync"
	"testing"

	"github.com/containerum/containerum/embark/pkg/installer"
)

func TestStaticInstallation(test *testing.T) {
	defer lock()()
	if err := (installer.Installer{
		Static: true,
	}).Install(); err != nil {
		test.Fatal(err)
	}
}

var installationMutex = &sync.Mutex{}

func lock() (unlock func()) {
	installationMutex.Lock()
	return func() { installationMutex.Unlock() }
}
