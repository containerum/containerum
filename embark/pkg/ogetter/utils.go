package ogetter

import (
	"path"
	"regexp"
	"strings"
)

func Ignored() []string {
	return []string{
		`\.helmignore`,
		`.+\.lock`,
	}
}

var ignoreRegex = regexp.MustCompile(strings.Join(Ignored(), "|"))

func IgnoreRegex() *regexp.Regexp {
	return ignoreRegex.Copy()
}

func IsIgnored(fname string) bool {
	return ignoreRegex.MatchString(fname)
}

func extractObjectNameFromFilename(fname string) string {
	var ext = path.Ext(fname)
	var _, objectName = path.Split(fname)
	objectName = strings.TrimSuffix(objectName, ext)
	return objectName
}
