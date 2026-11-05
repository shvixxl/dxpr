package ops

import (
	"math"
	"testing"
)

func TestBinaryApply(t *testing.T) {
	tests := []struct {
		name  string
		sym   string
		left  float64
		right float64
		want  float64
	}{
		{"add", "+", 2, 3, 5},
		{"sub", "-", 5, 3, 2},
		{"mul", "*", 4, 6, 24},
		{"div", "/", 8, 2, 4},
		{"div fraction", "/", 7, 2, 3.5},
		{"div negative fraction", "/", -7, 2, -3.5},
		{"floor div", "//", 7, 2, 3},
		{"floor div floors", "//", -7, 2, -4},
		{"floor div float", "//", 7.5, 2, 3},
		{"mod", "%", 7, 3, 1},
		{"pow", "^", 2, 10, 1024},
		{"pow alias", "**", 3, 3, 27},
		{"pow zero exponent", "^", 5, 0, 1},
		{"pow negative exponent", "^", 2, -1, 0.5},
		{"pow float base", "^", 2.5, 2, 6.25},
		{"add float", "+", 1, 2.5, 3.5},
		{"mul with float", "*", 2.5, 4, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := LookupBinary(tt.sym)
			if err != nil {
				t.Fatalf("LookupBinary(%q): %v", tt.sym, err)
			}
			got, err := op.Apply(tt.left, tt.right)
			if err != nil {
				t.Fatalf("Apply(%v, %v): %v", tt.left, tt.right, err)
			}
			if got != tt.want {
				t.Errorf("Apply(%v, %v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestBinaryErrors(t *testing.T) {
	tests := []struct {
		name  string
		sym   string
		left  float64
		right float64
	}{
		{"div by zero", "/", 5, 0},
		{"floor div by zero", "//", 5, 0},
		{"mod by zero", "%", 5, 0},
		{"mod float operand", "%", 1.5, 2},
		{"pow out of range", "^", 10, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := LookupBinary(tt.sym)
			if err != nil {
				t.Fatalf("LookupBinary(%q): %v", tt.sym, err)
			}
			if _, err := op.Apply(tt.left, tt.right); err == nil {
				t.Errorf("Apply(%v, %v) succeeded, want error", tt.left, tt.right)
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

func TestNegate(t *testing.T) {
	got, err := Negate.Apply(5)
	if err != nil || got != -5 {
		t.Errorf("Negate(5) = %v, %v; want -5", got, err)
	}
	got, err = Negate.Apply(2.5)
	if err != nil || got != -2.5 {
		t.Errorf("Negate(2.5) = %v, %v; want -2.5", got, err)
	}
}

func TestPowLargeExponent(t *testing.T) {
	// 2^63 is exactly representable in float64; beyond int64 the evaluator
	// rejects it, but the operator itself succeeds.
	got, err := Pow.Apply(2, 63)
	if err != nil || got != math.Pow(2, 63) {
		t.Errorf("Pow(2, 63) = %v, %v; want %v", got, err, math.Pow(2, 63))
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
