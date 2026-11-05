package ops

import (
	"fmt"
	"math"
)

// Add is the addition operator.
var Add = Binary{
	Symbol:       "+",
	BindingPower: BPAdditive,
	Apply:        func(a, b float64) (float64, error) { return a + b, nil },
}

// Sub is the subtraction operator.
var Sub = Binary{
	Symbol:       "-",
	BindingPower: BPAdditive,
	Apply:        func(a, b float64) (float64, error) { return a - b, nil },
}

// Mul is the multiplication operator.
var Mul = Binary{
	Symbol:       "*",
	Aliases:      []string{"x", "X"},
	BindingPower: BPMultiplicative,
	Apply:        func(a, b float64) (float64, error) { return a * b, nil },
}

// Div is true division: 7 / 2 = 3.5. Fractional results floor to int only
// at the end of evaluation.
var Div = Binary{
	Symbol:       "/",
	BindingPower: BPMultiplicative,
	Apply: func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	},
}

// FloorDiv is floor division: it always rounds down (7 // 2 = 3,
// -7 // 2 = -4).
var FloorDiv = Binary{
	Symbol:       "//",
	BindingPower: BPMultiplicative,
	Apply: func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return math.Floor(a / b), nil
	},
}

// Mod is the modulus operator. It requires integer operands.
var Mod = Binary{
	Symbol:       "%",
	BindingPower: BPMultiplicative,
	Apply: func(a, b float64) (float64, error) {
		if a != math.Trunc(a) || b != math.Trunc(b) {
			return 0, fmt.Errorf("modulus requires integer operands")
		}
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return math.Mod(a, b), nil
	},
}

// Negate is the unary minus operator.
// It binds looser than exponentiation: -2^2 is -(2^2).
var Negate = Unary{
	Symbol:       "-",
	BindingPower: BPMultiplicative,
	Apply:        func(operand float64) (float64, error) { return -operand, nil },
}

// Pow is the exponentiation operator. It is right-associative.
var Pow = Binary{
	Symbol:       "^",
	Aliases:      []string{"**"},
	BindingPower: BPPower,
	RightAssoc:   true,
	Apply: func(base, exp float64) (float64, error) {
		r := math.Pow(base, exp)
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return 0, fmt.Errorf("result out of range in %v^%v", base, exp)
		}
		return r, nil
	},
}
