package funcs

// Abs returns the absolute value of its argument.
var Abs = Function{
	Name:    "abs",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []int) (int, error) {
		if args[0] < 0 {
			return -args[0], nil
		}
		return args[0], nil
	},
}

// Round rounds its argument to the nearest integer. All values are integers
// today, so it is the identity; it exists for clarity in expressions and for
// future non-integer intermediates.
var Round = Function{
	Name:    "round",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []int) (int, error) {
		return args[0], nil
	},
}

// Ceil rounds its argument up to the nearest integer. All values are
// integers today, so it is the identity; it exists for clarity in expressions
// and for future non-integer intermediates.
var Ceil = Function{
	Name:    "ceil",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []int) (int, error) {
		return args[0], nil
	},
}

// Floor rounds its argument down to the nearest integer. All values are
// integers today, so it is the identity; it exists for clarity in expressions
// and for future non-integer intermediates.
var Floor = Function{
	Name:    "floor",
	MinArgs: 1,
	MaxArgs: 1,
	Apply: func(args []int) (int, error) {
		return args[0], nil
	},
}

// Min returns the smallest of its arguments.
var Min = Function{
	Name:    "min",
	MinArgs: 1,
	MaxArgs: -1,
	Apply: func(args []int) (int, error) {
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
	MinArgs: 1,
	MaxArgs: -1,
	Apply: func(args []int) (int, error) {
		largest := args[0]
		for _, v := range args[1:] {
			if v > largest {
				largest = v
			}
		}
		return largest, nil
	},
}
