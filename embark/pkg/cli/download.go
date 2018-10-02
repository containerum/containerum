package cli

import (
	"log"

	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/containerum/containerum/embark/pkg/installer"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Download() *cobra.Command {
	var config flags.Download
	var cmd = &cobra.Command{
		Use: "download",
		Run: func(cmd *cobra.Command, args []string) {
			var inst = installer.Installer{
				ContainerumConfigPath: config.Config,
				TempDir:               config.Dir,
			}
			var components, loadContainerumConfigErr = inst.LoadContainerumConfig()
			if loadContainerumConfigErr != nil {
				log.Fatal(loadContainerumConfigErr)
			}
			if err := inst.DownloadComponents(components); err != nil {
				log.Fatal(err)
			}
		},
	}
	if err := gpflag.ParseTo(&config, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}
