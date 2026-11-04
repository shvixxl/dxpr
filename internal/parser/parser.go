// Package parser converts a token stream into an abstract syntax tree.
package parser

import (
	"fmt"
	"strconv"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/tokenizer"
)

// Parse builds an AST from tokens using recursive descent.
//
// Grammar (lowest to highest precedence):
//
//	expression := term (('+' | '-') term)*
//	term       := power (('*' | '/' | '//' | '%') power)*
//	power      := dice (('^' | '**') power)?          // right-associative
//	dice       := [number] 'd' number
func Parse(toks []tokenizer.Token) (ast.Node, error) {
	p := &parser{toks: toks}
	node, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.peek().Kind != tokenizer.KindEOF {
		return nil, fmt.Errorf("unexpected %s after expression", p.peek().Kind.Describe())
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
		return tok, fmt.Errorf("expected %s, got %s", kind.Describe(), tok.Kind.Describe())
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

// expression := term (('+' | '-') term)*
func (p *parser) expression() (ast.Node, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.peek()
		if tok.Kind != tokenizer.KindPlus && tok.Kind != tokenizer.KindMinus {
			return left, nil
		}
		p.next()

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOp{Op: tok.Value, Left: left, Right: right}
	}
}

// term := power (('*' | '/' | '//' | '%') power)*
func (p *parser) term() (ast.Node, error) {
	left, err := p.power()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.peek()
		switch tok.Kind {
		case tokenizer.KindTimes, tokenizer.KindSlash,
			tokenizer.KindDoubleSlash, tokenizer.KindPercent:
			p.next()
		default:
			return left, nil
		}

		right, err := p.power()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOp{Op: tok.Value, Left: left, Right: right}
	}
}

// power := dice (('^' | '**') power)?   // right-associative
func (p *parser) power() (ast.Node, error) {
	base, err := p.dice()
	if err != nil {
		return nil, err
	}

	tok := p.peek()
	if tok.Kind != tokenizer.KindCaret && tok.Kind != tokenizer.KindStarStar {
		return base, nil
	}
	p.next()

	exp, err := p.power()
	if err != nil {
		return nil, err
	}

	return &ast.BinaryOp{Op: tok.Value, Left: base, Right: exp}, nil
}

// dice := [number] 'd' number
func (p *parser) dice() (ast.Node, error) {
	if p.peek().Kind == tokenizer.KindDice {
		p.next()
		sides, err := p.number()
		if err != nil {
			return nil, err
		}
		return &ast.Dice{Count: 1, Sides: sides}, nil
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
		return &ast.Dice{Count: n, Sides: sides}, nil
	}

	return &ast.Number{Value: n}, nil
}
