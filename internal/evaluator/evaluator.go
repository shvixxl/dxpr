// Package evaluator walks the AST and computes a result.
package evaluator

import (
	"fmt"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/funcs"
	"github.com/shvixxl/dxpr/internal/ops"
)

// Roller is satisfied by anything that can roll dice.
// It matches the public dxpr.Roller interface structurally.
type Roller interface {
	Roll(sides int) (int, error)
}

// Evaluate computes the value of an expression node.
func Evaluate(node ast.Node, roller Roller) (int, error) {
	switch n := node.(type) {
	case *ast.Number:
		return n.Value, nil

	case *ast.Dice:
		if n.Count < 1 {
			return 0, fmt.Errorf("invalid dice count: %d", n.Count)
		}
		if n.Sides < 1 {
			return 0, fmt.Errorf("invalid dice sides: %d", n.Sides)
		}
		total := 0
		for i := 0; i < n.Count; i++ {
			roll, err := roller.Roll(n.Sides)
			if err != nil {
				return 0, err
			}
			total += roll
		}
		return total, nil

	case *ast.BinaryOp:
		left, err := Evaluate(n.Left, roller)
		if err != nil {
			return 0, err
		}
		right, err := Evaluate(n.Right, roller)
		if err != nil {
			return 0, err
		}
		op, err := ops.LookupBinary(n.Op)
		if err != nil {
			return 0, err
		}
		return op.Apply(left, right)

	case *ast.UnaryOp:
		value, err := Evaluate(n.Operand, roller)
		if err != nil {
			return 0, err
		}
		op, err := ops.LookupUnary(n.Op)
		if err != nil {
			return 0, err
		}
		return op.Apply(value)

	case *ast.FnCall:
		args := make([]int, len(n.Args))
		for i, arg := range n.Args {
			v, err := Evaluate(arg, roller)
			if err != nil {
				return 0, err
			}
			args[i] = v
		}
		fn, err := funcs.Lookup(n.Name)
		if err != nil {
			return 0, err
		}
		return fn.Apply(args)

	default:
		return 0, fmt.Errorf("unknown node type %T", node)
	}
}
