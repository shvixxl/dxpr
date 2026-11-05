// Package evaluator walks the AST and computes a result.
package evaluator

import (
	"fmt"
	"math"
	"sort"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/dice"
	"github.com/shvixxl/dxpr/internal/funcs"
	"github.com/shvixxl/dxpr/internal/ops"
)

// Result bounds. Results are clamped to the int32 range so behavior is
// identical on 32-bit platforms, where int is 32 bits wide.
const (
	maxResult = math.MaxInt32
	minResult = math.MinInt32
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
	if math.IsNaN(v) || math.IsInf(v, 0) || v > maxResult || v < minResult {
		return 0, fmt.Errorf("result %v out of range (max %d)", v, maxResult)
	}
	return int(math.Floor(v)), nil
}

// evalValue computes the value of an expression node, preserving fractions.
func evalValue(node ast.Node, roller Dice) (float64, error) {
	switch n := node.(type) {
	case *ast.Number:
		return float64(n.Value), nil

	case *ast.Dice:
		if n.Count < 1 {
			return 0, fmt.Errorf("invalid dice count: %d", n.Count)
		}
		if n.Count > dice.MaxCount {
			return 0, fmt.Errorf("dice count %d too large (max %d)", n.Count, dice.MaxCount)
		}
		if n.Sides < 1 {
			return 0, fmt.Errorf("invalid dice sides: %d", n.Sides)
		}
		if n.Sides > dice.MaxSides {
			return 0, fmt.Errorf("dice sides %d too large (max %d)", n.Sides, dice.MaxSides)
		}
		total := 0.0
		for i := 0; i < n.Count; i++ {
			roll, err := roller.Roll(n.Sides)
			if err != nil {
				return 0, err
			}
			total += float64(roll)
		}
		return total, nil

	case *ast.Pool:
		if len(n.Elements) < 1 {
			return 0, fmt.Errorf("invalid pool: no elements")
		}
		if len(n.Elements) > dice.MaxCount {
			return 0, fmt.Errorf("pool size %d too large (max %d)", len(n.Elements), dice.MaxCount)
		}
		vals := make([]float64, len(n.Elements))
		for i, el := range n.Elements {
			v, err := evalValue(el, roller)
			if err != nil {
				return 0, err
			}
			vals[i] = v
		}
		if n.Keep.Keep > 0 {
			if n.Keep.Keep > len(vals) {
				return 0, fmt.Errorf("invalid keep count: %d (pool size %d)", n.Keep.Keep, len(vals))
			}
			sort.Float64s(vals)
			if n.Keep.Highest {
				vals = vals[len(vals)-n.Keep.Keep:]
			} else {
				vals = vals[:n.Keep.Keep]
			}
		}
		total := 0.0
		for _, v := range vals {
			total += v
		}
		return total, nil

	case *ast.BinaryOp:
		left, err := evalValue(n.Left, roller)
		if err != nil {
			return 0, err
		}
		right, err := evalValue(n.Right, roller)
		if err != nil {
			return 0, err
		}
		op, err := ops.LookupBinary(n.Op)
		if err != nil {
			return 0, err
		}
		return op.Apply(left, right)

	case *ast.UnaryOp:
		operand, err := evalValue(n.Operand, roller)
		if err != nil {
			return 0, err
		}
		op, err := ops.LookupUnary(n.Op)
		if err != nil {
			return 0, err
		}
		return op.Apply(operand)

	case *ast.FnCall:
		args := make([]float64, len(n.Args))
		for i, arg := range n.Args {
			v, err := evalValue(arg, roller)
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
