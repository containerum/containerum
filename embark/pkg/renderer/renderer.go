package renderer

import (
	"bytes"
	"html/template"
	"io"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/ogetter"
)

const (
	Values   = "values"
	Helpers  = "_helpers"
	Template = "object"
)

type ObjectConstructor func(io.Reader) (kube.Object, error)

type Renderer struct {
	Name         string
	Getter       ogetter.ObjectGetter
	Contstructor ObjectConstructor
	Values       map[string]interface{}
}

func (renderer Renderer) RenderComponent() (RenderedComponent, error) {
	var name = renderer.Name
	var getter = renderer.Getter
	var constructor = renderer.Contstructor
	var values = renderer.Values
	var templ = template.New(Template)
	var names = getter.ObjectNames()
	var objectExists = CheckIfObjectExists(names)
	var null RenderedComponent

	if objectExists(Helpers) {
		var helpersTextBuf = &bytes.Buffer{}
		if err := getter.Object(Helpers, helpersTextBuf); err != nil {
			return null, err
		}
		var parseHelperTemplErr error
		templ, parseHelperTemplErr = templ.Parse(helpersTextBuf.String())
		if parseHelperTemplErr != nil {
			return null, parseHelperTemplErr
		}
	}

	var objects = make([]kube.Object, 0, len(names))
	for _, name := range names {
		switch name {
		case Helpers:
			continue
		default:
			var objectTemplate, cloneTemplErr = templ.Clone()
			if cloneTemplErr != nil {
				panic(cloneTemplErr) // something really bad happened!
			}
			var objectTextBuf = &bytes.Buffer{}
			if err := getter.Object(name, objectTextBuf); err != nil {
				return null, err
			}
			var parseObjectTemplateErr error
			objectTemplate, parseObjectTemplateErr = objectTemplate.Parse(objectTextBuf.String())
			if parseObjectTemplateErr != nil {
				return null, parseObjectTemplateErr
			}
			objectTextBuf.Reset()
			if err := objectTemplate.Execute(objectTextBuf, values); err != nil {
				return null, err
			}
			var object, createObjectErr = constructor(objectTextBuf)
			if createObjectErr != nil {
				return null, createObjectErr
			}
			objects = append(objects, object)
		}
	}
	return NewRenderedObject(name, objects...), nil
}

func CheckIfObjectExists(names []string) func(name string) bool {
	var set = make(map[string]struct{}, len(names))
	var void = struct{}{}
	for _, name := range names {
		set[name] = void
	}
	return func(name string) bool {
		var _, ok = set[name]
		return ok
	}
}
