// Package ast defines the abstract syntax tree for dice expressions.
package ast

// Node is any node in the expression syntax tree.
type Node interface{}

// Number is an integer literal.
type Number struct {
	Value int
}

// Dice is a roll of Count dice, each with Sides.
// Count may be zero, meaning 1 (e.g. "d20" == "1d20").
type Dice struct {
	Count int
	Sides int
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
