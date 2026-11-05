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
// 999999999d6 from consuming unbounded memory and time. The parser also
// enforces them when resolving constant dice operands and building pools.
const (
	MaxDiceCount = 1000
	MaxDiceSides = 10000
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
		if n.Count > MaxDiceCount {
			return value.Value{}, fmt.Errorf("dice count %d too large (max %d)", n.Count, MaxDiceCount)
		}
		if n.Sides < 1 {
			return value.Value{}, fmt.Errorf("invalid dice sides: %d", n.Sides)
		}
		if n.Sides > MaxDiceSides {
			return value.Value{}, fmt.Errorf("dice sides %d too large (max %d)", n.Sides, MaxDiceSides)
		}
		total := int64(0)
		for i := 0; i < n.Count; i++ {
			roll, err := roller.Roll(n.Sides)
			if err != nil {
				return value.Value{}, err
			}
		total += int64(roll)
		}
		return value.FromInt(total), nil

	case *ast.Pool:
		if len(n.Elements) < 1 {
			return value.Value{}, fmt.Errorf("invalid pool: no elements")
		}
		if len(n.Elements) > MaxDiceCount {
			return value.Value{}, fmt.Errorf("pool size %d too large (max %d)", len(n.Elements), MaxDiceCount)
		}
		vals := make([]value.Value, len(n.Elements))
		for i, el := range n.Elements {
			v, err := evalValue(el, roller)
			if err != nil {
				return value.Value{}, err
			}
			vals[i] = v
		}
		if n.Keep != nil {
			if n.Keep.Keep < 1 || n.Keep.Keep > len(vals) {
				return value.Value{}, fmt.Errorf("invalid keep count: %d (pool size %d)", n.Keep.Keep, len(vals))
			}
			sort.SliceStable(vals, func(i, j int) bool {
				return vals[i].Float() < vals[j].Float()
			})
			if n.Keep.Highest {
				vals = vals[len(vals)-n.Keep.Keep:]
			} else {
				vals = vals[:n.Keep.Keep]
			}
		}
		allInt := true
		total := 0.0
		for _, v := range vals {
			if !v.IsInt() {
				allInt = false
			}
			total += v.Float()
		}
		if allInt {
			return value.FromInt(int64(total)), nil
		}
		return value.FromFloat(total), nil

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
