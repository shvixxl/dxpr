// Package funcs defines the function registry.
//
// Functions are called by name with parenthesized, comma-separated arguments:
// min(3d6, 10). They are pure: Apply maps evaluated integer arguments to an
// integer result and never touches randomness.
package funcs

import "fmt"

// Function is a named callable.
type Function struct {
	// Name is the identifier the function is called by.
	Name string
	// MinArgs is the minimum number of arguments.
	MinArgs int
	// MaxArgs is the maximum number of arguments, or -1 for unbounded.
	MaxArgs int
	// Apply computes the function's result.
	Apply func(args []int) (int, error)
}

// CheckArgCount reports whether n is a valid argument count for the function.
func (f Function) CheckArgCount(n int) error {
	if n < f.MinArgs {
		return fmt.Errorf("%s: expected at least %d argument(s), got %d", f.Name, f.MinArgs, n)
	}
	if f.MaxArgs >= 0 && n > f.MaxArgs {
		return fmt.Errorf("%s: expected at most %d argument(s), got %d", f.Name, f.MaxArgs, n)
	}
	return nil
}

// functions is the full set of functions.
var functions = []Function{Abs, Round, Ceil, Floor, Min, Max}

// ByName maps function names to their definitions.
var ByName = buildByName()

func buildByName() map[string]Function {
	m := make(map[string]Function)
	for _, fn := range functions {
		m[fn.Name] = fn
	}
	return m
}

// Lookup returns the function with the given name.
func Lookup(name string) (Function, error) {
	fn, ok := ByName[name]
	if !ok {
		return Function{}, fmt.Errorf("unknown function %q", name)
	}
	return fn, nil
}
