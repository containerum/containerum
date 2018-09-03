package install

import (
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Install() *cobra.Command {
	var installConfig = flags.Install{}
	var command = &cobra.Command{
		Use:     "install",
		Aliases: []string{"!"}, // embark !
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	if err := gpflag.ParseTo(&installConfig, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
