// Package ast defines the abstract syntax tree for dice expressions.
package ast

// Node is any node in the expression syntax tree.
type Node interface{}

// Number is an integer literal.
type Number struct {
	Value int
}

// KeepRule selects which rolls of a dice pool contribute to the total.
type KeepRule struct {
	// Keep is the number of rolls to keep.
	Keep int
	// Highest keeps the highest rolls instead of the lowest.
	Highest bool
}

// Dice is a roll of Count dice, each with Sides, summed. The parser sets
// Count to 1 when omitted ("d20" == "1d20").
//
// Count and Sides may instead come from expressions (CountExpr/SidesExpr),
// evaluated once per group at roll time: "(3d6kh)d(3d10kh)" first rolls the
// count and the sides, then rolls that many dice of that size. Constant
// operands are folded by the parser, leaving the Expr fields nil.
type Dice struct {
	Count     int
	CountExpr Node // dynamic count; nil when Count is folded
	Sides     int
	SidesExpr Node // dynamic sides; nil when Sides is folded
}

// Pool is a collection of expressions summed together. Keep selects a
// subset of the evaluated elements (kh/kl/dh/dl); its zero value keeps all.
// The parser always sets Keep, but the zero value behaves as "no selection"
// so hand-built pools stay valid.
type Pool struct {
	Elements []Node
	Keep     KeepRule
}

// BinaryOp is a binary operation (e.g. +, *, ^).
type BinaryOp struct {
	Op    string
	Left  Node
	Right Node
}

// UnaryOp is a prefix operation (e.g. unary -).
type UnaryOp struct {
	Op      string
	Operand Node
}

// FnCall is a function call, e.g. min(3d6, 10).
type FnCall struct {
	Name string
	Args []Node
}
