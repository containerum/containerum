package template

import (
	"bytes"
	"errors"
	"text/template"
)

// Template holds b0x and file template
type Template struct {
	template string

	name      string
	Variables interface{}
}

// Set the template to be used
// "files" or "file"
func (t *Template) Set(name string) error {
	t.name = name
	if name != "files" && name != "file" {
		return errors.New(`Error: Template must be "files" or "file"`)
	}

	if name == "files" {
		t.template = filesTemplate
	} else if name == "file" {
		t.template = fileTemplate
	}

	return nil
}

// Exec the template and return the final data as byte array
func (t *Template) Exec() ([]byte, error) {
	tmpl, err := template.New(t.name).Funcs(funcsTemplate).Parse(t.template)
	if err != nil {
		return nil, err
	}

	// exec template
	buff := bytes.NewBufferString("")
	err = tmpl.Execute(buff, t.Variables)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
