package main

import (
	"testing"
)

type versionTest struct {
	arg1, expected string
}

var versionTests = []versionTest{
	{"v1.1.0-1-gd37a9ba", version()},
	{"10.0.0", version()},
}

func TestVersion(t *testing.T) {
	for _, test := range versionTests {
		if output := test.arg1; output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		} else {
			t.Log("Success")
		}
	}
}
