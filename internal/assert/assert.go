package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	// Tells Go's test runner that the Equal() function is a test helper.
	// This means that if t.Errorf() is called the Go test runner will report
	// the filename and line number of the code that called Equal() in the output.
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
