package renderer

import (
	"sort"

	"github.com/containerum/containerum/embark/pkg/kube"
)

var objectPriorities = map[string]int{}

func DefaultOrder() []string {
	return []string{
		"configmap",
		"volume",
		"deployment",
		"service",
		"ingress",
	}
}

func init() {
	var order = DefaultOrder()
	for i, kind := range order {
		objectPriorities[kind] = i + 1
	}
}

func ObjectPriority(kind string) int {
	return objectPriorities[kind]
}

func SortObjects(objects []kube.Object) {
	sort.Slice(objects, func(i, j int) bool {
		var a = ObjectPriority(objects[i].Kind())
		var b = ObjectPriority(objects[j].Kind())
		if a == b {
			return objects[i].Kind() < objects[j].Kind()
		}
		return a < b
	})
}

func ObjectsToBatches(objects []kube.Object) [][]kube.Object {
	var batches [][]kube.Object
	objects = append([]kube.Object{}, objects...)
	if len(objects) < 2 {
		return append(batches, objects)
	}
	SortObjects(objects)
	var batch []kube.Object
	var prevPriority = ObjectPriority(objects[0].Kind())
	for _, object := range objects {
		var priority = ObjectPriority(object.Kind())
		if priority != prevPriority {
			var preallocSize = len(batch)
			batches = append(batches, batch)
			batch = make([]kube.Object, 0, preallocSize)
			prevPriority = priority
		}
		batch = append(batch, object)
	}
	batches = append(batches, batch)
	return batches
}
