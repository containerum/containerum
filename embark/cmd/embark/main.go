package main

import (
	"fmt"
	"os"
	"path"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/cli/flags"
	"github.com/containerum/containerum/embark/pkg/logger"
	"github.com/containerum/containerum/embark/pkg/tiller"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
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

			var log = logger.StdLogger()
			if install.Debug {
				log = logger.DebugLogger()
			}
			//_ = log
			var config, loadKubeConfigErr = loadKubeConfig(install.KubeConfig)
			if loadKubeConfigErr != nil {
				fmt.Println(loadKubeConfigErr)
				return
			}

			var tillerAddr, tillerInstallErr = tiller.InstallTiller(log, config)
			if tillerInstallErr != nil {
				fmt.Println(tillerInstallErr)
				os.Exit(1)
			}

			if install.Host == "" {
				for _, cluster := range config.Clusters {
					install.Host = cluster.Server
					break
				}
			}

			var client = builder.NewCLient(tillerAddr)
			//var client = builder.NewCLient("")

			if err := client.InstallChartWithDependencies(install.Namespace, install.Dir, install.Values); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Ok!")
		},
	}

	cmd.AddCommand(cmdInstallTiller())

	assertErr(gpflag.ParseTo(&install, cmd.PersistentFlags()))
	return cmd
}

func cmdInstallTiller() *cobra.Command {
	var install = flags.Install{}
	var cmd = &cobra.Command{
		Use: "tiller",
		Run: func(cmd *cobra.Command, args []string) {

			var log = logger.StdLogger()
			if install.Debug {
				log = logger.DebugLogger()
			}

			var config, loadKubeConfigErr = loadKubeConfig(install.KubeConfig)
			if loadKubeConfigErr != nil {
				fmt.Println(loadKubeConfigErr)
				return
			}

			var addr, installTillerErr = tiller.InstallTiller(log, config)
			if installTillerErr != nil {
				fmt.Println(installTillerErr)
				os.Exit(1)
			}

			fmt.Println(addr)
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

func loadKubeConfig(configPath string) (api.Config, error) {
	if configPath == "" {
		var ok = false
		configPath, ok = os.LookupEnv("KUBECONFIG")
		if !ok {
			configPath = path.Join(os.Getenv("HOME"), ".kube", "config")
		}
	}
	var config, configLoadFilerErr = clientcmd.LoadFromFile(configPath)
	if configLoadFilerErr != nil {
		fmt.Printf("unable to load kube config from %q: %v\n", configPath, configLoadFilerErr)
		os.Exit(1)
	}
	return *config, nil
}
