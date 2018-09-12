package kube

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ericchiang/k8s"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/thoas/go-funk"
)

var (
	_ json.Marshaler   = Object{}
	_ json.Unmarshaler = new(Object)

	//	_ yaml.Marshaler   = Object{}
	//	_ yaml.Unmarshaler = new(Object)
	_ k8s.Resource = Object{}
)

func init() {
	k8s.Register("embark", "v1", "embark-generic-object", true, new(Object))
}

type Object struct {
	meta *metav1.ObjectMeta
	body map[string]interface{}
}

func (object Object) Query(path string) interface{} {
	return funk.Get(object.body, path)
}

func (object Object) String() string {
	var data, err = yaml.Marshal(object)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (object Object) ToStruct(target interface{}) error {
	return mapstructure.WeakDecode(object.body, target)
}

func ObjectFromJSON(re io.Reader) (Object, error) {
	var object = Object{
		body: make(map[string]interface{}),
	}
	return object, json.NewDecoder(re).Decode(object.body)
}

func ObjectFromYAML(re io.Reader) (Object, error) {
	var object = Object{
		body: make(map[string]interface{}),
	}
	var buf = &bytes.Buffer{}
	if _, err := buf.ReadFrom(re); err != nil {
		return Object{}, err
	}
	if err := yaml.Unmarshal(buf.Bytes(), object.body); err != nil {
		return Object{}, err
	}
	return object, nil
}

func (object *Object) PatchMeta(patchers ...func(meta *metav1.ObjectMeta)) {
	object.initMeta()
	for _, patcher := range patchers {
		patcher(object.meta)
	}
}

func ObjectFromStruct(data interface{}) (Object, error) {
	var object Object
	return object, mapstructure.Decode(data, &object.body)
}

func (object *Object) initMeta() *metav1.ObjectMeta {
	if object.meta == nil {
		object.meta = new(metav1.ObjectMeta)
	}
	if err := mapstructure.WeakDecode(object.body["metadata"], object.meta); err != nil {
		panic(err)
	}
	return object.meta
}

func (object Object) GetMetadata() *metav1.ObjectMeta {
	if object.meta == nil {
		object.initMeta()
	}
	return object.meta
}

func (object Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(object.body)
}

func (object *Object) UnmarshalJSON(data []byte) error {
	return yaml.Unmarshal(data, &object.body)
}

func (object Object) MarshalYAML() (interface{}, error) {
	return object.body, nil
}

func (object *Object) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	return unmarshaler(&object.body)
}
