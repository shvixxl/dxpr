// Package value provides the numeric value type used during evaluation.
//
// Values are integers until an operation produces a fraction (e.g. 7 / 2),
// at which point they promote to floats. Integer arithmetic stays exact and
// fast; floats appear only where the math demands them.
package value

import "math"

// Value is a number: either an exact integer or a float.
type Value struct {
	i     int64
	f     float64
	isInt bool
}

// FromInt returns an integer Value.
func FromInt(n int64) Value { return Value{i: n, isInt: true} }

// FromFloat returns a float Value.
func FromFloat(x float64) Value { return Value{f: x} }

// IsInt reports whether the value is an exact integer.
func (v Value) IsInt() bool { return v.isInt }

// Int returns the value's integer part, truncating floats toward zero.
func (v Value) Int() int64 {
	if v.isInt {
		return v.i
	}
	return int64(v.f)
}

// Float returns the value as a float.
func (v Value) Float() float64 {
	if v.isInt {
		return float64(v.i)
	}
	return v.f
}

// Floor returns the value rounded down to an integer.
func (v Value) Floor() int64 {
	if v.isInt {
		return v.i
	}
	return int64(math.Floor(v.f))
}

// Less reports whether v is smaller than o. Integers compare exactly;
// mixed comparisons use floats.
func (v Value) Less(o Value) bool {
	if v.isInt && o.isInt {
		return v.i < o.i
	}
	return v.Float() < o.Float()
}
