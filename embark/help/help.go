package help

import (
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//go:generate go build -v -o ../bin/fileb0x ../vendor/github.com/UnnoTed/fileb0x
//go:generate ../bin/fileb0x b0x.toml

type Help struct {
	Short    string
	Long     string
	Examples []string
}

func GetHelp(command string) (help Help, ok bool) {
	var data, readFileErr = ReadFile(command + ".yaml")
	if readFileErr != nil {
		return Help{}, false
	}
	return help, yaml.Unmarshal(data, &help) == nil
}

func FillHelps(root *cobra.Command) {
	var stack = []*cobra.Command{root}
	for len(stack) > 0 {
		var peak = len(stack) - 1
		var command = stack[peak]
		stack[peak] = nil
		stack = stack[:peak]
		if command == nil {
			continue
		}
		if help, ok := GetHelp(command.Use); ok {
			command.Short = defaultStr(help.Short, command.Short)
			command.Long = defaultStr(help.Long, command.Long)
			command.Example = strings.Join(help.Examples, "\n")
		}
		stack = append(stack, command.Commands()...)
	}
}

func defaultStr(str, def string) string {
	if str == "" {
		return def
	}
	return str
}
