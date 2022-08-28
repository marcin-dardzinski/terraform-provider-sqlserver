package internal

import (
	"testing"
)

func TestProviderSchema(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
