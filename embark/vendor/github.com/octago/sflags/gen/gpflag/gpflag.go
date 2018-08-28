package gpflag

import (
	"os"

	"github.com/octago/sflags"
	"github.com/spf13/pflag"
)

// flagSet describes interface,
// that's implemented by pflag library and required by sflags.
type flagSet interface {
	VarPF(value pflag.Value, name, shorthand, usage string) *pflag.Flag
}

var _ flagSet = (*pflag.FlagSet)(nil)

// GenerateTo takes a list of sflag.Flag,
// that are parsed from some config structure, and put it to dst.
func GenerateTo(src []*sflags.Flag, dst flagSet) {
	for _, srcFlag := range src {
		flag := dst.VarPF(srcFlag.Value, srcFlag.Name, srcFlag.Short, srcFlag.Usage)
		if boolFlag, casted := srcFlag.Value.(sflags.BoolFlag); casted && boolFlag.IsBoolFlag() {
			// pflag uses -1 in this case,
			// we will use the same behaviour as in flag library
			flag.NoOptDefVal = "true"
		}
		flag.Hidden = srcFlag.Hidden
		if srcFlag.Deprecated {
			// we use Usage as Deprecated message for a pflag
			flag.Deprecated = srcFlag.Usage
			if flag.Deprecated == "" {
				flag.Deprecated = "Deprecated"
			}
		}
	}
}

// ParseTo parses cfg, that is a pointer to some structure,
// and puts it to dst.
func ParseTo(cfg interface{}, dst flagSet, optFuncs ...sflags.OptFunc) error {
	flags, err := sflags.ParseStruct(cfg, optFuncs...)
	if err != nil {
		return err
	}
	GenerateTo(flags, dst)
	return nil
}

// Parse parses cfg, that is a pointer to some structure,
// puts it to the new pflag.FlagSet and returns it.
func Parse(cfg interface{}, optFuncs ...sflags.OptFunc) (*pflag.FlagSet, error) {
	fs := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	err := ParseTo(cfg, fs, optFuncs...)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// ParseToDef parses cfg, that is a pointer to some structure and
// puts it to the default pflag.CommandLine.
func ParseToDef(cfg interface{}, optFuncs ...sflags.OptFunc) error {
	err := ParseTo(cfg, pflag.CommandLine, optFuncs...)
	if err != nil {
		return err
	}
	return nil
}
