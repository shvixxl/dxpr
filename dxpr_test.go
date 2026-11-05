package dxpr_test

import (
	"fmt"
	"testing"

	dxpr "github.com/shvixxl/dxpr"
)

// seqRoller returns rolls from a fixed sequence, cycling if exhausted.
type seqRoller struct {
	rolls []int
	next  int
}

func (s *seqRoller) Roll(sides int) (int, error) {
	roll := s.rolls[s.next%len(s.rolls)]
	s.next++
	return roll, nil
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"number", "42", 42},
		{"addition", "1 + 2", 3},
		{"precedence", "2 + 3 * 4", 14},
		{"power", "2^10", 1024},
		{"power alias", "2 ** 10", 1024},
		{"multiply alias", "2 x 10", 20},
		{"multiply alias uppercase, no spaces", "2X10", 20},
		{"floor division", "7 // 2", 3},
		{"modulus", "7 % 3", 1},
		{"dice with count", "3d6", 9},        // seqRoller returns all 3s
		{"dice arithmetic", "d20 + 5", 8},    // seqRoller returns all 3s
		{"complex", "2d4 * 2 + 1d8 - 3", 12}, // (3+3)*2 + 3 - 3
		{"unary minus", "-3 + 5", 2},
		{"unary minus on dice", "-2d6 + 10", 4},  // -(3+3) + 10
		{"unary after operator", "d20 + -5", -2}, // 3 + -5
		{"unary minus with power", "-2^2", -4},
		{"double negation", "--5", 5},
		{"grouping", "(2 + 3) * 4", 20},
		{"grouping with dice", "(d6 + d6) * 2", 12}, // (3+3)*2
		{"abs", "abs(-5)", 5},
		{"round", "round(2^10)", 1024},
		{"min", "min(3d6, 10)", 9},
		{"max", "max(d20, 10)", 10},
		{"min variadic", "min(4, 2, 3)", 2},
		{"nested functions", "abs(min(1, -2))", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dxpr.EvaluateWithDice(tt.in, &seqRoller{rolls: []int{3}})
			if err != nil {
				t.Fatalf("EvaluateWithDice(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("EvaluateWithDice(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"unknown character", "d20 + foo"},
		{"dangling operator", "1 +"},
		{"double operator", "1 + * 2"},
		{"division by zero", "5 / 0"},
		{"floor division by zero", "5 // 0"},
		{"modulus by zero", "5 % 0"},
		{"modulus with float", "7 / 2 % 2"},
		{"dice without sides", "3d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := dxpr.Evaluate(tt.in); err == nil {
				t.Errorf("Evaluate(%q) succeeded, want error", tt.in)
			}
		})
	}
}

func TestEvaluateRealRandomness(t *testing.T) {
	// With the real roller, results must be within valid bounds.
	got, err := dxpr.Evaluate("d20")
	if err != nil {
		t.Fatalf("Evaluate(d20): %v", err)
	}
	if got < 1 || got > 20 {
		t.Errorf("Evaluate(d20) = %d, out of range [1, 20]", got)
	}
}

func TestEvaluateWithDiceCustomSequence(t *testing.T) {
	// A roller returning distinct values exercises multi-dice summation.
	roller := &seqRoller{rolls: []int{1, 2, 3}}
	got, err := dxpr.EvaluateWithDice("3d6", roller)
	if err != nil {
		t.Fatalf("EvaluateWithDice: %v", err)
	}
	if got != 6 {
		t.Errorf("EvaluateWithDice(3d6) = %d, want 6", got)
	}
	if roller.next != 3 {
		t.Errorf("Roll called %d times, want 3", roller.next)
	}
}

func TestFloatSemantics(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"exact division stays int", "8 / 2", 4},
		{"fraction floors at the end", "7 / 2", 3},
		{"negative fraction floors down", "-7 / 2", -4},
		{"round", "round(7 / 2)", 4},
		{"floor", "floor(7 / 2)", 3},
		{"ceil", "ceil(7 / 2)", 4},
		{"negative exponent is float", "2 ^ -1", 0},
		{"float flows through arithmetic", "7 / 2 + 1", 4},
		{"float floors only at the end", "ceil(7 / 2 * 2)", 7},
		{"half a roll", "d20 / 2", 1}, // seqRoller returns 3 → 1.5 → 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dxpr.EvaluateWithDice(tt.in, &seqRoller{rolls: []int{3}})
			if err != nil {
				t.Fatalf("EvaluateWithDice(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("EvaluateWithDice(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDynamicDice(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"expression count and sides", "(3+6)d(3-1)", 27}, // 9d2, seqRoller returns all 3s
		{"expression sides only", "2d(3-1)", 6},           // 2d2
		{"expression count only", "(3+6)d2", 27},           // 9d2
		{"count converted to int", "(7/2)d6", 9},            // 3d6
		{"constant function operand", "max(3, 9)d2", 27},    // 9d2
		{"dynamic dice in arithmetic", "(3+6)d(3-1) + 1", 28},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dxpr.EvaluateWithDice(tt.in, &seqRoller{rolls: []int{3}})
			if err != nil {
				t.Fatalf("EvaluateWithDice(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("EvaluateWithDice(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDynamicDiceErrors(t *testing.T) {
	for _, in := range []string{"(2d6)d6", "2d(3d1)", "(1, 2)", "(1, 2)kz", "(3+6)d", "(0)d6"} {
		if _, err := dxpr.EvaluateWithDice(in, &seqRoller{rolls: []int{3}}); err == nil {
			t.Errorf("EvaluateWithDice(%q) succeeded, want error", in)
		}
	}
}

func TestKeepDrop(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		rolls []int
		want  int
	}{
		{"keep highest is advantage", "2d20kh1", []int{5, 12}, 12},
		{"modifier without count defaults to 1", "2d20kh", []int{5, 12}, 12},
		{"drop without count defaults to 1", "4d6dl", []int{1, 6, 2, 4}, 12},
		{"keep lowest is disadvantage", "2d20kl1", []int{5, 12}, 5},
		{"drop lowest keeps rest", "4d6dl1", []int{1, 6, 2, 4}, 12},
		{"drop highest keeps rest", "3d6dh1", []int{1, 6, 2}, 3},
		{"keep in arithmetic", "2d20kh1 + 5", []int{5, 12}, 17},
		{"keep inside function", "min(2d20kh1, 10)", []int{5, 12}, 10},
		{"heterogeneous pool keep highest", "(d6, d10, 10)kh", []int{5, 12}, 12},
		{"heterogeneous pool keep lowest", "(d6, d10, 10)kl", []int{5, 12}, 5},
		{"pool modifier count as expression", "(d6, d10, 10)kh(2-1)", []int{5, 12}, 12},
		{"pool drop lowest keeps rest", "(d6, d10, 10)dl", []int{5, 12}, 22},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dxpr.EvaluateWithDice(tt.in, &seqRoller{rolls: tt.rolls})
			if err != nil {
				t.Fatalf("EvaluateWithDice(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("EvaluateWithDice(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDiceInterface(t *testing.T) {
	// The public Dice interface must accept any implementation.
	var _ dxpr.Dice = (*errRoller)(nil)
	_, err := dxpr.EvaluateWithDice("d6", &errRoller{})
	if err == nil {
		t.Error("EvaluateWithDice succeeded with failing roller, want error")
	}
}

type errRoller struct{}

func (errRoller) Roll(sides int) (int, error) {
	return 0, fmt.Errorf("boom: %d sides", sides)
}
