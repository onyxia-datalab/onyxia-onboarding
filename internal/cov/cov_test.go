package cov_test

import "testing"
import "github.com/onyxia-datalab/onyxia-onboarding/internal/cov"

func TestTested(t *testing.T) {
	if res := cov.Tested(); res != 2 {
		t.Fatalf("Tested() = %q, expected %q", res, 2)
	}
}
