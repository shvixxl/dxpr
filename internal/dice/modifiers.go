package dice

import (
	"fmt"

	"github.com/shvixxl/dxpr/internal/ast"
)

// Modifier is a dice modifier suffix such as kh1 (keep highest 1).
type Modifier struct {
	// Name is the modifier keyword (e.g. "kh").
	Name string
	// TakesNumber reports whether the modifier takes an explicit count
	// (kh1) as opposed to taking none.
	TakesNumber bool
	// Resolve converts the pool size and the modifier's count into a
	// KeepRule, validating ranges.
	Resolve func(count, n int) (*ast.KeepRule, error)
}

// keepRange validates that n is in [1, count].
func keepRange(mod string, n, count int) error {
	if n < 1 || n > count {
		return fmt.Errorf("%s: count %d out of range [1, %d]", mod, n, count)
	}
	return nil
}

// KeepHighest keeps the highest n rolls (kh).
var KeepHighest = Modifier{
	Name:        "kh",
	TakesNumber: true,
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("kh", n, count); err != nil {
			return nil, err
		}
		return &ast.KeepRule{Keep: n, Highest: true}, nil
	},
}

// KeepLowest keeps the lowest n rolls (kl).
var KeepLowest = Modifier{
	Name:        "kl",
	TakesNumber: true,
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("kl", n, count); err != nil {
			return nil, err
		}
		return &ast.KeepRule{Keep: n, Highest: false}, nil
	},
}

// DropHighest drops the highest n rolls, keeping the rest (dh).
var DropHighest = Modifier{
	Name:        "dh",
	TakesNumber: true,
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("dh", n, count-1); err != nil {
			return nil, fmt.Errorf("dh: count %d out of range [1, %d]", n, count-1)
		}
		return &ast.KeepRule{Keep: count - n, Highest: false}, nil
	},
}

// DropLowest drops the lowest n rolls, keeping the rest (dl).
var DropLowest = Modifier{
	Name:        "dl",
	TakesNumber: true,
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("dl", n, count-1); err != nil {
			return nil, fmt.Errorf("dl: count %d out of range [1, %d]", n, count-1)
		}
		return &ast.KeepRule{Keep: count - n, Highest: true}, nil
	},
}

// modifiers is the full set of dice modifiers.
var modifiers = []Modifier{KeepHighest, KeepLowest, DropHighest, DropLowest}

// ModifiersByName maps modifier keywords to their definitions.
var ModifiersByName = buildModifiersByName()

func buildModifiersByName() map[string]Modifier {
	m := make(map[string]Modifier)
	for _, mod := range modifiers {
		m[mod.Name] = mod
	}
	return m
}

// LookupModifier returns the dice modifier with the given keyword.
func LookupModifier(name string) (Modifier, error) {
	mod, ok := ModifiersByName[name]
	if !ok {
		return Modifier{}, fmt.Errorf("unknown dice modifier %q", name)
	}
	return mod, nil
}
