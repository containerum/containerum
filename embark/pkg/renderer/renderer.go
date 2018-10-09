package renderer

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/containerum/containerum/embark/pkg/emberr"

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

type ObjectConstructor = func(io.Reader) (kube.Object, error)

type Renderer struct {
	Name            string
	ObjectsToRender []string
	DependsOn       []string
	ObjectGetter    ogetter.ObjectGetter
	Constructor     ObjectConstructor
	Values          Values
}

func (renderer Renderer) RenderComponent() (RenderedComponent, error) {
	var name = renderer.Name
	var getter = renderer.ObjectGetter
	var constructor = renderer.Constructor
	var values = renderer.Values
	var templ = template.New(Template).Funcs(engine.FuncMap())
	var names = renderer.ObjectsToRender
	var objectExists = CheckIfObjectExists(renderer.ObjectGetter.ObjectNames())
	var null RenderedComponent

	var buf = &bytes.Buffer{}
	if objectExists(Helpers) {
		buf.Reset()
		var helpersTextBuf = buf
		if err := getter.Object(Helpers, helpersTextBuf); err != nil {
			return null, err
		}
		var parseHelperTemplErr error
		templ, parseHelperTemplErr = templ.New(Helpers).
			Funcs(engine.FuncMap()).
			Parse(helpersTextBuf.String())
		if parseHelperTemplErr != nil {
			return null, emberr.ErrUnableToRenderObject{
				Name:   name,
				Reason: parseHelperTemplErr,
			}
		}
	}
	var objects = make([]kube.Object, 0, len(names))
	for _, name := range names {
		switch name {
		case Helpers, Notes:
			continue
		default:
			var objectTextBuf = buf
			objectTextBuf.Reset()
			if err := getter.Object(name, objectTextBuf); err != nil {
				return null, emberr.ErrUnableToRenderObject{
					Name:   name,
					Reason: err,
				}
			}
			var parseObjectTemplateErr error
			templ, parseObjectTemplateErr = templ.
				New(name).
				Parse(objectTextBuf.String())
			if parseObjectTemplateErr != nil {
				return null, emberr.ErrUnableToRenderObject{
					Name:   name,
					Reason: parseObjectTemplateErr,
				}
			}
		}
	}
	for _, name := range names {
		switch {
		case strings.HasPrefix(name, "_"), name == Notes:
			continue
		}
		var objectTextBuf = buf
		objectTextBuf.Reset()
		if err := templ.ExecuteTemplate(objectTextBuf, name, values); err != nil {
			return null, emberr.ErrUnableToRenderObject{
				Name:   name,
				Reason: err,
			}
		}
		var object, createObjectErr = constructor(objectTextBuf)
		if createObjectErr != nil {
			return null, emberr.ErrUnableToRenderObject{
				Name:   name,
				Reason: createObjectErr,
			}
		}
		objects = append(objects, object)
	}
	return NewRenderedObject(name, renderer.DependsOn, objects...), nil
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
