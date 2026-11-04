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

// BinaryOp is a binary operation: Add (+) or Subtract (-).
type BinaryOp struct {
	Op    string // "+" or "-"
	Left  Node
	Right Node
}
