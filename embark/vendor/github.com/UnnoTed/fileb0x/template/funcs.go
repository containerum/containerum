package template

import (
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/UnnoTed/fileb0x/config"
)

var safeNameBlacklist = map[string]string{}
var blacklist = map[string]int{}

// taken from golint @ https://github.com/golang/lint/blob/master/lint.go#L702
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

var r = regexp.MustCompile(`[^a-zA-Z0-9]`)

var funcsTemplate = template.FuncMap{
	"exported":         exported,
	"buildTags":        buildTags,
	"exportedTitle":    exportedTitle,
	"buildSafeVarName": buildSafeVarName,
}

var unexported bool

// SetUnexported variables, functions and types
func SetUnexported(e bool) {
	unexported = e
}

func exported(field string) string {
	if !unexported {
		return strings.ToUpper(field)
	}

	return strings.ToLower(field)
}

func exportedTitle(field string) string {
	if !unexported {
		return strings.Title(field)
	}

	return strings.ToLower(field[0:1]) + field[1:]
}

func buildSafeVarName(path string) string {
	name, exists := safeNameBlacklist[path]
	if exists {
		return name
	}

	n := config.SafeVarName.ReplaceAllString(path, "$")
	words := strings.Split(n, "$")

	name = ""
	// check for uppercase words
	for _, word := range words {
		upper := strings.ToUpper(word)

		if commonInitialisms[upper] {
			name += upper
		} else {
			name += strings.Title(word)
		}
	}

	// avoid redeclaring variables
	//
	// _file.txt
	// file.txt
	_, blacklisted := blacklist[name]

	if blacklisted {
		blacklist[name]++
		name += strconv.Itoa(blacklist[name])
	}

	safeNameBlacklist[path] = name
	blacklist[name]++
	return name
}

func buildTags(tags string) string {
	if tags != "" {
		tags = "// +build " + tags + "\n"
	}

	return tags
}
