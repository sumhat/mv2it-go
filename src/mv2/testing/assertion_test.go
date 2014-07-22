package testing

import (
	"testing"
)

func TestAssertEqualsIntValid(t *testing.T) {
	Assert(t).That(1).IsEqualTo(1)
	Assert(t).That(-1).IsEqualTo(-1)
}

func TestAssertLessThanValid(t *testing.T) {
	Assert(t).That(1).IsLessThan(2)
	Assert(t).That(1.5).IsLessThan(3)
}

func TestAssertLessThanInvalid(t *testing.T) {
	defer func() {
		_ = recover().(*UnsupportedTypeError)
	}()
	Assert(t).That(1).IsLessThan("abcd")
}

func TestThatRune(t *testing.T) {
	Assert(t).ThatRune(rune('a')).IsEqualTo(rune('a'))
}
