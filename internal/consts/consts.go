// Package consts evaluates constant expressions: sub-expressions whose
// value is known before any dice are rolled, such as dice counts and sides.
//
// The parser uses it to resolve operands like "(3+6)d(3-1)". Dice rolls and
// pools are rejected: a count or sides must not depend on randomness.
package consts

import (
	"fmt"
	"math"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/funcs"
	"github.com/shvixxl/dxpr/internal/ops"
)

// maxInt bounds results so the int conversion is portable across platforms.
const maxInt = math.MaxInt32

// Evaluate computes a constant expression, flooring fractional results.
func Evaluate(node ast.Node) (int, error) {
	v, err := evalConst(node)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) || v > maxInt || v < -maxInt-1 {
		return 0, fmt.Errorf("result %v out of range", v)
	}
	return int(math.Floor(v)), nil
}

func evalConst(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.Number:
		return float64(n.Value), nil

	case *ast.BinaryOp:
		left, err := evalConst(n.Left)
		if err != nil {
			return 0, err
		}
		right, err := evalConst(n.Right)
		if err != nil {
			return 0, err
		}
		op, err := ops.LookupBinary(n.Op)
		if err != nil {
			return 0, err
		}
		return op.Apply(left, right)

	case *ast.UnaryOp:
		operand, err := evalConst(n.Operand)
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
			v, err := evalConst(arg)
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

	case *ast.Dice, *ast.Pool:
		return 0, fmt.Errorf("dice rolls are not allowed in constant expressions")

	default:
		return 0, fmt.Errorf("unknown node type %T", node)
	}
}
