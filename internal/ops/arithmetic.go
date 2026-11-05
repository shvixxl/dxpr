package ops

import (
	"fmt"
	"math"

	"github.com/shvixxl/dxpr/internal/value"
)

// maxInt64 is the largest int64 value.
const maxInt64 = int64(math.MaxInt64)

// Add is the addition operator.
var Add = Binary{
	Symbol:       "+",
	BindingPower: BPAdditive,
	Apply: func(a, b value.Value) (value.Value, error) {
		if a.IsInt() && b.IsInt() {
			x, y := a.Int(), b.Int()
			r := x + y
			if (x > 0) == (y > 0) && (r > 0) != (x > 0) {
				return value.Value{}, fmt.Errorf("integer overflow in %d + %d", x, y)
			}
			return value.FromInt(r), nil
		}
		return value.FromFloat(a.Float() + b.Float()), nil
	},
}

// Sub is the subtraction operator.
var Sub = Binary{
	Symbol:       "-",
	BindingPower: BPAdditive,
	Apply: func(a, b value.Value) (value.Value, error) {
		if a.IsInt() && b.IsInt() {
			x, y := a.Int(), b.Int()
			r := x - y
			if x != 0 && y != 0 && (x > 0) != (y > 0) && (r > 0) != (x > 0) {
				return value.Value{}, fmt.Errorf("integer overflow in %d - %d", x, y)
			}
			return value.FromInt(r), nil
		}
		return value.FromFloat(a.Float() - b.Float()), nil
	},
}

// Mul is the multiplication operator.
var Mul = Binary{
	Symbol:       "*",
	Aliases:      []string{"x", "X"},
	BindingPower: BPMultiplicative,
	Apply: func(a, b value.Value) (value.Value, error) {
		if a.IsInt() && b.IsInt() {
			x, y := a.Int(), b.Int()
			if x != 0 && y != 0 && x*y/y != x {
				return value.Value{}, fmt.Errorf("integer overflow in %d * %d", x, y)
			}
			return value.FromInt(x * y), nil
		}
		return value.FromFloat(a.Float() * b.Float()), nil
	},
}

// Div is true division: it produces an integer when the division is exact
// and a float otherwise (7 / 2 = 3.5).
var Div = Binary{
	Symbol:       "/",
	BindingPower: BPMultiplicative,
	Apply: func(a, b value.Value) (value.Value, error) {
		if a.IsInt() && b.IsInt() {
			x, y := a.Int(), b.Int()
			if y == 0 {
				return value.Value{}, fmt.Errorf("division by zero")
			}
			if x%y == 0 {
				return value.FromInt(x / y), nil
			}
		}
		if b.Float() == 0 {
			return value.Value{}, fmt.Errorf("division by zero")
		}
		return value.FromFloat(a.Float() / b.Float()), nil
	},
}

// FloorDiv is floor division: it always rounds down to an integer.
var FloorDiv = Binary{
	Symbol:       "//",
	BindingPower: BPMultiplicative,
	Apply: func(a, b value.Value) (value.Value, error) {
		if b.Float() == 0 {
			return value.Value{}, fmt.Errorf("division by zero")
		}
		q := math.Floor(a.Float() / b.Float())
		if q > float64(maxInt64) || q < float64(-maxInt64-1) {
			return value.Value{}, fmt.Errorf("integer overflow in floor division")
		}
		return value.FromInt(int64(q)), nil
	},
}

// Mod is the modulus operator. It requires integer operands.
var Mod = Binary{
	Symbol:       "%",
	BindingPower: BPMultiplicative,
	Apply: func(a, b value.Value) (value.Value, error) {
		if !a.IsInt() || !b.IsInt() {
			return value.Value{}, fmt.Errorf("modulus requires integer operands")
		}
		x, y := a.Int(), b.Int()
		if y == 0 {
			return value.Value{}, fmt.Errorf("division by zero")
		}
		return value.FromInt(x % y), nil
	},
}

// Negate is the unary minus operator.
// It binds looser than exponentiation: -2^2 is -(2^2).
var Negate = Unary{
	Symbol:       "-",
	BindingPower: BPMultiplicative,
	Apply: func(operand value.Value) (value.Value, error) {
		if operand.IsInt() {
			return value.FromInt(-operand.Int()), nil
		}
		return value.FromFloat(-operand.Float()), nil
	},
}

// powMaxExp bounds integer exponents to keep evaluation fast and prevent
// unbounded loops on adversarial input such as 2^2000000000.
const powMaxExp = 63

// Pow is the exponentiation operator. It is right-associative.
// Integer bases with non-negative integer exponents stay exact; anything
// else (negative exponents, float operands) produces a float.
var Pow = Binary{
	Symbol:       "^",
	Aliases:      []string{"**"},
	BindingPower: BPPower,
	RightAssoc:   true,
	Apply: func(base, exp value.Value) (value.Value, error) {
		if base.IsInt() && exp.IsInt() && exp.Int() >= 0 {
			b, e := base.Int(), exp.Int()
			if e > powMaxExp {
				return value.Value{}, fmt.Errorf("exponent %d too large (max %d)", e, powMaxExp)
			}
			result := int64(1)
			for i := int64(0); i < e; i++ {
				if result != 0 && b != 0 && result*b/b != result {
					return value.Value{}, fmt.Errorf("integer overflow in %d^%d", b, e)
				}
				result *= b
			}
			return value.FromInt(result), nil
		}
		r := math.Pow(base.Float(), exp.Float())
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return value.Value{}, fmt.Errorf("result out of range in %v^%v", base.Float(), exp.Float())
		}
		return value.FromFloat(r), nil
	},
}
