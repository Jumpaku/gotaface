package assert

import (
	"testing"
)

func anyEqual(actual any, expect any) bool {
	return actual == expect
}
func Equal[T comparable](t *testing.T, actual T, expect T) {
	t.Helper()

	if !anyEqual(actual, expect) {
		t.Errorf("ASSERT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}
