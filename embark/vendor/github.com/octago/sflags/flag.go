// Package sflags helps to generate flags by parsing structure
package sflags

// Flag structure might be used by cli/flag libraries for their flag generation.
type Flag struct {
	Name       string // name as it appears on command line
	Short      string // optional short name
	EnvName    string
	Usage      string // help message
	Value      Value  // value as set
	DefValue   string // default value (as text); for usage message
	Hidden     bool
	Deprecated bool
}
