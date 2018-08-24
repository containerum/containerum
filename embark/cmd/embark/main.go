package main

import (
	"fmt"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"github.com/spf13/cobra"
)

func main() {
	cmdDownloadRequirements.Execute()
}

var cmdDownloadRequirements = &cobra.Command{
	Use: "requirements",
	Run: func(cmd *cobra.Command, args []string) {
		var requirementsFile, dir = args[0], args[1]
		var client = builder.NewCLient("")
		var req requirements.Requirements
		builder.LoadYAML(requirementsFile, &req)
		var gr, err = client.FetchAllDeps(req, dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(gr.Execute("containerum"))
	},
}
