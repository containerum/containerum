package emberr

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestLevenstein(test *testing.T) {
	var nearest = ErrObjectNotFound{
		Name:              "deployment",
		ObjectsWhichExist: []string{"svc", "deploy", "net"},
	}.findNearest()
	assert.Equal(test, nearest, "deploy")
}
