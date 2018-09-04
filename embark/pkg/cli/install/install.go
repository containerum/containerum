package install

import (
	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/containerum/containerum/embark/pkg/utils/fer"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Install(defaultInstallConfig *flags.Install) *cobra.Command {
	var installConfig = flags.Install{}
	if defaultInstallConfig != nil {
		installConfig = *defaultInstallConfig
	}
	var command = &cobra.Command{
		Use:     "install",
		Aliases: []string{"!"}, // embark !
		Run: func(cmd *cobra.Command, args []string) {
			var clientOptions = builder.DefaultClientOptionsPtr()
			if installConfig.Debug {
				clientOptions.Merge(builder.Debug())
			}
			if installConfig.KubeConfig != "" {
				clientOptions.Merge(builder.KubeConfigPath(installConfig.KubeConfig))
			}

			var client, newBuilderClientErr = builder.NewCLient(*clientOptions)
			if newBuilderClientErr != nil {
				fer.Fatal("unable to init client:\n%v\n", newBuilderClientErr)
			}
			var installChartErr = client.InstallChartWithDependencies(
				installConfig.Namespace,
				installConfig.Dir,
				installConfig.Values)
			if installChartErr != nil {
				fer.Fatal("unable to install chart:\n%v\n", installChartErr)
			}
		},
	}
	if err := gpflag.ParseTo(&installConfig, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
