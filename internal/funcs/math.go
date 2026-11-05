package funcs

import (
	"math"
)

// Abs returns the absolute value of its argument.
var Abs = Function{
	Name:    "abs",
	MinArgs: 1,
	MaxArgs: 1,
	Apply:   func(args []float64) (float64, error) { return math.Abs(args[0]), nil },
}

// Round rounds its argument to the nearest integer (half away from zero).
var Round = Function{
	Name:    "round",
	MinArgs: 1,
	MaxArgs: 1,
	Apply:   func(args []float64) (float64, error) { return math.Round(args[0]), nil },
}

// Floor rounds its argument down to the nearest integer.
var Floor = Function{
	Name:    "floor",
	MinArgs: 1,
	MaxArgs: 1,
	Apply:   func(args []float64) (float64, error) { return math.Floor(args[0]), nil },
}

// Ceil rounds its argument up to the nearest integer.
var Ceil = Function{
	Name:    "ceil",
	MinArgs: 1,
	MaxArgs: 1,
	Apply:   func(args []float64) (float64, error) { return math.Ceil(args[0]), nil },
}

// Min returns the smallest of its arguments.
var Min = Function{
	Name:    "min",
	MinArgs: 2,
	MaxArgs: -1,
	Apply: func(args []float64) (float64, error) {
		smallest := args[0]
		for _, v := range args[1:] {
			if v < smallest {
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
	Apply: func(args []float64) (float64, error) {
		largest := args[0]
		for _, v := range args[1:] {
			if v > largest {
				largest = v
			}
		}
		return largest, nil
	},
}
