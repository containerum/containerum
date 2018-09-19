package install

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/containerum/containerum/embark/pkg/models/containerum"
	"github.com/containerum/containerum/embark/pkg/static"
	"github.com/containerum/containerum/embark/pkg/utils/fer"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
			if installConfig.Dir == "" {
				installConfig.Dir = path.Join(os.TempDir(), "embark")
				if err := os.MkdirAll(installConfig.Dir, os.ModePerm|os.ModeDir); !os.IsExist(err) {
					fer.Fatal("unable to create temp dir %q:\n%v\n", installConfig.Dir, err)
				}
			}
			var contData []byte
			var loadContDataErr error

			if installConfig.KubeConfig != "" {
				contData, loadContDataErr = ioutil.ReadFile(installConfig.Containerum)
				if loadContDataErr != nil {
					fer.Fatal("unable to load containerum file: %v", loadContDataErr)
				}
			} else {
				contData, loadContDataErr = static.ReadFile("containerum.yaml")
			}

			var cont containerum.Containerum
			if err := yaml.Unmarshal(contData, &cont); err != nil {
				fer.Fatal("unable to load containerum data: %v", err)
			}
			if err := builder.DowloadComponents(installConfig.Dir, cont); err != nil {
				fer.Fatal("unable to download containerum components: %v", err)
			}
			var components, renderErr = builder.RenderComponents(installConfig.Dir, cont)
			if renderErr != nil {
				fer.Fatal("unable to render containerum components: %v", renderErr)
			}
			var installationGraph, buildGraphErr = builder.BuildGraph(installConfig.Dir, components)
			if buildGraphErr != nil {
				fer.Fatal("unable to build installation graph: %v", installationGraph)
			}
			const rootComponent = "containerum"
			if err := installationGraph.Execute(rootComponent); err != nil {
				fer.Fatal("unable to install %s: %v", rootComponent, err)
			}
		},
	}
	if err := gpflag.ParseTo(&installConfig, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
