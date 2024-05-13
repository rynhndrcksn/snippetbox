package assert

import "testing"

func Equal[T comparable](t *testing.T, got, want T) {
	// Tells Go's test runner that the Equal() function is a test helper.
	// This means that if t.Errorf() is called the Go test runner will report
	// the filename and line number of the code that called Equal() in the output.
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
