package l

import (
	"testing"
)

func TestPurpose(t *testing.T) {
	for p := INFO; p <= FATAL; p++ {
		// Let it crash when an enum is not represented.
		_ = p.String()
	}
}
