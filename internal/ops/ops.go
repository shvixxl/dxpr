// Package ops defines the operator registries used by the parser and evaluator.
//
// Each operator is a self-contained unit describing its symbol, precedence and
// behavior. Operator definitions are grouped by theme (e.g. arithmetic.go);
// the machinery for each operator family (type + registry) lives in its own
// file (binary.go, unary.go).
package ops

// Binding power levels shared across operator families. Higher binds tighter.
const (
	BPAdditive       = 10 // + -
	BPMultiplicative = 20 // * / // % (and unary -)
	BPPower          = 30 // ^ **
)
