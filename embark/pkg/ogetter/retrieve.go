package ogetter

import "bytes"

// Loads content of selected objects to map[string]*bytes.Buffer
func RetrieveObjects(getter ObjectGetter, objectNames ...string) (map[string]*bytes.Buffer, error) {
	var objects = make(map[string]*bytes.Buffer, len(objectNames))
	for _, objectName := range objectNames {
		var buf = &bytes.Buffer{}
		if err := getter.Object(objectName, buf); err != nil {
			return nil, err
		}
		objects[objectName] = buf
	}
	return objects, nil
}
