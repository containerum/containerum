package cli

import (
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/containerum/containerum/embark/pkg/installer"
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
			var installErr = installer.Installer{
				ContainerumConfigPath: installConfig.Containerum,
				KubectlConfigPath:     installConfig.KubeConfig,
				TempDir:               installConfig.Dir,
			}.Install()
			if installErr != nil {
				fer.Fatal("unable to install containerum:\n\t%v", installErr)
			}
		},
	}
	if err := gpflag.ParseTo(&installConfig, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
