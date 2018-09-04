package flags

import "testing"

func TestInstall_Validate(test *testing.T) {
	var testCases = []struct {
		name  string
		flags Install
		err   bool
	}{
		{name: "empty", flags: Install{}, err: true},
		{name: "non-empty dir", flags: Install{Dir: "~/.local"}, err: false},
	}
	for i, testCase := range testCases {
		var err = testCase.flags.Validate()
		switch {
		case (err == nil) && testCase.err:
			test.Fatalf("%2d %s expects error, got nothing", i, testCase.name)
		case (err != nil) && !testCase.err:
			test.Fatalf("%2d %s expects no errors, got %q", i, testCase.name, err)
		}
	}
}
