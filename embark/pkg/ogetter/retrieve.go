package ogetter

import (
	"bytes"

	"github.com/containerum/containerum/embark/pkg/emberr"
)

func extras() []string {
	return []string{
		"_helpers",
	}
}

var isExtras = make(map[string]bool)

// Loads content of selected objects to map[string]*bytes.Buffer
func RetrieveObjects(getter ObjectGetter, objectNames ...string) (Objects, error) {
	objectNames = append(objectNames, extras()...)
	var objects = make(map[string]*bytes.Buffer, len(objectNames))
	for _, objectName := range objectNames {
		var buf = &bytes.Buffer{}
		switch err := getter.Object(objectName, buf).(type) {
		case emberr.ErrObjectNotFound:
			if !isExtras[objectName] {
				return nil, err
			}
		default:
			objects[objectName] = buf
			continue
		}
	}
	return objects, nil
}

func init() {
	for _, item := range extras() {
		isExtras[item] = true
	}
}
