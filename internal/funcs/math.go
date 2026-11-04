package funcs

import (
	"math"

	"github.com/shvixxl/dxpr/internal/value"
)

// Abs returns the absolute value of its argument.
var Abs = Function{
	Name:    "abs",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []value.Value) (value.Value, error) {
		v := args[0]
		if v.IsInt() {
			if v.Int() < 0 {
				return value.FromInt(-v.Int()), nil
			}
			return v, nil
		}
		return value.FromFloat(math.Abs(v.Float())), nil
	},
}

// Round rounds its argument to the nearest integer (half away from zero).
var Round = Function{
	Name:    "round",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []value.Value) (value.Value, error) {
		return value.FromInt(int64(math.Round(args[0].Float()))), nil
	},
}

// Floor rounds its argument down to the nearest integer.
var Floor = Function{
	Name:    "floor",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []value.Value) (value.Value, error) {
		return value.FromInt(args[0].Floor()), nil
	},
}

// Ceil rounds its argument up to the nearest integer.
var Ceil = Function{
	Name:    "ceil",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []value.Value) (value.Value, error) {
		return value.FromInt(int64(math.Ceil(args[0].Float()))), nil
	},
}

// Min returns the smallest of its arguments.
var Min = Function{
	Name:    "min",
	MinArgs: 2,
	MaxArgs: -1,
	Apply: func(args []value.Value) (value.Value, error) {
		smallest := args[0]
		for _, v := range args[1:] {
			if v.Less(smallest) {
				smallest = v
			}
		}
		return smallest, nil
	},
}

// Max returns the largest of its arguments.
var Max = Function{
	Name:    "max",
	MinArgs: 2,
	MaxArgs: -1,
	Apply: func(args []value.Value) (value.Value, error) {
		largest := args[0]
		for _, v := range args[1:] {
			if largest.Less(v) {
				largest = v
			}
		}
		return largest, nil
	},
}
