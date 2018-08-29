package main

import (
	"fmt"
	"os"
	"path"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := cmdInstall().Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func cmdInstall() *cobra.Command {
	var install = flags.Install{}
	var cmd = &cobra.Command{
		Use: "embark",
		Run: func(cmd *cobra.Command, args []string) {
			var client = builder.NewCLient(install.Host)
			if install.KubeConfig == "" {
				var ok = false
				install.KubeConfig, ok = os.LookupEnv("KUBECONFIG")
				if !ok {
					install.KubeConfig = path.Join(os.Getenv("HOME"), ".kube", "config")
				}
			}
			var config, configLoadFilerErr = clientcmd.LoadFromFile(install.KubeConfig)
			if configLoadFilerErr != nil {
				fmt.Printf("unable to load kube config from %q: %v\n", install.KubeConfig, configLoadFilerErr)
				os.Exit(1)
			}
			if _, err := client.InstallTiller(*config); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if install.Host == "" {
				for _, cluster := range config.Clusters {
					install.Host = cluster.Server
					break
				}
			}
			if err := client.InstallChartWithDependencies(install.Namespace, install.Dir, install.Values); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Ok!")
		},
	}
	assertErr(gpflag.ParseTo(&install, cmd.PersistentFlags()))
	return cmd
}

func assertErr(err error) {
	if err != nil {
		panic(err)
	}
}
