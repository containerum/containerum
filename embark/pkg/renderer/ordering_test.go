package renderer

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/containerum/containerum/embark/pkg/utils/why"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/magiconair/properties/assert"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func TestSortObject(test *testing.T) {
	var kinds = append([]string{
		"anything",
		"weird object",
	}, DefaultOrder()...)
	var objects = mocks(kinds)
	SortObjects(objects)
	var sorted = extractKinds(objects)
	test.Log(sorted)
	assert.Equal(test, sorted, kinds)
}

func TestObjectsToBatches(test *testing.T) {
	var kinds = append([]string{
		"anything",
		"weird object",
	}, DefaultOrder()...)
	var objects = mocks(kinds)
	var batches = ObjectsToBatches(objects)
	var paragon = [][]string{
		{"anything", "weird object"},
		{"configmap"},
		{"volume"},
		{"deployment"},
		{"service"},
		{"ingress"},
	}
	why.PrintFromIter("batches", len(batches), func(i int) (string, error) {
		var batch = batches[i]
		var kinds = extractKinds(batch)
		return strings.Join(kinds, ", "), nil
	})
	for i, batch := range paragon {
		assert.Equal(test, extractKinds(batches[i]), batch)
	}
}

func mock(kind string) kube.ObjectMock {
	return kube.ObjectMock{ObjectKind: kind}
}

func mocks(kinds []string) []kube.Object {
	kinds = append([]string{}, kinds...)
	rnd.Shuffle(len(kinds), func(i, j int) {
		kinds[i], kinds[j] = kinds[j], kinds[i]
	})
	var objects = make([]kube.Object, 0, len(kinds))
	for _, kind := range kinds {
		objects = append(objects, mock(kind))
	}
	return objects
}

func extractKinds(objects []kube.Object) []string {
	var kinds = make([]string, 0, len(objects))
	for _, object := range objects {
		kinds = append(kinds, object.Kind())
	}
	return kinds
}
