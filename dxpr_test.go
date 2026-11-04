package dxpr_test

import (
	"errors"
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
			got, err := dxpr.EvaluateWithRoller(tt.in, &seqRoller{rolls: []int{3}})
			if err != nil {
				t.Fatalf("EvaluateWithRoller(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("EvaluateWithRoller(%q) = %d, want %d", tt.in, got, tt.want)
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
		{"negative exponent", "2 ^ -1"},
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

func TestEvaluateWithRollerCustomSequence(t *testing.T) {
	// A roller returning distinct values exercises multi-dice summation.
	roller := &seqRoller{rolls: []int{1, 2, 3}}
	got, err := dxpr.EvaluateWithRoller("3d6", roller)
	if err != nil {
		t.Fatalf("EvaluateWithRoller: %v", err)
	}
	if got != 6 {
		t.Errorf("EvaluateWithRoller(3d6) = %d, want 6", got)
	}
	if roller.next != 3 {
		t.Errorf("Roll called %d times, want 3", roller.next)
	}
}

func TestRollerInterface(t *testing.T) {
	// The public Roller interface must accept any implementation.
	var _ dxpr.Roller = (*errRoller)(nil)
	_, err := dxpr.EvaluateWithRoller("d6", &errRoller{})
	if err == nil {
		t.Error("EvaluateWithRoller succeeded with failing roller, want error")
	}
}

type errRoller struct{}

func (errRoller) Roll(sides int) (int, error) {
	return 0, errors.New(fmt.Sprintf("boom: %d sides", sides))
}
