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
		{"multiply alias same tree", "2x3", bin("*", num(2), num(3))},
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
		{"dice with expression count", "(3+6)d(3-1)", &ast.Dice{Count: 9, Sides: 2}},
		{"dice with expression sides", "2d(3-1)", &ast.Dice{Count: 2, Sides: 2}},
		{"dice count converted to int", "(7/2)d6", &ast.Dice{Count: 3, Sides: 6}},
		{"keep highest", "2d20kh1", dicePool(20, 2, &ast.KeepRule{Keep: 1, Highest: true})},
		{"keep lowest", "4d6kl2", dicePool(6, 4, &ast.KeepRule{Keep: 2, Highest: false})},
		{"drop lowest", "4d6dl1", dicePool(6, 4, &ast.KeepRule{Keep: 3, Highest: true})},
		{"drop highest", "3d6dh1", dicePool(6, 3, &ast.KeepRule{Keep: 2, Highest: false})},
		{"modifier with spaces", "4d6 dl 1", dicePool(6, 4, &ast.KeepRule{Keep: 3, Highest: true})},
		{"modifier on countless dice", "d20kh1", dicePool(20, 1, &ast.KeepRule{Keep: 1, Highest: true})},
		{"modifier without count defaults to 1", "2d20kh", dicePool(20, 2, &ast.KeepRule{Keep: 1, Highest: true})},
		{"drop on single die is a no-op", "d20dl1", dicePool(20, 1, &ast.KeepRule{Keep: 1, Highest: true})},
		{"drop without count defaults to 1", "4d6dl", dicePool(6, 4, &ast.KeepRule{Keep: 3, Highest: true})},
		{"modifier count as expression", "4d6kh(2-1)", dicePool(6, 4, &ast.KeepRule{Keep: 1, Highest: true})},
		{"modifier in arithmetic", "2d20kh1 + 5", bin("+", dicePool(20, 2, &ast.KeepRule{Keep: 1, Highest: true}), num(5))},
		{"heterogeneous pool keep highest", "(d6, d10, 10)kh", &ast.Pool{
			Elements: []ast.Node{&ast.Dice{Count: 1, Sides: 6}, &ast.Dice{Count: 1, Sides: 10}, num(10)},
			Keep:     ast.KeepRule{Keep: 1, Highest: true},
		}},
		{"heterogeneous pool drop lowest", "(d6, d10, 10)dl", &ast.Pool{
			Elements: []ast.Node{&ast.Dice{Count: 1, Sides: 6}, &ast.Dice{Count: 1, Sides: 10}, num(10)},
			Keep:     ast.KeepRule{Keep: 2, Highest: true},
		}},
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
		{"modifier keep too many", "4d6kh5"},
		{"modifier keep zero", "4d6kh0"},
		{"modifier drop all dice", "4d6dl4"},

		{"pool without modifier", "(1, 2)"},
		{"pool with unknown modifier", "(1, 2)kz"},
		{"dice in dice count", "(2d6)d6"},
		{"dice in dice sides", "2d(3d1)"},
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

// dicePool builds a Pool of n dice with the given sides and keep rule.
func dicePool(sides, n int, keep *ast.KeepRule) *ast.Pool {
	elements := make([]ast.Node, n)
	for i := range elements {
		elements[i] = &ast.Dice{Count: 1, Sides: sides}
	}
	return &ast.Pool{Elements: elements, Keep: *keep}
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
	case *ast.Pool:
		b, ok := b.(*ast.Pool)
		if !ok || len(a.Elements) != len(b.Elements) {
			return false
		}
		for i := range a.Elements {
			if !equal(a.Elements[i], b.Elements[i]) {
				return false
			}
		}
		return a.Keep == b.Keep
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
