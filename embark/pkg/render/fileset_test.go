package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFilesetFromDir(test *testing.T) {
	type testCase struct {
		dir         string
		expectErr   bool
		expectNames []string
	}
	var testCases = []testCase{
		{dir: "testdata", expectNames: []string{"svc.yaml", "svc", "deployment.yaml", "deployment", "deployment.yml", "service.yml", "service"}},
		{dir: "testdata/svc.yaml", expectErr: true},
	}
	for i, t := range testCases {
		var set, err = FileSetFromDir(t.dir)
		switch {
		case (err != nil) && !t.expectErr:
			test.Fatalf("%2d error is not expected, got %q", i, err)
		case (err == nil) && t.expectErr:
			test.Fatalf("%2d expect error, got nothing", i)
		case (err == nil) && !t.expectErr:
			assert.ElementsMatch(test, set.Names(), t.expectNames)
		}
	}
}
