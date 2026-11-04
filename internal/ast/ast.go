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

// Dice is a roll of Count dice, each with Sides. The parser sets Count to 1
// when omitted ("d20" == "1d20").
// Keep optionally selects a subset of the rolls (kh/kl/dh/dl).
type Dice struct {
	Count int
	Sides int
	Keep  *KeepRule
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
