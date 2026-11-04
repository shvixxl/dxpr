// Package ops defines the operator registry used by the parser and evaluator.
//
// Each operator is a self-contained unit describing its symbol, precedence and
// behavior. To add an operator, define it in this package (grouped by theme in
// its own file) and reference it in BinaryBySymbol below.
package ops

import "fmt"

// Binding power levels for binary operators. Higher binds tighter.
const (
	BPAdditive       = 10 // + -
	BPMultiplicative = 20 // * / // %
	BPPower          = 30 // ^ **
)

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

// BinaryBySymbol maps operator symbols (including aliases) to their definitions.
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

// Unary is a prefix (unary) operator.
type Unary struct {
	// Symbol is the token text of the operator.
	Symbol string
	// BindingPower is the operator's precedence; higher binds tighter.
	BindingPower int
	// Apply computes the operator's result.
	Apply func(value int) (int, error)
}

// unaries is the full set of unary operators.
var unaries = []Unary{Negate}

// UnaryBySymbol maps unary operator symbols to their definitions.
var UnaryBySymbol = buildUnaryBySymbol()

func buildUnaryBySymbol() map[string]Unary {
	m := make(map[string]Unary)
	for _, op := range unaries {
		m[op.Symbol] = op
	}
	return m
}

// LookupUnary returns the unary operator for a symbol.
func LookupUnary(symbol string) (Unary, error) {
	op, ok := UnaryBySymbol[symbol]
	if !ok {
		return Unary{}, fmt.Errorf("unknown unary operator %q", symbol)
	}
	return op, nil
}

// LookupBinary returns the binary operator for a symbol.
func LookupBinary(symbol string) (Binary, error) {
	op, ok := BinaryBySymbol[symbol]
	if !ok {
		return Binary{}, fmt.Errorf("unknown operator %q", symbol)
	}
	return op, nil
}
