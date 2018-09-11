package objects

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"k8s.io/helm/pkg/engine"
)

const (
	Helpers = "_helpers"
)

type Objects map[string]*bytes.Buffer

func (objects Objects) Names() []string {
	var names = make([]string, 0, len(objects))
	for name := range objects {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (objects Objects) Render(bootstrap *template.Template, output map[string]*bytes.Buffer, values map[string]interface{}) error {
	var tmpl, _ = bootstrap.Clone()
	tmpl = tmpl.Funcs(alterFuncMap(tmpl))
	if objects.Contains(Helpers) {
		var err error
		tmpl, err = tmpl.Parse(objects[Helpers].String())
		if err != nil {
			return emberr.ErrUnableToParseObject{
				Name:   Helpers,
				Reason: err,
			}
		}
	}

	for name, object := range objects {
		switch {
		case strings.HasPrefix(name, "_"):
			continue
		default:
			var vals = make(map[string]interface{}, len(values))
			for k, v := range values {
				vals[k] = v
			}
			vals["Template"] = map[string]interface{}{"Name": name, "BasePath": name}
			var objectTmpl, _ = tmpl.Clone()
			objectTmpl, parseObjectTmplErr := objectTmpl.Parse(object.String())
			if parseObjectTmplErr != nil {
				return emberr.ErrUnableToParseObject{
					Name:   name,
					Reason: parseObjectTmplErr,
				}
			}
			var buf = &bytes.Buffer{}
			var renderErr = objectTmpl.Execute(buf, values)
			if renderErr != nil {
				return emberr.ErrUnableToRenderObject{
					Name:   name,
					Reason: renderErr,
				}
			}
			output[name] = buf
		}
	}
	return nil
}

func (objects Objects) Contains(name string) bool {
	var _, exists = objects[name]
	return exists
}

func (objects Objects) Object(name string) (Object, error) {
	var buf, ok = objects[name]
	if !ok {
		return Object{}, emberr.ErrObjectNotFound{
			Name:              name,
			ObjectsWhichExist: objects.Names(),
		}
	}
	return Object{
		Name: name,
		Body: bytes.NewBufferString(buf.String()),
	}, nil
}

func (objects Objects) Slice() []Object {
	var slice = make([]Object, 0, len(objects))
	for name, body := range objects {
		slice = append(slice, Object{
			Name: name,
			Body: bytes.NewBufferString(body.String()),
		})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Name < slice[j].Name
	})
	return slice
}

type Object struct {
	Name string
	Body *bytes.Buffer
}

func (object Object) Render(tmpl *template.Template, output io.Writer, values map[string]interface{}) error {
	var err error
	tmpl, err = tmpl.Parse(object.Body.String())
	if err != nil {
		return err
	}
	return tmpl.Execute(output, values)
}

// The resulting FuncMap is only valid for the passed-in template.
func alterFuncMap(t *template.Template) template.FuncMap {
	// Clone the func map because we are adding context-specific functions.
	var funcMap template.FuncMap = map[string]interface{}{}
	for k, v := range engine.FuncMap() {
		funcMap[k] = v
	}

	// Add the 'include' function here so we can close over t.
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := t.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// Add the 'required' function here
	funcMap["required"] = func(warn string, val interface{}) (interface{}, error) {
		if val == nil {
			return val, fmt.Errorf(warn)
		} else if _, ok := val.(string); ok {
			if val == "" {
				return val, fmt.Errorf(warn)
			}
		}
		return val, nil
	}
	return funcMap
}
