package common

import "testing"

// TestFake is here to allow 'make test' to pass for our CI job
func TestFake(t *testing.T) {
	t.Log("fake test for CI to pass, since we have no tests in our project yet!")
	t.Log("remove after we have our first real test!")
}
