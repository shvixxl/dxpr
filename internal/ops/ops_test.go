package ops

import "testing"

func TestBinaryApply(t *testing.T) {
	tests := []struct {
		name  string
		sym   string
		left  int
		right int
		want  int
	}{
		{"add", "+", 2, 3, 5},
		{"sub", "-", 5, 3, 2},
		{"mul", "*", 4, 6, 24},
		{"div", "/", 7, 2, 3},
		{"div truncates toward zero", "/", -7, 2, -3},
		{"floor div", "//", 7, 2, 3},
		{"floor div floors", "//", -7, 2, -4},
		{"mod", "%", 7, 3, 1},
		{"pow", "^", 2, 10, 1024},
		{"pow alias", "**", 3, 3, 27},
		{"pow zero exponent", "^", 5, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := LookupBinary(tt.sym)
			if err != nil {
				t.Fatalf("LookupBinary(%q): %v", tt.sym, err)
			}
			got, err := op.Apply(tt.left, tt.right)
			if err != nil {
				t.Fatalf("Apply(%d, %d): %v", tt.left, tt.right, err)
			}
			if got != tt.want {
				t.Errorf("Apply(%d, %d) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestBinaryErrors(t *testing.T) {
	tests := []struct {
		name  string
		sym   string
		left  int
		right int
	}{
		{"div by zero", "/", 5, 0},
		{"floor div by zero", "//", 5, 0},
		{"mod by zero", "%", 5, 0},
		{"negative exponent", "^", 2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := LookupBinary(tt.sym)
			if err != nil {
				t.Fatalf("LookupBinary(%q): %v", tt.sym, err)
			}
			if _, err := op.Apply(tt.left, tt.right); err == nil {
				t.Errorf("Apply(%d, %d) succeeded, want error", tt.left, tt.right)
			}
		})
	}
}

func TestLookupBinaryUnknown(t *testing.T) {
	if _, err := LookupBinary("<>"); err == nil {
		t.Error("LookupBinary(\"<>\") succeeded, want error")
	}
}

func TestRegistryProperties(t *testing.T) {
	for symbol, op := range BinaryBySymbol {
		if op.Symbol != symbol && !contains(op.Aliases, symbol) {
			t.Errorf("symbol %q maps to operator %q without matching Symbol or Aliases", symbol, op.Symbol)
		}
		if op.BindingPower <= 0 {
			t.Errorf("operator %q has non-positive binding power %d", symbol, op.BindingPower)
		}
		if op.Apply == nil {
			t.Errorf("operator %q has nil Apply", symbol)
		}
	}
}

func contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
