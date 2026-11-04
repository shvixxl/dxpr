// Package dxpr evaluates dice expressions such as "3d6 + 2" or "d20".
//
// It is usable both as a library and via the dxpr command-line tool.
// Only this package is public; the tokenizer, parser and evaluator are internal.
package dxpr

import (
	"github.com/shvixxl/dxpr/internal/dice"
	"github.com/shvixxl/dxpr/internal/evaluator"
	"github.com/shvixxl/dxpr/internal/parser"
	"github.com/shvixxl/dxpr/internal/tokenizer"
)

// Dice rolls dice. Implement it to control randomness (e.g. in tests).
type Dice interface {
	// Roll returns a random integer in [1, sides].
	Roll(sides int) (int, error)
}

// Evaluate parses and evaluates a dice expression, returning the result.
//
// Dice are rolled using the default roller, which is cryptographically secure.
func Evaluate(expr string) (int, error) {
	return EvaluateWithDice(expr, dice.NewDice())
}

// EvaluateWithDice is like Evaluate but rolls dice using the given Dice.
// This allows callers to control randomness (e.g. for tests).
func EvaluateWithDice(expr string, d Dice) (int, error) {
	toks, err := tokenizer.Tokenize(expr)
	if err != nil {
		return 0, err
	}

	node, err := parser.Parse(toks)
	if err != nil {
		return 0, err
	}

	return evaluator.Evaluate(node, d)
}
