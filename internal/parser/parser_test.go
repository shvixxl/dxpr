package parser

import (
	"testing"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/tokenizer"
)

// num returns a Number node.
func num(n int) ast.Node { return &ast.Number{Value: n} }

// bin builds a nested BinaryOp from a symbol and operands.
func bin(op string, left, right ast.Node) ast.Node {
	return &ast.BinaryOp{Op: op, Left: left, Right: right}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want ast.Node
	}{
		{"number", "42", num(42)},
		{"dice", "d20", &ast.Dice{Count: 1, Sides: 20}},
		{"dice with count", "3d6", &ast.Dice{Count: 3, Sides: 6}},
		{"addition", "1+2", bin("+", num(1), num(2))},
		{"left assoc", "1-2-3", bin("-", bin("-", num(1), num(2)), num(3))},
		{"precedence", "2+3*4", bin("+", num(2), bin("*", num(3), num(4)))},
		{"power right assoc", "2^3^2", bin("^", num(2), bin("^", num(3), num(2)))},
		{"power binds tightest", "2*3^2", bin("*", num(2), bin("^", num(3), num(2)))},
		{"power alias same tree", "2**3**2", bin("^", num(2), bin("^", num(3), num(2)))},
		{"dice as operand", "d6+5", bin("+", &ast.Dice{Count: 1, Sides: 6}, num(5))},
		{"dice in arithmetic", "2d6*3", bin("*", &ast.Dice{Count: 2, Sides: 6}, num(3))},
		{"unary minus", "-5", &ast.UnaryOp{Op: "-", Operand: num(5)}},
		{"double negation", "--5", &ast.UnaryOp{Op: "-", Operand: &ast.UnaryOp{Op: "-", Operand: num(5)}}},
		{"unary on dice", "-3d6", &ast.UnaryOp{Op: "-", Operand: &ast.Dice{Count: 3, Sides: 6}}},
		{"unary binds looser than power", "-2^2", &ast.UnaryOp{Op: "-", Operand: bin("^", num(2), num(2))}},
		{"unary binds tighter than additive", "-2+3", bin("+", &ast.UnaryOp{Op: "-", Operand: num(2)}, num(3))},
		{"grouping overrides precedence", "(2+3)*4", bin("*", bin("+", num(2), num(3)), num(4))},
		{"nested grouping", "((5))", num(5)},
		{"grouping with unary", "-(2+3)", &ast.UnaryOp{Op: "-", Operand: bin("+", num(2), num(3))}},
		{"function call", "abs(5)", &ast.FnCall{Name: "abs", Args: []ast.Node{num(5)}}},
		{"function variadic", "min(3, 1, 2)", &ast.FnCall{Name: "min", Args: []ast.Node{num(3), num(1), num(2)}}},
		{"function with expression args", "max(1+1, 2^3)", &ast.FnCall{Name: "max", Args: []ast.Node{bin("+", num(1), num(1)), bin("^", num(2), num(3))}}},
		{"function in arithmetic", "abs(-5) + 1", bin("+", &ast.FnCall{Name: "abs", Args: []ast.Node{&ast.UnaryOp{Op: "-", Operand: num(5)}}}, num(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toks, err := tokenizer.Tokenize(tt.in)
			if err != nil {
				t.Fatalf("Tokenize(%q): %v", tt.in, err)
			}
			got, err := Parse(toks)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.in, err)
			}
			if !equal(got, tt.want) {
				t.Errorf("Parse(%q)\n got  %#v\n want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"dangling operator", "d20+"},
		{"leading operator", "*3"},
		{"double operator", "2+*3"},
		{"dice without sides", "3d"},
		{"trailing dice", "3d6d5"},
		{"operator after dice", "3d6*"},
		{"unbalanced paren", "(1+2"},
		{"close paren without open", "1+2)"},
		{"unknown function", "nosuch(1)"},
		{"function too few args", "abs()"},
		{"function too many args", "abs(1, 2)"},
		{"function missing paren", "abs(1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toks, err := tokenizer.Tokenize(tt.in)
			if err != nil {
				return // lexical errors are fine; we only care about parse errors
			}
			if _, err := Parse(toks); err == nil {
				t.Errorf("Parse(%q) succeeded, want error", tt.in)
			}
		})
	}
}

// equal compares two ASTs structurally.
func equal(a, b ast.Node) bool {
	switch a := a.(type) {
	case *ast.Number:
		b, ok := b.(*ast.Number)
		return ok && a.Value == b.Value
	case *ast.Dice:
		b, ok := b.(*ast.Dice)
		return ok && a.Count == b.Count && a.Sides == b.Sides
	case *ast.BinaryOp:
		b, ok := b.(*ast.BinaryOp)
		return ok && a.Op == b.Op && equal(a.Left, b.Left) && equal(a.Right, b.Right)
	case *ast.UnaryOp:
		b, ok := b.(*ast.UnaryOp)
		return ok && a.Op == b.Op && equal(a.Operand, b.Operand)
	case *ast.FnCall:
		b, ok := b.(*ast.FnCall)
		if !ok || a.Name != b.Name || len(a.Args) != len(b.Args) {
			return false
		}
		for i := range a.Args {
			if !equal(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
