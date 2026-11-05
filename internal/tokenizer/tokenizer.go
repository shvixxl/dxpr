// Package tokenizer converts an expression string into a stream of tokens.
//
// Operators are not hardcoded here: any punctuation sequence is matched
// against the ops registry (longest match first), making the registry the
// single source of truth for operator symbols.
package tokenizer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/shvixxl/dxpr/internal/ops"
)

// Kind identifies the type of a token.
type Kind int

const (
	KindNumber   Kind = iota
	KindDice          // 'd' or 'D'
	KindOperator      // any symbol registered in ops
	KindIdent         // function name
	KindLParen
	KindRParen
	KindComma
	KindEOF
)

// Token is a single lexical unit of an expression.
type Token struct {
	Kind  Kind
	Value string // literal text (digits for numbers, operator symbol for operators)
}

// Tokenize splits expr into tokens.
func Tokenize(expr string) ([]Token, error) {
	var toks []Token
	runes := []rune(expr)

	for i := 0; i < len(runes); {
		r := runes[i]
		switch {
		case unicode.IsSpace(r):
			i++
		case unicode.IsDigit(r):
			start := i
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
			}
			toks = append(toks, Token{Kind: KindNumber, Value: string(runes[start:i])})
		case (r == 'd' || r == 'D') && i+1 < len(runes) && unicode.IsDigit(runes[i+1]):
			toks = append(toks, Token{Kind: KindDice, Value: string(r)})
			i++
		case unicode.IsLetter(r):
			start := i
			for i < len(runes) && unicode.IsLetter(runes[i]) {
				i++
			}
			run := string(runes[start:i])
			// A letter run may be a registered operator alias (e.g. 'x' for
			// '*'); the ops registry is the single source of truth.
			if op, ok := ops.BinaryBySymbol[run]; ok {
				toks = append(toks, Token{Kind: KindOperator, Value: op.Symbol})
				continue
			}
			if op, ok := ops.UnaryBySymbol[run]; ok {
				toks = append(toks, Token{Kind: KindOperator, Value: op.Symbol})
				continue
			}
			toks = append(toks, Token{Kind: KindIdent, Value: run})
		case r == '(':
			toks = append(toks, Token{Kind: KindLParen, Value: "("})
			i++
		case r == ')':
			toks = append(toks, Token{Kind: KindRParen, Value: ")"})
			i++
		case r == ',':
			toks = append(toks, Token{Kind: KindComma, Value: ","})
			i++
		default:
			// Operator: consume the maximal punctuation run and match the
			// longest registered symbol against it.
			start := i
			for i < len(runes) && isPunct(runes[i]) {
				i++
			}
			run := string(runes[start:i])
			sym, ok := longestOperatorMatch(run)
			if !ok {
				return nil, fmt.Errorf("unexpected character %q in expression %q", string(r), expr)
			}
			toks = append(toks, Token{Kind: KindOperator, Value: sym})
			i = start + len([]rune(sym))
		}
	}

	toks = append(toks, Token{Kind: KindEOF})
	return toks, nil
}

// isPunct reports whether r may appear in an operator symbol.
func isPunct(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsDigit(r) && !unicode.IsLetter(r)
}

// longestOperatorMatch returns the longest registered operator symbol that is
// a prefix of run. It iterates runes so multi-byte characters are never split.
func longestOperatorMatch(run string) (string, bool) {
	runes := []rune(run)
	for size := len(runes); size > 0; size-- {
		prefix := string(runes[:size])
		if _, ok := ops.BinaryBySymbol[prefix]; ok {
			return prefix, true
		}
	}
	return "", false
}

// Describe returns a human-readable name for a token kind.
func (k Kind) Describe() string {
	switch k {
	case KindNumber:
		return "number"
	case KindDice:
		return "'d'"
	case KindOperator:
		return "operator"
	case KindIdent:
		return "identifier"
	case KindLParen:
		return "'('"
	case KindRParen:
		return "')'"
	case KindComma:
		return "','"
	case KindEOF:
		return "end of expression"
	default:
		return strings.ToLower(fmt.Sprintf("token(%d)", k))
	}
}

// Describe returns a human-readable description of the token, including the
// operator symbol when applicable.
func (t Token) Describe() string {
	if t.Kind == KindOperator {
		return "'" + t.Value + "'"
	}
	return t.Kind.Describe()
}
