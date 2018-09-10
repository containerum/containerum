package objects

import (
	"bytes"
	"text/template"
)

type Objects map[string]*bytes.Buffer

func (objects Objects) Names() []string {
	var names = make([]string, 0, len(objects))
	for name := range objects {
		names = append(names, name)
	}
	return names
}

func (objects Objects) Render(bootstrap *template.Template, output map[string]*bytes.Buffer, values map[string]interface{}) error {
	var tmpl, _ = bootstrap.Clone()
	for name, object := range objects {
		var objectTmpl, _ = tmpl.Clone()
		objectTmpl, parseObjectTmplErr := objectTmpl.Parse(object.String())
		if parseObjectTmplErr != nil {
			return parseObjectTmplErr
		}
		var buf = &bytes.Buffer{}
		var renderErr = objectTmpl.Execute(buf, values)
		if renderErr != nil {
			return renderErr
		}
		output[name] = buf
	}
	return nil
}

func (objects Objects) Contains(name string) bool {
	var _, exists = objects[name]
	return exists
}
