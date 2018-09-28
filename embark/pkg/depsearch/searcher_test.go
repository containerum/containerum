package depsearch

import (
	"testing"
)

func TestStatic(test *testing.T) {
	var searcher = Static()
	type testCase struct {
		Component string
		MustExist bool
	}
	var testCases = []testCase{
		{Component: "auth", MustExist: true},
		{Component: "alkf", MustExist: false},
		{Component: "postgresql", MustExist: true},
	}
	for i, tc := range testCases {
		var componentExists = searcher.Contains(tc.Component)
		if componentExists != tc.MustExist {
			var mustExist = "must not exist"
			if tc.MustExist {
				mustExist = "must exist"
			}
			test.Fatalf("test case %d: component %q %s", i, tc.Component, mustExist)
		}
	}
}
