package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "A1+AZ209*     max(1, 2.5)",
			expected: []Token{
				{Type: CELL_ADDRESS, Literal: "A1"},
				{Type: PLUS, Literal: "+"},
				{Type: CELL_ADDRESS, Literal: "AZ209"},
				{Type: MULT, Literal: "*"},
				{Type: FUNCTION, Literal: "max"},
				{Type: LPAREN, Literal: "("},
				{Type: INTEGER, Literal: "1"},
				{Type: COMMA, Literal: ","},
				{Type: FLOAT, Literal: "2.5"},
				{Type: RPAREN, Literal: ")"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "B12 - 45.6 / sum(A1,B12)",
			expected: []Token{
				{Type: CELL_ADDRESS, Literal: "B12"},
				{Type: MINUS, Literal: "-"},
				{Type: FLOAT, Literal: "45.6"},
				{Type: DIV, Literal: "/"},
				{Type: FUNCTION, Literal: "sum"},
				{Type: LPAREN, Literal: "("},
				{Type: CELL_ADDRESS, Literal: "A1"},
				{Type: COMMA, Literal: ","},
				{Type: CELL_ADDRESS, Literal: "B12"},
				{Type: RPAREN, Literal: ")"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expectedToken := range tt.expected {
			token := l.NextToken()

			if token.Type != expectedToken.Type {
				t.Fatalf("test[%d] - token type wrong. expected=%q, got=%q", i, expectedToken.Type, token.Type)
			}

			if token.Literal != expectedToken.Literal {
				t.Fatalf("test[%d] - token literal wrong. expected=%q, got=%q", i, expectedToken.Literal, token.Literal)
			}
		}
	}
}
