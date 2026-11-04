package funcs

import "testing"

func TestApply(t *testing.T) {
	tests := []struct {
		name string
		fn   string
		args []int
		want int
	}{
		{"abs positive", "abs", []int{5}, 5},
		{"abs negative", "abs", []int{-5}, 5},
		{"abs zero", "abs", []int{0}, 0},
		{"round", "round", []int{7}, 7},
		{"ceil", "ceil", []int{7}, 7},
		{"floor", "floor", []int{7}, 7},
		{"min", "min", []int{3, 1, 2}, 1},
		{"min single", "min", []int{4}, 4},
		{"min negative", "min", []int{-3, -1}, -3},
		{"max", "max", []int{3, 1, 2}, 3},
		{"max single", "max", []int{4}, 4},
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
				t.Errorf("Apply(%v) = %d, want %d", tt.args, got, tt.want)
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
