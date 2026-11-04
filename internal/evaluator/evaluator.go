// Package evaluator walks the AST and computes a result.
package evaluator

import (
	"fmt"
	"math"

	"github.com/shvixxl/dxpr/internal/ast"
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
		switch n.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		case "//":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return int(math.Floor(float64(left) / float64(right))), nil
		case "%":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left % right, nil
		case "^", "**":
			if right < 0 {
				return 0, fmt.Errorf("negative exponent %d", right)
			}
			return int(math.Pow(float64(left), float64(right))), nil
		default:
			return 0, fmt.Errorf("unknown operator %q", n.Op)
		}

	default:
		return 0, fmt.Errorf("unknown node type %T", node)
	}
}
