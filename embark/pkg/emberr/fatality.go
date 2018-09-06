package emberr

import (
	"os"

	"github.com/containerum/containerum/embark/pkg/utils/fer"
)

type Fatal interface {
	Error
	ExitCoder
}

type ExitCoder interface {
	ExitCode() int
}

var (
	_ ExitCoder = defaultExitCoder{}
)

type defaultExitCoder struct{}

func (defaultExitCoder) ExitCode() (exitCode int) {
	return 1
}

func IsFatal(err error) bool {
	var _, isFatal = err.(Fatal)
	return isFatal
}

func Terminate(err Fatal) {
	fer.Println(err.Error())
	os.Exit(err.ExitCode())
}

func TerminateIfFatal(err error) {
	if IsFatal(err) {
		Terminate(err.(Fatal))
	}
}
