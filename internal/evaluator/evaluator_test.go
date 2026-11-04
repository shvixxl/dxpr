package evaluator

import (
	"fmt"
	"testing"

	"github.com/shvixxl/dxpr/internal/ast"
)

// fakeRoller returns rolls from a fixed sequence, cycling if exhausted.
type fakeRoller struct {
	rolls []int
	next  int
}

func (f *fakeRoller) Roll(sides int) (int, error) {
	if sides < 1 {
		return 0, fmt.Errorf("invalid dice: %d sides", sides)
	}
	roll := f.rolls[f.next%len(f.rolls)]
	f.next++
	return roll, nil
}

func TestEvaluateNumbersAndOps(t *testing.T) {
	roller := &fakeRoller{rolls: []int{1}}
	tests := []struct {
		name string
		node ast.Node
		want int
	}{
		{"number", &ast.Number{Value: 42}, 42},
		{"negative number", &ast.Number{Value: -5}, -5},
		{"add", &ast.BinaryOp{Op: "+", Left: &ast.Number{Value: 2}, Right: &ast.Number{Value: 3}}, 5},
		{"mul", &ast.BinaryOp{Op: "*", Left: &ast.Number{Value: 4}, Right: &ast.Number{Value: 6}}, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.node, roller)
			if err != nil {
				t.Fatalf("Evaluate: %v", err)
			}
			if got != tt.want {
				t.Errorf("Evaluate = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEvaluateDice(t *testing.T) {
	tests := []struct {
		name  string
		node  ast.Node
		rolls []int
		want  int
	}{
		{"single die", &ast.Dice{Count: 1, Sides: 6}, []int{4}, 4},
		{"multiple dice summed", &ast.Dice{Count: 3, Sides: 6}, []int{1, 5, 3}, 9},
		{"dice in arithmetic", &ast.BinaryOp{
			Op:    "+",
			Left:  &ast.Dice{Count: 2, Sides: 6},
			Right: &ast.Number{Value: 10},
		}, []int{6, 2}, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.node, &fakeRoller{rolls: tt.rolls})
			if err != nil {
				t.Fatalf("Evaluate: %v", err)
			}
			if got != tt.want {
				t.Errorf("Evaluate = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	roller := &fakeRoller{rolls: []int{1}}
	tests := []struct {
		name string
		node ast.Node
	}{
		{"zero count", &ast.Dice{Count: 0, Sides: 6}},
		{"negative count", &ast.Dice{Count: -1, Sides: 6}},
		{"zero sides", &ast.Dice{Count: 1, Sides: 0}},
		{"unknown operator", &ast.BinaryOp{Op: "<>", Left: &ast.Number{Value: 1}, Right: &ast.Number{Value: 2}}},
		{"unknown unary operator", &ast.UnaryOp{Op: "<>", Operand: &ast.Number{Value: 1}}},
		{"unknown function", &ast.FnCall{Name: "nosuch", Args: []ast.Node{&ast.Number{Value: 1}}}},
		{"division by zero", &ast.BinaryOp{Op: "/", Left: &ast.Number{Value: 1}, Right: &ast.Number{Value: 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Evaluate(tt.node, roller); err == nil {
				t.Error("Evaluate succeeded, want error")
			}
		})
	}
}

func TestEvaluateRollsCorrectSides(t *testing.T) {
	var requested []int
	roller := &sidesRecorder{requested: &requested}

	_, err := Evaluate(&ast.Dice{Count: 2, Sides: 20}, roller)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if len(requested) != 2 || requested[0] != 20 || requested[1] != 20 {
		t.Errorf("Roll called with %v, want [20 20]", requested)
	}
}

type sidesRecorder struct {
	requested *[]int
}

func (s *sidesRecorder) Roll(sides int) (int, error) {
	*s.requested = append(*s.requested, sides)
	return 1, nil
}
