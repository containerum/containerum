package main

import (
	"fmt"
	"os"
	"path"

	"io/ioutil"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	kubeClientAPI "k8s.io/client-go/tools/clientcmd/api"
)

func main() {
	if err := cmdDownloadRequirements().Execute(); err != nil {
		fmt.Println(err)
	}
}

func cmdDownloadRequirements() *cobra.Command {
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
			var kubeConfigData, err = ioutil.ReadFile(install.KubeConfig)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			var config kubeClientAPI.Config
			if err := yaml.Unmarshal(kubeConfigData, &config); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := client.InstallTiller(config); err != nil {
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
