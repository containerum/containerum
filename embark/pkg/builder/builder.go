package builder

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/helm"
)

type Builder struct {
	Template string
	Values   string
	Helpers  string
	Output   io.Writer
}

func (builder Builder) Build() error {
	helm.NewClient()
	var targetTemplate, err = template.
		New("target").
		Funcs(template.FuncMap{
			"toYaml":     ToYaml,
			"default":    Default,
			"indent":     Indent,
			"trunc":      Truncate,
			"trimSuffix": strings.TrimSuffix,
			"contains":   strings.Contains,
			"replace":    strings.Replace,
		}).Parse(builder.Helpers)
	if err != nil {
		return fmt.Errorf("unable to parse helpers: %v", err)
	}
	targetTemplate, err = targetTemplate.Parse(builder.Template)
	if err != nil {
		return fmt.Errorf("unable to parse target template: %v", err)
	}
	var values = make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(builder.Values), &values); err != nil {
		return err
	}
	type data struct {
		Release struct {
			Name string
		}
		Values map[string]interface{}
	}
	return targetTemplate.ExecuteTemplate(builder.Output, "target", data{
		Values: values,
	})
}

func ToYaml(value interface{}) string {
	var data, err = yaml.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Indent(text string, width int) string {
	var lines = strings.Split(text, "\n")
	return strings.Join(lines, "\n"+strings.Repeat(" ", width))
}

func Default(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

func Truncate(text string, width int) string {
	var rstr = []rune(text)
	if len(rstr) < 63 {
		return string(rstr[:width])
	}
	return text
}
