package emberr

import "testing"

func TestLevenstein(test *testing.T) {
	test.Log(ErrObjectNotFound{
		Name:              "deployment",
		ObjectsWhichExist: []string{"svc", "deploy", "net"},
	}.Error())
}
