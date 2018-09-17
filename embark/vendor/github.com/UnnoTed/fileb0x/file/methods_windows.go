package file

import "strings"

// GetRemap returns a map's params with
// info required to load files directly
// from the hard drive when using prefix
// and base while debug mode is activaTed
func (f *File) GetRemap() string {
	if f.Base == "" && f.Prefix == "" {
		return ""
	}

	return `"` + strings.Replace(f.OriginalPath, `\`, `\\`, -1) + `": {
		"prefix": "` + f.Prefix + `",
		"base": "` + f.Base + `",
	},`
}
