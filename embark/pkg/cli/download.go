package cli

import (
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	gpflag "github.com/octago/sflags/gen/gpflag"
	cobra "github.com/spf13/cobra"
)

func Download() *cobra.Command {
	var config flags.Download
	var cmd = &cobra.Command{
		Use: "Download",
		Run: func(cmd *cobra.Command, args []string) {
			// write your code here
		},
	}
	if err := gpflag.ParseTo(&config, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}
