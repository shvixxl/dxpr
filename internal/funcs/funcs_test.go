package funcs

import (
	"testing"

	"github.com/shvixxl/dxpr/internal/value"
)

func i(n int64) value.Value { return value.FromInt(n) }

func f(x float64) value.Value { return value.FromFloat(x) }

func TestApply(t *testing.T) {
	tests := []struct {
		name string
		fn   string
		args []value.Value
		want value.Value
	}{
		{"abs positive", "abs", []value.Value{i(5)}, i(5)},
		{"abs negative", "abs", []value.Value{i(-5)}, i(5)},
		{"abs zero", "abs", []value.Value{i(0)}, i(0)},
		{"abs float", "abs", []value.Value{f(-2.5)}, f(2.5)},
		{"round int", "round", []value.Value{i(7)}, i(7)},
		{"round half up", "round", []value.Value{f(3.5)}, i(4)},
		{"round half away from zero", "round", []value.Value{f(-3.5)}, i(-4)},
		{"floor int", "floor", []value.Value{i(7)}, i(7)},
		{"floor fraction", "floor", []value.Value{f(3.5)}, i(3)},
		{"floor negative fraction", "floor", []value.Value{f(-3.5)}, i(-4)},
		{"ceil int", "ceil", []value.Value{i(7)}, i(7)},
		{"ceil fraction", "ceil", []value.Value{f(3.5)}, i(4)},
		{"ceil negative fraction", "ceil", []value.Value{f(-3.5)}, i(-3)},
		{"min", "min", []value.Value{i(3), i(1), i(2)}, i(1)},
		{"min mixed int float", "min", []value.Value{i(3), f(1.5)}, f(1.5)},
		{"min negative", "min", []value.Value{i(-3), i(-1)}, i(-3)},
		{"max", "max", []value.Value{i(3), i(1), i(2)}, i(3)},
		{"max mixed", "max", []value.Value{i(3), f(3.5)}, f(3.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := Lookup(tt.fn)
			if err != nil {
				t.Fatalf("Lookup(%q): %v", tt.fn, err)
			}
			if err := fn.CheckArgCount(len(tt.args)); err != nil {
				t.Fatalf("CheckArgCount: %v", err)
			}
			got, err := fn.Apply(tt.args)
			if err != nil {
				t.Fatalf("Apply(%v): %v", tt.args, err)
			}
			if got != tt.want {
				t.Errorf("Apply(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestArgCounts(t *testing.T) {
	tests := []struct {
		name  string
		fn    string
		count int
		ok    bool
	}{
		{"abs ok", "abs", 1, true},
		{"abs too few", "abs", 0, false},
		{"abs too many", "abs", 2, false},
		{"round ok", "round", 1, true},
		{"round too few", "round", 0, false},
		{"ceil too many", "ceil", 2, false},
		{"floor too few", "floor", 0, false},
		{"min variadic ok", "min", 5, true},
		{"min single arg is error", "min", 1, false},
		{"min too few", "min", 0, false},
		{"max variadic ok", "max", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := Lookup(tt.fn)
			if err != nil {
				t.Fatalf("Lookup(%q): %v", tt.fn, err)
			}
			err = fn.CheckArgCount(tt.count)
			if tt.ok && err != nil {
				t.Errorf("CheckArgCount(%d) = %v, want nil", tt.count, err)
			}
			if !tt.ok && err == nil {
				t.Errorf("CheckArgCount(%d) = nil, want error", tt.count)
			}
		})
	}
}

func TestLookupUnknown(t *testing.T) {
	if _, err := Lookup("nosuch"); err == nil {
		t.Error("Lookup(\"nosuch\") succeeded, want error")
	}
}

func TestRegistryProperties(t *testing.T) {
	for name, fn := range ByName {
		if fn.Name != name {
			t.Errorf("name %q maps to function %q", name, fn.Name)
		}
		if fn.MinArgs < 1 {
			t.Errorf("function %q has MinArgs %d", name, fn.MinArgs)
		}
		if fn.MaxArgs >= 0 && fn.MaxArgs < fn.MinArgs {
			t.Errorf("function %q has MaxArgs %d < MinArgs %d", name, fn.MaxArgs, fn.MinArgs)
		}
		if fn.Apply == nil {
			t.Errorf("function %q has nil Apply", name)
		}
	}
}
