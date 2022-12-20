package version

import (
	"testing"
)

// arg 1 means expected value and actual stands for 'actual result of command'
type formatVersionTest struct {
	arg1, actual string
}

// test function
func TestFormatVersion(t *testing.T) {

	// declaring mock input values for function to be tested

	var version string
	var dateStr string

	version = "v1.1.0-1"
	dateStr = "2021-10-12"

	// array containing a truthy and falsy test case
	var formatVersionTests = []formatVersionTest{
		{"v1.1.0-1", Format(version, dateStr)},
		{"lr version 1.1.0-1 (2021-10-12)\nhttps://github.com/loginradius/lr-cli/releases/tag/v1.1.0-1\n", Format(version, dateStr)},
	}

	for _, test := range formatVersionTests {
		if output := test.arg1; output != test.actual {
			t.Errorf("Output %q not equal to actual value %q", output, test.actual)
		} else {
			t.Log("Success")
		}
	}
}
