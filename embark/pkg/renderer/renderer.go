package renderer

import (
	"bytes"
	"io"
	"sort"
	"text/template"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/ogetter"
	"k8s.io/helm/pkg/engine"
)

const (
	ValuesName = "values"
	Helpers    = "_helpers"
	Template   = "object"
	Notes      = "NOTES"
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
	var templ = template.New(Template).Funcs(engine.FuncMap())
	var names = getter.ObjectNames()
	var objectExists = CheckIfObjectExists(names)
	var null RenderedComponent

	var buf = &bytes.Buffer{}
	if objectExists(Helpers) {
		var helpersTextBuf = buf
		buf.Reset()
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
	sort.Strings(names)
	for _, name := range names {
		switch name {
		case Helpers, Notes:
			continue
		default:
			var objectTemplate, cloneTemplErr = templ.Clone()
			if cloneTemplErr != nil {
				panic(cloneTemplErr) // something really bad happened!
			}
			var objectTextBuf = buf
			buf.Reset()
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
