package ops

import (
	"testing"

	"github.com/shvixxl/dxpr/internal/value"
)

func i(n int64) value.Value { return value.FromInt(n) }

func f(x float64) value.Value { return value.FromFloat(x) }

func TestBinaryApply(t *testing.T) {
	tests := []struct {
		name  string
		sym   string
		left  value.Value
		right value.Value
		want  value.Value
	}{
		{"add", "+", i(2), i(3), i(5)},
		{"sub", "-", i(5), i(3), i(2)},
		{"mul", "*", i(4), i(6), i(24)},
		{"div exact", "/", i(8), i(2), i(4)},
		{"div fraction", "/", i(7), i(2), f(3.5)},
		{"div negative fraction", "/", i(-7), i(2), f(-3.5)},
		{"floor div", "//", i(7), i(2), i(3)},
		{"floor div floors", "//", i(-7), i(2), i(-4)},
		{"floor div float", "//", f(7.5), i(2), i(3)},
		{"mod", "%", i(7), i(3), i(1)},
		{"pow", "^", i(2), i(10), i(1024)},
		{"pow alias", "**", i(3), i(3), i(27)},
		{"pow zero exponent", "^", i(5), i(0), i(1)},
		{"pow negative exponent is float", "^", i(2), i(-1), f(0.5)},
		{"pow float base", "^", f(2.5), i(2), f(6.25)},
		{"add promotes to float", "+", i(1), f(2.5), f(3.5)},
		{"mul with float", "*", f(2.5), i(4), f(10)},
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
		left  value.Value
		right value.Value
	}{
		{"div by zero", "/", i(5), i(0)},
		{"div by zero float", "/", i(5), f(0)},
		{"floor div by zero", "//", i(5), i(0)},
		{"mod by zero", "%", i(5), i(0)},
		{"mod float operand", "%", f(1.5), i(2)},
		{"exponent too large", "^", i(2), i(64)},
		{"overflow", "^", i(2), i(63)},
		{"add overflow", "+", i(9223372036854775807), i(1)},
		{"mul overflow", "*", i(9223372036854775807), i(2)},
		{"pow float out of range", "^", i(10), i(400)},
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
	got, err := Negate.Apply(i(5))
	if err != nil || got != i(-5) {
		t.Errorf("Negate(5) = %v, %v; want -5", got, err)
	}
	got, err = Negate.Apply(f(2.5))
	if err != nil || got != f(-2.5) {
		t.Errorf("Negate(2.5) = %v, %v; want -2.5", got, err)
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
