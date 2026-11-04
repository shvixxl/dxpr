package ops

import "fmt"

// Add is the addition operator.
var Add = Binary{
	Symbol:       "+",
	BindingPower: BPAdditive,
	Apply: func(left, right int) (int, error) {
		return left + right, nil
	},
}

// Sub is the subtraction operator.
var Sub = Binary{
	Symbol:       "-",
	BindingPower: BPAdditive,
	Apply: func(left, right int) (int, error) {
		return left - right, nil
	},
}

// Mul is the multiplication operator.
var Mul = Binary{
	Symbol:       "*",
	BindingPower: BPMultiplicative,
	Apply: func(left, right int) (int, error) {
		return left * right, nil
	},
}

// Div is integer (truncating) division.
var Div = Binary{
	Symbol:       "/",
	BindingPower: BPMultiplicative,
	Apply: func(left, right int) (int, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	},
}

// FloorDiv is floor division.
var FloorDiv = Binary{
	Symbol:       "//",
	BindingPower: BPMultiplicative,
	Apply: func(left, right int) (int, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		q := left / right
		if left%right != 0 && (left < 0) != (right < 0) {
			q--
		}
		return q, nil
	},
}

// Mod is the modulus operator.
var Mod = Binary{
	Symbol:       "%",
	BindingPower: BPMultiplicative,
	Apply: func(left, right int) (int, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left % right, nil
	},
}

// Negate is the unary minus operator.
// It binds looser than exponentiation: -2^2 is -(2^2).
var Negate = Unary{
	Symbol:       "-",
	BindingPower: BPMultiplicative,
	Apply: func(value int) (int, error) {
		return -value, nil
	},
}

// Pow is the exponentiation operator. It is right-associative.
var Pow = Binary{
	Symbol:       "^",
	Aliases:      []string{"**"},
	BindingPower: BPPower,
	RightAssoc:   true,
	Apply: func(base, exp int) (int, error) {
		if exp < 0 {
			return 0, fmt.Errorf("negative exponent %d", exp)
		}
		result := 1
		for i := 0; i < exp; i++ {
			result *= base
		}
		return result, nil
	},
}
