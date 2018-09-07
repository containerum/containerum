package sflags

import (
	"strings"
)

// transform s from CamelCase to flag-case
func camelToFlag(s, flagDivider string) string {
	splitted := split(s)
	return strings.ToLower(strings.Join(splitted, flagDivider))
}

// transform s from flag-case to CAMEL_CASE
func flagToEnv(s, flagDivider, envDivider string) string {
	return strings.ToUpper(strings.Replace(s, flagDivider, envDivider, -1))
}
