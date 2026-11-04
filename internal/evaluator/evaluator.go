// Package evaluator walks the AST and computes a result.
package evaluator

import (
	"fmt"
	"sort"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/funcs"
	"github.com/shvixxl/dxpr/internal/ops"
	"github.com/shvixxl/dxpr/internal/value"
)

// Resource limits for dice evaluation. They keep adversarial input such as
// 999999999d6 from consuming unbounded memory and time.
const (
	maxDiceCount = 1000
	maxDiceSides = 10000
)

// Dice is satisfied by anything that can roll dice.
// It matches the public dxpr.Dice interface structurally.
type Dice interface {
	Roll(sides int) (int, error)
}

// Evaluate computes the value of an expression node. Fractional results are
// rounded down, matching the tabletop convention for halved values.
func Evaluate(node ast.Node, roller Dice) (int, error) {
	v, err := evalValue(node, roller)
	if err != nil {
		return 0, err
	}
	return int(v.Floor()), nil
}

// evalValue computes the value of an expression node, preserving fractions.
func evalValue(node ast.Node, roller Dice) (value.Value, error) {
	switch n := node.(type) {
	case *ast.Number:
		return value.FromInt(int64(n.Value)), nil

	case *ast.Dice:
		if n.Count < 1 {
			return value.Value{}, fmt.Errorf("invalid dice count: %d", n.Count)
		}
		if n.Count > maxDiceCount {
			return value.Value{}, fmt.Errorf("dice count %d too large (max %d)", n.Count, maxDiceCount)
		}
		if n.Sides < 1 {
			return value.Value{}, fmt.Errorf("invalid dice sides: %d", n.Sides)
		}
		if n.Sides > maxDiceSides {
			return value.Value{}, fmt.Errorf("dice sides %d too large (max %d)", n.Sides, maxDiceSides)
		}
		if n.Keep != nil && (n.Keep.Keep < 1 || n.Keep.Keep > n.Count) {
			return value.Value{}, fmt.Errorf("invalid keep count: %d (dice count %d)", n.Keep.Keep, n.Count)
		}
		rolls := make([]int, n.Count)
		for i := range rolls {
			roll, err := roller.Roll(n.Sides)
			if err != nil {
				return value.Value{}, err
			}
			rolls[i] = roll
		}
		if n.Keep != nil {
			sort.Ints(rolls)
			if n.Keep.Highest {
				rolls = rolls[len(rolls)-n.Keep.Keep:]
			} else {
				rolls = rolls[:n.Keep.Keep]
			}
		}
		total := int64(0)
		for _, roll := range rolls {
			total += int64(roll)
		}
		return value.FromInt(total), nil

	case *ast.BinaryOp:
		left, err := evalValue(n.Left, roller)
		if err != nil {
			return value.Value{}, err
		}
		right, err := evalValue(n.Right, roller)
		if err != nil {
			return value.Value{}, err
		}
		op, err := ops.LookupBinary(n.Op)
		if err != nil {
			return value.Value{}, err
		}
		return op.Apply(left, right)

	case *ast.UnaryOp:
		operand, err := evalValue(n.Operand, roller)
		if err != nil {
			return value.Value{}, err
		}
		op, err := ops.LookupUnary(n.Op)
		if err != nil {
			return value.Value{}, err
		}
		return op.Apply(operand)

	case *ast.FnCall:
		args := make([]value.Value, len(n.Args))
		for i, arg := range n.Args {
			v, err := evalValue(arg, roller)
			if err != nil {
				return value.Value{}, err
			}
			args[i] = v
		}
		fn, err := funcs.Lookup(n.Name)
		if err != nil {
			return value.Value{}, err
		}
		return fn.Apply(args)

	default:
		return value.Value{}, fmt.Errorf("unknown node type %T", node)
	}
}
