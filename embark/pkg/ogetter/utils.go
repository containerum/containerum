package ogetter

import (
	"regexp"
	"strings"
)

func Ingnored() []string {
	var patterns = []string{
		"\\.helmignore",
		`.+\.lock`,
	}
	return patterns
}

var ignoreRegex = regexp.MustCompile(strings.Join(Ingnored(), "|"))

func IgnoreRegex() *regexp.Regexp {
	return ignoreRegex.Copy()
}

func IsIgnored(fname string) bool {
	return ignoreRegex.MatchString(fname)
}
