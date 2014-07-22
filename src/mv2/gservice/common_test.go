package gservice

import (
	"fmt"
	"testing"
)

func TestIsValidClientKey(t *testing.T) {
	for _, v := range clientKeys {
		if !IsValidClientKey(v) {
			t.Error(fmt.Sprintf("%s should be a valid client key", v))
		}
	}
}

func TestIsValidClientKeyInvalidCase(t *testing.T) {
	fKey := clientKeys[0] + "x"
	if IsValidClientKey(fKey) {
		t.Error(fmt.Sprintf("%s should not be a valid client key", fKey))
	}
}
