package dice

import (
	"testing"
)

func TestCryptoDiceRollBounds(t *testing.T) {
	d := NewDice()
	for _, sides := range []int{1, 2, 6, 20, 100} {
		for i := 0; i < 100; i++ {
			roll, err := d.Roll(sides)
			if err != nil {
				t.Fatalf("Roll(%d): %v", sides, err)
			}
			if roll < 1 || roll > sides {
				t.Fatalf("Roll(%d) = %d, out of range [1, %d]", sides, roll, sides)
			}
		}
	}
}

func TestCryptoDiceRollInvalidSides(t *testing.T) {
	d := NewDice()
	for _, sides := range []int{0, -1, -20} {
		if _, err := d.Roll(sides); err == nil {
			t.Errorf("Roll(%d) succeeded, want error", sides)
		}
	}
}

func TestCryptoDiceRollCoversRange(t *testing.T) {
	// A d2 rolled enough times must produce both possible values.
	d := NewDice()
	seen := map[int]bool{}
	for i := 0; i < 100; i++ {
		roll, err := d.Roll(2)
		if err != nil {
			t.Fatalf("Roll(2): %v", err)
		}
		seen[roll] = true
	}
	for _, v := range []int{1, 2} {
		if !seen[v] {
			t.Errorf("Roll(2) never produced %d in 100 rolls", v)
		}
	}
}

func TestLookupModifier(t *testing.T) {
	for _, name := range []string{"kh", "k", "kl", "dh", "dl"} {
		if _, err := LookupModifier(name); err != nil {
			t.Errorf("LookupModifier(%q): %v", name, err)
		}
	}
	if _, err := LookupModifier("nosuch"); err == nil {
		t.Error("LookupModifier(\"nosuch\") succeeded, want error")
	}
}

func TestModifierRegistryProperties(t *testing.T) {
	for name, mod := range ModifiersByName {
		if mod.Name != name && !contains(mod.Aliases, name) {
			t.Errorf("name %q maps to modifier %q without matching Name or Aliases", name, mod.Name)
		}
		if mod.Resolve == nil {
			t.Errorf("modifier %q has nil Resolve", name)
		}
	}
}

func TestModifierResolve(t *testing.T) {
	tests := []struct {
		name    string
		mod     string
		count   int
		n       int
		want    *int  // expected Keep value, nil for error
		highest *bool // expected Highest value, nil for error
	}{
		{"kh min", "kh", 4, 1, intPtr(1), boolPtr(true)},
		{"kh max", "kh", 4, 4, intPtr(4), boolPtr(true)},
		{"kl min", "kl", 4, 1, intPtr(1), boolPtr(false)},
		{"kl max", "kl", 4, 4, intPtr(4), boolPtr(false)},
		{"dh one", "dh", 4, 1, intPtr(3), boolPtr(false)},
		{"dh all but one", "dh", 4, 3, intPtr(1), boolPtr(false)},
		{"dl one", "dl", 4, 1, intPtr(3), boolPtr(true)},
		{"dl all but one", "dl", 4, 3, intPtr(1), boolPtr(true)},
		{"kh zero", "kh", 4, 0, nil, nil},
		{"kh too many", "kh", 4, 5, nil, nil},
		{"kl zero", "kl", 4, 0, nil, nil},
		{"kh on single die", "kh", 1, 1, intPtr(1), boolPtr(true)},
		{"dh on single die is a no-op", "dh", 1, 1, intPtr(1), boolPtr(true)},
		{"dl on single die is a no-op", "dl", 1, 1, intPtr(1), boolPtr(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod, err := LookupModifier(tt.mod)
			if err != nil {
				t.Fatalf("LookupModifier(%q): %v", tt.mod, err)
			}
			keep, err := mod.Resolve(tt.count, tt.n)
			if tt.want == nil {
				if err == nil {
					t.Errorf("Resolve(%d, %d) succeeded, want error", tt.count, tt.n)
				}
				return
			}
			if err != nil {
				t.Fatalf("Resolve(%d, %d): %v", tt.count, tt.n, err)
			}
			if keep.Keep != *tt.want || keep.Highest != *tt.highest {
				t.Errorf("Resolve(%d, %d) = KeepRule{Keep: %d, Highest: %v}, want {Keep: %d, Highest: %v}",
					tt.count, tt.n, keep.Keep, keep.Highest, *tt.want, *tt.highest)
			}
		})
	}
}

func intPtr(n int) *int    { return &n }
func boolPtr(b bool) *bool { return &b }

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
