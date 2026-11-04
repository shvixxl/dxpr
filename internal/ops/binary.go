package ops

import "fmt"

// Binary is a binary operator.
type Binary struct {
	// Symbol is the canonical token text of the operator (e.g. "+", "//").
	Symbol string
	// Aliases are alternative spellings of the same operator (e.g. "**" for "^").
	Aliases []string
	// BindingPower is the operator's precedence; higher binds tighter.
	BindingPower int
	// RightAssoc makes the operator right-associative (e.g. exponentiation).
	RightAssoc bool
	// Apply computes the operator's result.
	Apply func(left, right int) (int, error)
}

// binaries is the full set of binary operators.
var binaries = []Binary{Add, Sub, Mul, Div, FloorDiv, Mod, Pow}

// BinaryBySymbol maps binary operator symbols (including aliases) to their
// definitions.
var BinaryBySymbol = buildBinaryBySymbol()

func buildBinaryBySymbol() map[string]Binary {
	m := make(map[string]Binary)
	for _, op := range binaries {
		m[op.Symbol] = op
		for _, alias := range op.Aliases {
			m[alias] = op
		}
	}
	return m
}

// LookupBinary returns the binary operator for a symbol.
func LookupBinary(symbol string) (Binary, error) {
	op, ok := BinaryBySymbol[symbol]
	if !ok {
		return Binary{}, fmt.Errorf("unknown operator %q", symbol)
	}
	return op, nil
}
