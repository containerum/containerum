package builder

import (
	"path"
	"strings"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"gopkg.in/yaml.v2"
	kubeApsV1 "k8s.io/api/apps/v1"
	kubeBatchAPIv1 "k8s.io/api/batch/v1"
	kubeCoreV1 "k8s.io/api/core/v1"
	kubeExtensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const notesFileSuffix = "NOTES.txt"

type renderOptions struct {
	Values map[string]interface{}
}

func (options renderOptions) Merge(another ...renderOptions) renderOptions {
	for _, anotherOption := range another {
		if anotherOption.Values != nil {
			options.Values = anotherOption.Values
		}
	}
	return options
}

func RenderWithValues(values map[string]interface{}) renderOptions {
	return renderOptions{
		Values: values,
	}
}

type RenderedChart struct {
	Deployments []kubeApsV1.Deployment
	Ingresses   []kubeExtensionsV1beta1.Ingress
	Services    []kubeCoreV1.Service
	Secrets     []kubeCoreV1.Secret
	Configs     []kubeCoreV1.ConfigMap
	Volumes     []kubeCoreV1.Volume
	Jobs        []kubeBatchAPIv1.Job
}

func RenderChart(ch *chart.Chart, options ...renderOptions) (*RenderedChart, error) {
	//var chartValues, readChartValuesErr = chartutil.ReadValues([]byte(ch.GetValues().Raw))
	//if readChartValuesErr != nil {
	//	return readChartValuesErr
	//}

	var renderConfig = renderOptions{}
	var chartValuesCapsErr error
	renderConfig.Values, chartValuesCapsErr = chartutil.ToRenderValuesCaps(ch,
		&chart.Config{
			Raw: ch.GetValues().Raw,
			//	Values: chartValues,
		},
		chartutil.ReleaseOptions{
			Name:      Containerum,
			Namespace: Containerum,
			IsInstall: true,
		},
		&chartutil.Capabilities{})
	if chartValuesCapsErr != nil {
		return nil, chartValuesCapsErr
	}

	for _, option := range options {
		if option.Values != nil {
			renderConfig.Values = option.Values
		}
	}

	var renderEngine = engine.New()
	var targets, renderErr = renderEngine.Render(ch, renderConfig.Values)
	if renderErr != nil {
		return nil, renderErr
	}

	var notes = make([]string, 0)
	var rendered = RenderedChart{}
	for k, v := range targets {
		if strings.HasSuffix(k, notesFileSuffix) {
			// Only apply the notes if it belongs to the parent ch
			// Note: Do not use filePath.Join since it creates a path with \ which is not expected
			if k == path.Join(ch.Metadata.Name, "templates", notesFileSuffix) {
				notes = append(notes, v)
			}
			delete(targets, k)
		} else {
			var meta v1.TypeMeta
			var metaUnmarshalErr = yaml.Unmarshal([]byte(v), &meta)
			if metaUnmarshalErr != nil {
				return nil, nil
			}
			switch strings.ToLower(meta.Kind) {
			case "deployment":
				var depl, err = parseDeployment(serializedKubeObject)
				if err != nil {
					return nil, err
				}
				rendered.Deployments = append(rendered.Deployments, depl)
			case "service":
				var serv, err = parseService(v)
				if err != nil {
					return nil, err
				}
				rendered.Services = append(rendered.Services, serv)
			case "volume":
				var volume, err = parseVolume(v)
				if err != nil {
					return nil, err
				}
				rendered.Volumes = append(rendered.Volumes, volume)
			case "configmap":
				var configmap, err = parseConfigmap(v)
				if err != nil {
					return nil, err
				}
				rendered.Configs = append(rendered.Configs, configmap)
			case "job":
				var job, err = parseJob(v)
				if err != nil {
					return nil, err
				}
				rendered.Jobs = append(rendered.Jobs, job)
			case "secret":
				var secret, err = parseSecret(v)
				if err != nil {
					return nil, err
				}
				rendered.Secrets = append(rendered.Secrets, secret)
			case "ingress":
				var ingr, err = parseIngress(v)
				if err != nil {
					return nil, err
				}
				rendered.Ingresses = append(rendered.Ingresses, ingr)
			default:
				return nil, emberr.ErrUnsupportedKubeObjectType(meta.Kind)
			}
		}
	}
	return &rendered, nil
}

func parseDeployment(data string) (kubeApsV1.Deployment, error) {
	var obj = kubeApsV1.Deployment{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseService(data string) (kubeCoreV1.Service, error) {
	var obj = kubeCoreV1.Service{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseVolume(data string) (kubeCoreV1.Volume, error) {
	var obj = kubeCoreV1.Volume{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseIngress(data string) (kubeExtensionsV1beta1.Ingress, error) {
	var obj = kubeExtensionsV1beta1.Ingress{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseSecret(data string) (kubeCoreV1.Secret, error) {
	var obj = kubeCoreV1.Secret{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseJob(data string) (kubeBatchAPIv1.Job, error) {
	var obj = kubeBatchAPIv1.Job{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseConfigmap(data string) (kubeCoreV1.ConfigMap, error) {
	var obj = kubeCoreV1.ConfigMap{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}
