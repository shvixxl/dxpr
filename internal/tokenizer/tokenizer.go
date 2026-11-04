// Package tokenizer converts an expression string into a stream of tokens.
package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
)

// Kind identifies the type of a token.
type Kind int

const (
	KindNumber Kind = iota
	KindDice   // 'd' or 'D'
	KindPlus
	KindMinus
	KindTimes       // '*'
	KindStarStar    // '**'
	KindSlash       // '/'
	KindDoubleSlash // '//'
	KindPercent     // '%'
	KindCaret       // '^'
	KindEOF
)

// Token is a single lexical unit of an expression.
type Token struct {
	Kind  Kind
	Value string // literal text (digits for KindNumber)
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
		case r == 'd' || r == 'D':
			toks = append(toks, Token{Kind: KindDice, Value: string(r)})
			i++
		case r == '+':
			toks = append(toks, Token{Kind: KindPlus, Value: "+"})
			i++
		case r == '-':
			toks = append(toks, Token{Kind: KindMinus, Value: "-"})
			i++
		case r == '*':
			if i+1 < len(runes) && runes[i+1] == '*' {
				toks = append(toks, Token{Kind: KindStarStar, Value: "**"})
				i += 2
			} else {
				toks = append(toks, Token{Kind: KindTimes, Value: "*"})
				i++
			}
		case r == '/':
			if i+1 < len(runes) && runes[i+1] == '/' {
				toks = append(toks, Token{Kind: KindDoubleSlash, Value: "//"})
				i += 2
			} else {
				toks = append(toks, Token{Kind: KindSlash, Value: "/"})
				i++
			}
		case r == '%':
			toks = append(toks, Token{Kind: KindPercent, Value: "%"})
			i++
		case r == '^':
			toks = append(toks, Token{Kind: KindCaret, Value: "^"})
			i++
		default:
			return nil, fmt.Errorf("unexpected character %q", string(r))
		}
	}

	toks = append(toks, Token{Kind: KindEOF})
	return toks, nil
}

// Describe returns a human-readable name for a token kind.
func (k Kind) Describe() string {
	switch k {
	case KindNumber:
		return "number"
	case KindDice:
		return "'d'"
	case KindPlus:
		return "'+'"
	case KindMinus:
		return "'-'"
	case KindTimes:
		return "'*'"
	case KindStarStar:
		return "'**'"
	case KindSlash:
		return "'/'"
	case KindDoubleSlash:
		return "'//'"
	case KindPercent:
		return "'%'"
	case KindCaret:
		return "'^'"
	case KindEOF:
		return "end of expression"
	default:
		return strings.ToLower(fmt.Sprintf("token(%d)", k))
	}
}
