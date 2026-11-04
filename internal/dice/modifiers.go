package dice

import (
	"fmt"

	"github.com/shvixxl/dxpr/internal/ast"
)

// Modifier is a dice modifier suffix such as kh1 (keep highest 1).
type Modifier struct {
	// Name is the modifier keyword (e.g. "kh").
	Name string
	// Resolve converts the pool size and the modifier's count into a
	// KeepRule, validating ranges. The count defaults to 1 when omitted
	// (2d20kh == 2d20kh1).
	Resolve func(count, n int) (*ast.KeepRule, error)
}

// keepRange validates that n is in [1, max].
func keepRange(mod string, n, max int) error {
	if n < 1 || n > max {
		return fmt.Errorf("%s: count %d out of range [1, %d]", mod, n, max)
	}
	return nil
}

// KeepHighest keeps the highest n rolls (kh).
var KeepHighest = Modifier{
	Name:        "kh",
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
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("dh", n, count-1); err != nil {
			return nil, err
		}
		return &ast.KeepRule{Keep: count - n, Highest: false}, nil
	},
}

// DropLowest drops the lowest n rolls, keeping the rest (dl).
var DropLowest = Modifier{
	Name:        "dl",
	Resolve: func(count, n int) (*ast.KeepRule, error) {
		if err := keepRange("dl", n, count-1); err != nil {
			return nil, err
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
