// Package parser converts a token stream into an abstract syntax tree using a
// Pratt (precedence-climbing) parser driven by the operator registry.
//
// Grammar:
//
//	expression := atom binaryOp expression?    // precedence from ops registry
//	atom       := [number] 'd' number | number
package parser

import (
	"fmt"
	"strconv"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/dice"
	"github.com/shvixxl/dxpr/internal/funcs"
	"github.com/shvixxl/dxpr/internal/ops"
	"github.com/shvixxl/dxpr/internal/tokenizer"
)

// Parse builds an AST from tokens.
func Parse(toks []tokenizer.Token) (ast.Node, error) {
	p := &parser{toks: toks}
	node, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if p.peek().Kind != tokenizer.KindEOF {
		return nil, fmt.Errorf("unexpected %s after expression", p.peek().Describe())
	}
	return node, nil
}

type parser struct {
	toks []tokenizer.Token
	pos  int
}

func (p *parser) peek() tokenizer.Token { return p.toks[p.pos] }

func (p *parser) next() tokenizer.Token {
	tok := p.toks[p.pos]
	if tok.Kind != tokenizer.KindEOF {
		p.pos++
	}
	return tok
}

func (p *parser) expect(kind tokenizer.Kind) (tokenizer.Token, error) {
	tok := p.next()
	if tok.Kind != kind {
		return tok, fmt.Errorf("expected %s, got %s", kind.Describe(), tok.Describe())
	}
	return tok, nil
}

func (p *parser) number() (int, error) {
	tok, err := p.expect(tokenizer.KindNumber)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(tok.Value)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", tok.Value)
	}
	return n, nil
}

// expression parses atoms and binary operators with binding power of at
// least minBP, climbing precedence as operators allow.
func (p *parser) expression(minBP int) (ast.Node, error) {
	left, err := p.atom()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.peek()
		if tok.Kind != tokenizer.KindOperator {
			return left, nil
		}
		op, ok := ops.BinaryBySymbol[tok.Value]
		if !ok || op.BindingPower < minBP {
			return left, nil
		}
		p.next()

		// Left-associative operators allow equal binding power on the right;
		// right-associative ones require strictly greater.
		nextBP := op.BindingPower
		if !op.RightAssoc {
			nextBP++
		}

		right, err := p.expression(nextBP)
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOp{Op: op.Symbol, Left: left, Right: right}
	}
}

// atom := unaryOp atom | functionCall | '(' expression ')' | [number] 'd' number | number
func (p *parser) atom() (ast.Node, error) {
	// Parenthesized grouping.
	if p.peek().Kind == tokenizer.KindLParen {
		p.next()
		inner, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokenizer.KindRParen); err != nil {
			return nil, err
		}
		return inner, nil
	}

	// Function calls.
	if p.peek().Kind == tokenizer.KindIdent {
		return p.functionCall()
	}

	// Prefix (unary) operators.
	if tok := p.peek(); tok.Kind == tokenizer.KindOperator {
		if op, ok := ops.UnaryBySymbol[tok.Value]; ok {
			p.next()
			operand, err := p.expression(op.BindingPower)
			if err != nil {
				return nil, err
			}
			return &ast.UnaryOp{Op: op.Symbol, Operand: operand}, nil
		}
	}

	if p.peek().Kind == tokenizer.KindDice {
		p.next()
		sides, err := p.number()
		if err != nil {
			return nil, err
		}
		return p.diceModifiers(&ast.Dice{Count: 1, Sides: sides})
	}

	n, err := p.number()
	if err != nil {
		return nil, err
	}

	// "NdM": number followed by dice
	if p.peek().Kind == tokenizer.KindDice {
		p.next()
		sides, err := p.number()
		if err != nil {
			return nil, err
		}
		return p.diceModifiers(&ast.Dice{Count: n, Sides: sides})
	}

	return &ast.Number{Value: n}, nil
}

// diceModifiers applies an optional modifier suffix (kh, kl, dh, dl) to a
// dice node, e.g. 4d6dl1.
func (p *parser) diceModifiers(d *ast.Dice) (ast.Node, error) {
	tok := p.peek()
	if tok.Kind != tokenizer.KindIdent {
		return d, nil
	}
	mod, ok := dice.ModifiersByName[tok.Value]
	if !ok {
		return d, nil // not a modifier; let the caller report it if invalid
	}
	p.next()

	n := 1
	if mod.TakesNumber {
		var err error
		n, err = p.number()
		if err != nil {
			return nil, err
		}
	}

	keep, err := mod.Resolve(d.Count, n)
	if err != nil {
		return nil, err
	}
	d.Keep = keep
	return d, nil
}

// functionCall := ident '(' [expression (',' expression)*] ')'
func (p *parser) functionCall() (ast.Node, error) {
	name := p.next().Value
	fn, err := funcs.Lookup(name)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(tokenizer.KindLParen); err != nil {
		return nil, err
	}

	var args []ast.Node
	if p.peek().Kind != tokenizer.KindRParen {
		for {
			arg, err := p.expression(0)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.peek().Kind != tokenizer.KindComma {
				break
			}
			p.next()
		}
	}

	if _, err := p.expect(tokenizer.KindRParen); err != nil {
		return nil, err
	}
	if err := fn.CheckArgCount(len(args)); err != nil {
		return nil, err
	}
	return &ast.FnCall{Name: fn.Name, Args: args}, nil
}
