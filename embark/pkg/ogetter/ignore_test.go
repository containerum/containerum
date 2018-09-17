package ogetter

import "testing"

func TestIsIgnored(test *testing.T) {
	test.Logf("IsIgnore regexp: %q", ignoreRegex)
	type testCase struct {
		name       string
		isIngnored bool
	}
	var testCases = []testCase{
		{name: ".helmignore", isIngnored: true},
		{name: "Chart.yaml", isIngnored: false},
		{name: "requirements.lock", isIngnored: true},
	}

	for i, tec := range testCases {
		var got = IsIgnored(tec.name)
		if got != tec.isIngnored {
			test.Fatalf("test case %d %q: IsIgnored expected %v, got %v", i, tec.name, tec.isIngnored, got)
		}
	}
}
