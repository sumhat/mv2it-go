package testing

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

const (
	Rune = iota + 100
	Nil
)

type AssertT struct {
	t *testing.T
}

func Assert(t *testing.T) *AssertT {
	return &AssertT{t: t}
}

type ThatInterface struct {
	t            *testing.T
	name string
	boolValue    bool
	runeValue    rune
	int32Value   int32
	int64Value   int64
	uint32Value  uint32
	uint64Value  uint64
	float64Value float64
	stringValue  string
	origin       interface{}
	kind         reflect.Kind
}

type UnsupportedTypeError struct {
	Type  reflect.Type
	Value interface{}
}

func newUnsupportedTypeError(v interface{}) *UnsupportedTypeError {
	return &UnsupportedTypeError{Type: reflect.TypeOf(v), Value: v}
}

func (this *UnsupportedTypeError) String() string {
	return fmt.Sprintf("%v: %s", this.Value, this.Type.String())
}

func toInt32(v interface{}) int32 {
	switch v.(type) {
	case int:
		return int32(v.(int))
	case int8:
		return int32(v.(int8))
	case int16:
		return int32(v.(int16))
	case int32:
		return v.(int32)
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func toUInt32(v interface{}) uint32 {
	switch v.(type) {
	case uint:
		return uint32(v.(uint))
	case uint8:
		return uint32(v.(uint8))
	case uint16:
		return uint32(v.(uint16))
	case uint32:
		return v.(uint32)
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func toFloat64(v interface{}) float64 {
	switch v.(type) {
	case uint, uint8, uint16, uint32:
		return float64(toUInt32(v))
	case uint64:
		return float64(v.(uint64))
	case int, int8, int16, int32:
		return float64(toInt32(v))
	case int64:
		return float64(v.(int64))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func toString(v interface{}) string {
	switch v.(type) {
	case rune:
		return string(v.(rune))
	case string:
		return v.(string)
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func (this *AssertT) ThatRune(v rune) *ThatInterface {
	return &ThatInterface{
		t:         this.t,
		origin:    v,
		kind:      Rune,
		runeValue: v,
	}
}

func (this *AssertT) That(v interface{}) *ThatInterface {
	that := &ThatInterface{t: this.t, origin: v}
	switch v.(type) {
	case bool:
		that.boolValue = v.(bool)
		that.kind = reflect.Bool
	case int, int8, int16, int32:
		that.int32Value = toInt32(v)
		that.kind = reflect.Int32
	case int64:
		that.int64Value = v.(int64)
		that.kind = reflect.Int64
	case uint, uint8, uint16, uint32:
		that.uint32Value = toUInt32(v)
		that.kind = reflect.Uint32
	case uint64:
		that.uint64Value = v.(uint64)
		that.kind = reflect.Uint64
	case float32, float64:
		that.float64Value = toFloat64(v)
		that.kind = reflect.Float64
	case string:
		that.stringValue = toString(v)
		that.kind = reflect.String
	default:
		if v == nil {
			that.kind = Nil
		} else {
			that.kind = reflect.TypeOf(v).Kind()
		}
	}
	return that
}

func (this *ThatInterface) Named(name string) *ThatInterface {
	this.name = name
	return this
}

func (this *ThatInterface) prependName(format string) string {
	if len(this.name) > 0 {
		format = this.name + ": " + format
	}
	return format
}

func (this *ThatInterface) IsEqualTo(v interface{}) {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not equal to %s"))
	switch this.kind {
	case reflect.Bool:
		if this.boolValue != v.(bool) {
			aError.raise2Bools(this.boolValue, v.(bool))
		}
	case Rune:
		if this.runeValue != v.(rune) {
			aError.raise2Runes(this.runeValue, v.(rune))
		}
	case reflect.Int32:
		int32Value := toInt32(v)
		if this.int32Value != int32Value {
			aError.raise2Int32s(this.int32Value, int32Value)
		}
	case reflect.Int64:
		if this.int64Value != v.(int64) {
			aError.raise2Int64s(this.int64Value, v.(int64))
		}
	case reflect.Uint32:
		if this.uint32Value != toUInt32(v) {
			aError.raise2UInt32s(this.uint32Value, toUInt32(v))
		}
	case reflect.Uint64:
		if this.uint64Value != v.(uint64) {
			aError.raise2UInt64s(this.uint64Value, v.(uint64))
		}
	case reflect.String:
		if this.stringValue != toString(v) {
			aError.raise2Strings(this.stringValue, v.(string))
		}
	case reflect.Interface:
		if !reflect.DeepEqual(this.origin, v) {
			aError.raise2Interfaces(this.origin, v)
		}
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func (this *ThatInterface) IsLessThan(v interface{}) {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not less than %s"))
	switch this.kind {
	case reflect.Int32:
		int32Value := toInt32(v)
		if !(this.int32Value < int32Value) {
			aError.raise2Int32s(this.int32Value, int32Value)
		}
	case reflect.Int64:
		if !(this.int64Value < v.(int64)) {
			aError.raise2Int64s(this.int64Value, v.(int64))
		}
	case reflect.Uint32:
		if !(this.uint32Value < toUInt32(v)) {
			aError.raise2UInt32s(this.uint32Value, toUInt32(v))
		}
	case reflect.Uint64:
		if !(this.uint64Value < v.(uint64)) {
			aError.raise2UInt64s(this.uint64Value, v.(uint64))
		}
	case reflect.Float64:
		if !(this.float64Value < toFloat64(v)) {
			aError.raise2Float64s(this.float64Value, toFloat64(v))
		}
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func (this *ThatInterface) IsGreaterThan(v interface{}) {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not greater than %s"))
	switch this.kind {
	case reflect.Int32:
		int32Value := toInt32(v)
		if !(this.int32Value > int32Value) {
			aError.raise2Int32s(this.int32Value, int32Value)
		}
	case reflect.Int64:
		if !(this.int64Value > v.(int64)) {
			aError.raise2Int64s(this.int64Value, v.(int64))
		}
	case reflect.Uint32:
		if !(this.uint32Value > toUInt32(v)) {
			aError.raise2UInt32s(this.uint32Value, toUInt32(v))
		}
	case reflect.Uint64:
		if !(this.uint64Value > v.(uint64)) {
			aError.raise2UInt64s(this.uint64Value, v.(uint64))
		}
	case reflect.Float64:
		if !(this.float64Value > toFloat64(v)) {
			aError.raise2Float64s(this.float64Value, toFloat64(v))
		}
	default:
		panic(newUnsupportedTypeError(v))
	}
}

func (this *ThatInterface) IsTrue() {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not true"))
	switch this.kind {
	case reflect.Bool:
		if this.boolValue != true {
			aError.raiseBool(this.boolValue)
		}
	default:
		panic(newUnsupportedTypeError(this.origin))
	}
}

func (this *ThatInterface) IsFalse() {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not false"))
	switch this.kind {
	case reflect.Bool:
		if this.boolValue != false {
			aError.raiseBool(this.boolValue)
		}
	default:
		panic(newUnsupportedTypeError(this.origin))
	}
}

func (this *ThatInterface) IsNil() {
	aError := newAssertError(this.t, this.prependName("%s (%s) is not nil"))
	if this.kind != Nil {
		aError.raiseInterface(this.origin)
	}
}

func (this *ThatInterface) IsNotNil() {
	aError := newAssertError(this.t, this.prependName("%s (%s) is nil"))
	if this.kind == Nil {
		aError.raiseInterface(this.origin)
	}
}

type AssertionError struct {
	t      *testing.T
	format string
}

func newAssertError(t *testing.T, format string) *AssertionError {
	return &AssertionError{t: t, format: format}
}

func (this *AssertionError) raise2TypedStrings(s1 string, s2 string, t string) {
	this.t.Errorf(this.format, s1, t, s2)
}

func (this *AssertionError) raiseTypedString(s1 string, t string) {
	this.t.Errorf(this.format, s1, t)
}

func (this *AssertionError) raise2Strings(s1 string, s2 string) {
	this.raise2TypedStrings(strconv.Quote(s1), strconv.Quote(s2), "string")
}

func (this *AssertionError) raise2Bools(b1 bool, b2 bool) {
	this.raise2TypedStrings(strconv.FormatBool(b1), strconv.FormatBool(b2), "bool")
}

func (this *AssertionError) raise2Int32s(i1 int32, i2 int32) {
	this.raise2TypedStrings(strconv.FormatInt(int64(i1), 10), strconv.FormatInt(int64(i2), 10), "int32")
}

func (this *AssertionError) raise2Int64s(i1 int64, i2 int64) {
	this.raise2TypedStrings(strconv.FormatInt(i1, 10), strconv.FormatInt(i2, 10), "int64")
}

func (this *AssertionError) raise2UInt32s(i1 uint32, i2 uint32) {
	this.raise2TypedStrings(strconv.FormatUint(uint64(i1), 10), strconv.FormatUint(uint64(i2), 10), "uint32")
}

func (this *AssertionError) raise2UInt64s(i1 uint64, i2 uint64) {
	this.raise2TypedStrings(strconv.FormatUint(i1, 10), strconv.FormatUint(i2, 10), "uint64")
}

func (this *AssertionError) raise2Float64s(f1 float64, f2 float64) {
	this.raise2TypedStrings(strconv.FormatFloat(f1, 'f', -1, 64), strconv.FormatFloat(f2, 'f', -1, 64), "float64")
}

func (this *AssertionError) raise2Interfaces(i1 interface{}, i2 interface{}) {
	this.raise2TypedStrings(fmt.Sprintf("%v", i1), fmt.Sprintf("%v", i2), reflect.TypeOf(i1).String())
}

func (this *AssertionError) raise2Runes(r1 rune, r2 rune) {
	this.raise2TypedStrings(strconv.QuoteRune(r1), strconv.QuoteRune(r2), "rune")
}

func (this *AssertionError) raiseBool(b bool) {
	this.raiseTypedString(strconv.FormatBool(b), "bool")
}

func (this *AssertionError) raiseString(s string) {
	this.raiseTypedString(strconv.Quote(s), "string")
}

func (this *AssertionError) raiseInterface(v interface{}) {
	this.raiseTypedString(fmt.Sprintf("%v", v), reflect.TypeOf(v).String())
}