package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Token
	}{
		{
			"single number",
			"42",
			[]Token{{KindNumber, "42"}, {KindEOF, ""}},
		},
		{
			"dice without count",
			"d20",
			[]Token{{KindDice, "d"}, {KindNumber, "20"}, {KindEOF, ""}},
		},
		{
			"dice with count, uppercase D",
			"3D6",
			[]Token{{KindNumber, "3"}, {KindDice, "D"}, {KindNumber, "6"}, {KindEOF, ""}},
		},
		{
			"all operators",
			"1+2-3*4/5%6^7",
			[]Token{
				{KindNumber, "1"}, {KindOperator, "+"},
				{KindNumber, "2"}, {KindOperator, "-"},
				{KindNumber, "3"}, {KindOperator, "*"},
				{KindNumber, "4"}, {KindOperator, "/"},
				{KindNumber, "5"}, {KindOperator, "%"},
				{KindNumber, "6"}, {KindOperator, "^"},
				{KindNumber, "7"}, {KindEOF, ""},
			},
		},
		{
			"multi-char operators longest match",
			"2**3//4",
			[]Token{
				{KindNumber, "2"}, {KindOperator, "**"},
				{KindNumber, "3"}, {KindOperator, "//"},
				{KindNumber, "4"}, {KindEOF, ""},
			},
		},
		{
			"adjacent single-char operators split",
			"2* *3",
			[]Token{
				{KindNumber, "2"}, {KindOperator, "*"},
				{KindOperator, "*"},
				{KindNumber, "3"}, {KindEOF, ""},
			},
		},
		{
			"whitespace ignored",
			"  d20\t+\n5  ",
			[]Token{
				{KindDice, "d"}, {KindNumber, "20"},
				{KindOperator, "+"}, {KindNumber, "5"}, {KindEOF, ""},
			},
		},
		{
			"empty expression",
			"",
			[]Token{{KindEOF, ""}},
		},
		{
			"d followed by digit is dice",
			"d20",
			[]Token{{KindDice, "d"}, {KindNumber, "20"}, {KindEOF, ""}},
		},
		{
			"d followed by letter is identifier",
			"4d6dl1",
			[]Token{
				{KindNumber, "4"}, {KindDice, "d"}, {KindNumber, "6"},
				{KindIdent, "dl"}, {KindNumber, "1"}, {KindEOF, ""},
			},
		},
		{
			"identifier with digits split",
			"kh10",
			[]Token{{KindIdent, "kh"}, {KindNumber, "10"}, {KindEOF, ""}},
		},
		{
			"letter run containing alias letter stays identifier",
			"max(1,2)",
			[]Token{
				{KindIdent, "max"}, {KindLParen, "("}, {KindNumber, "1"},
				{KindComma, ","}, {KindNumber, "2"}, {KindRParen, ")"}, {KindEOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.in)
			if err != nil {
				t.Fatalf("Tokenize(%q): %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	for _, in := range []string{"d20!", "3 + $"} {
		if _, err := Tokenize(in); err == nil {
			t.Errorf("Tokenize(%q) succeeded, want error", in)
		}
	}
}

func TestTokenDescribe(t *testing.T) {
	tests := []struct {
		tok  Token
		want string
	}{
		{Token{KindNumber, "42"}, "number"},
		{Token{KindDice, "d"}, "'d'"},
		{Token{KindOperator, "//"}, "'//'"},
		{Token{KindEOF, ""}, "end of expression"},
	}
	for _, tt := range tests {
		if got := tt.tok.Describe(); got != tt.want {
			t.Errorf("Token(%v).Describe() = %q, want %q", tt.tok, got, tt.want)
		}
	}
}
