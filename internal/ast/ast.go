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
// Count to 1 when omitted ("d20" == "1d20"). Count and Sides may come from
// constant expressions: "(3+6)d(3-1)" == "9d2".
type Dice struct {
	Count int
	Sides int
}

// Pool is a collection of expressions summed together. Keep optionally
// selects a subset of the evaluated elements (kh/kl/dh/dl), e.g.
// "4d6kh1" is a pool of four d6 rolls keeping the highest, and
// "(d6, d10, 10)kh" is a heterogeneous pool keeping the highest.
type Pool struct {
	Elements []Node
	Keep     *KeepRule
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
