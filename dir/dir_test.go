package dir

import (
	"testing"
)

func TestList(t *testing.T) {
	// Assume that / can always be tested. Other than that we don't have many tests to run.
	if _, err := List("/"); err != nil {
		t.Errorf("List(/) = _, %v, need nil error", err)
	}
}
