// Package parser converts a token stream into an abstract syntax tree using a
// Pratt (precedence-climbing) parser driven by the operator registry.
//
// Grammar:
//
//	expression := atom ('d' diceSuffix | binaryOp expression)?  // 'd' binds tightest
//	atom       := pool | functionCall | '(' expression ')' | 'd' diceSuffix | number
//	pool       := '(' expression (',' expression)+ ')' modifier
//	diceSuffix := diceSides modifier?
package parser

import (
	"fmt"
	"strconv"

	"github.com/shvixxl/dxpr/internal/ast"
	"github.com/shvixxl/dxpr/internal/consts"
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

		// A 'd' suffix turns any constant expression into a dice count:
		// "(3+6)d(3-1)" == "9d2". Dice binds tighter than any operator.
		if tok.Kind == tokenizer.KindDice {
			p.next()
			left, err = p.diceSuffix(left)
			if err != nil {
				return nil, err
			}
			continue
		}

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
	// Parenthesized grouping or pool: "(a, b, c)kh1" is a pool and
	// requires a modifier; "(a)" is plain grouping.
	if p.peek().Kind == tokenizer.KindLParen {
		p.next()
		first, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		elements := []ast.Node{first}
		isPool := false
		for p.peek().Kind == tokenizer.KindComma {
			p.next()
			el, err := p.expression(0)
			if err != nil {
				return nil, err
			}
			elements = append(elements, el)
			isPool = true
		}
		if _, err := p.expect(tokenizer.KindRParen); err != nil {
			return nil, err
		}
		if !isPool {
			return first, nil
		}
		return p.poolSuffix(elements)
	}

	// A bare 'd' lexes as an identifier; give a helpful error instead of
	// "unknown function".
	if tok := p.peek(); tok.Kind == tokenizer.KindIdent && (tok.Value == "d" || tok.Value == "D") {
		return nil, fmt.Errorf("expected number of sides after 'd'")
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

	// Bare "d20" / "d(3-1)": count defaults to 1.
	if p.peek().Kind == tokenizer.KindDice {
		p.next()
		return p.diceSuffix(&ast.Number{Value: 1})
	}

	n, err := p.number()
	if err != nil {
		return nil, err
	}
	return &ast.Number{Value: n}, nil
}

// diceSuffix parses what follows a consumed 'd' token: the sides (a number
// or a parenthesized expression) and an optional pool modifier (kh/kl/dh/dl).
// countNode is an expression for the number of dice. Constant operands are
// folded to ints; dice-containing ones are kept as dynamic expressions,
// evaluated once per group at roll time. The result is a Dice node, or a Pool
// when a modifier is present (modifiers require constant operands).
func (p *parser) diceSuffix(countNode ast.Node) (ast.Node, error) {
	d := &ast.Dice{}
	if consts.IsConstant(countNode) {
		count, err := consts.Evaluate(countNode)
		if err != nil {
			return nil, fmt.Errorf("invalid dice count: %v", err)
		}
		d.Count = count
	} else {
		d.CountExpr = countNode
	}

	sidesNode, err := p.diceSides()
	if err != nil {
		return nil, err
	}
	if consts.IsConstant(sidesNode) {
		sides, err := consts.Evaluate(sidesNode)
		if err != nil {
			return nil, fmt.Errorf("invalid dice sides: %v", err)
		}
		d.Sides = sides
	} else {
		d.SidesExpr = sidesNode
	}

	tok := p.peek()
	if tok.Kind != tokenizer.KindIdent {
		return d, nil
	}
	mod, ok := dice.ModifiersByName[tok.Value]
	if !ok {
		return d, nil // not a modifier; let the caller report it if invalid
	}
	p.next()

	// Modifiers select among individual rolls, so the dice must expand into
	// a pool at parse time; that requires constant count and sides.
	if d.CountExpr != nil || d.SidesExpr != nil {
		return nil, fmt.Errorf("dice modifiers require constant count and sides")
	}
	count, sides := d.Count, d.Sides
	if count > dice.MaxCount {
		return nil, fmt.Errorf("dice count %d too large (max %d)", count, dice.MaxCount)
	}
	// A modifier selects among individual rolls, so the dice expand into a
	// pool: "4d6kh1" becomes a pool of four d6 rolls.
	elements := make([]ast.Node, count)
	for i := range elements {
		elements[i] = &ast.Dice{Count: 1, Sides: sides}
	}
	return p.applyModifier(mod, elements)
}

// diceSides parses the sides of a dice roll: a number ("d6") or a
// parenthesized expression ("d(3-1)", "d(3d10kh)"). It returns the
// expression node; the caller folds constants and keeps the rest dynamic.
func (p *parser) diceSides() (ast.Node, error) {
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
	sides, err := p.number()
	if err != nil {
		return nil, fmt.Errorf("expected number of sides after 'd'")
	}
	return &ast.Number{Value: sides}, nil
}

// poolSuffix applies a mandatory pool modifier (kh/kl/dh/dl) to elements.
func (p *parser) poolSuffix(elements []ast.Node) (ast.Node, error) {
	tok := p.peek()
	mod, ok := dice.ModifiersByName[tok.Value]
	if tok.Kind != tokenizer.KindIdent || !ok {
		return nil, fmt.Errorf("expected pool modifier (kh, kl, dh, dl) after pool, got %s", tok.Describe())
	}
	p.next()
	if len(elements) > dice.MaxCount {
		return nil, fmt.Errorf("pool size %d too large (max %d)", len(elements), dice.MaxCount)
	}
	return p.applyModifier(mod, elements)
}

// applyModifier parses the modifier's optional count (defaulting to 1, and
// itself a number or parenthesized constant expression) and builds a Pool.
func (p *parser) applyModifier(mod dice.Modifier, elements []ast.Node) (ast.Node, error) {
	n := 1
	switch p.peek().Kind {
	case tokenizer.KindNumber:
		var err error
		n, err = p.number()
		if err != nil {
			return nil, err
		}
	case tokenizer.KindLParen:
		p.next()
		inner, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokenizer.KindRParen); err != nil {
			return nil, err
		}
		n, err = p.constInt(inner)
		if err != nil {
			return nil, err
		}
	}

	keep, err := mod.Resolve(len(elements), n)
	if err != nil {
		return nil, err
	}
	return &ast.Pool{Elements: elements, Keep: *keep}, nil
}

// constInt evaluates a constant expression to an int. Dice rolls are
// rejected: operands like counts and sides must be known before rolling.
// Fractional results are converted to int (floored).
func (p *parser) constInt(n ast.Node) (int, error) {
	v, err := consts.Evaluate(n)
	if err != nil {
		return 0, fmt.Errorf("must be a constant expression: %v", err)
	}
	return v, nil
}

// diceModifiers applies an optional modifier suffix (kh, kl, dh, dl) to a
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
