package gservice

import (
	T "mv2/testing"
	"testing"
)

func TestIsValidClientKey(t *testing.T) {
	for _, v := range clientKeys {
		T.Assert(t).That(IsValidClientKey(v)).IsTrue()
	}
}

func TestIsValidClientKeyInvalidCase(t *testing.T) {
	fKey := clientKeys[0] + "x"
	T.Assert(t).That(IsValidClientKey(fKey)).IsFalse()
}
