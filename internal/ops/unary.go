package ops

import (
	"fmt"
)

// Unary is a prefix (unary) operator.
type Unary struct {
	// Symbol is the token text of the operator.
	Symbol string
	// BindingPower is the operator's precedence; higher binds tighter.
	BindingPower int
	// Apply computes the operator's result.
	Apply func(operand float64) (float64, error)
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
