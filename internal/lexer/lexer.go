package lexer

import (
	"unicode"
)

// TokenType defines the type of token
type TokenType string

const (
	INTEGER      TokenType = "INTEGER"
	FLOAT        TokenType = "FLOAT"
	STRING       TokenType = "STRING"
	CELL_ADDRESS TokenType = "CELL_ADDRESS"
	FUNCTION     TokenType = "FUNCTION"
	LPAREN       TokenType = "LPAREN"
	RPAREN       TokenType = "RPAREN"
	COMMA        TokenType = "COMMA"
	EOF          TokenType = "EOF"
	PLUS         TokenType = "PLUS"
	MINUS        TokenType = "MINUS"
	MULT         TokenType = "MULT"
	DIV          TokenType = "DIV"
)

// Token represents a token in the formula
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer holds the state for lexing
type Lexer struct {
	input   string
	pos     int
	current rune
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character in the input
func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.current = 0 // End of input
	} else {
		l.current = rune(l.input[l.pos])
	}
	l.pos++
}

// NextToken processes and returns the next token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.current == 0 {
		return Token{Type: EOF, Literal: ""}
	}

	switch l.current {
	case '(':
		l.readChar()
		return Token{Type: LPAREN, Literal: "("}
	case '+':
		l.readChar()
		return Token{Type: PLUS, Literal: "+"}
	case '-':
		l.readChar()
		return Token{Type: MINUS, Literal: "-"}
	case '*':
		l.readChar()
		return Token{Type: MULT, Literal: "*"}
	case '/':
		l.readChar()
		return Token{Type: DIV, Literal: "/"}
	case ')':
		l.readChar()
		return Token{Type: RPAREN, Literal: ")"}
	case ',':
		l.readChar()
		return Token{Type: COMMA, Literal: ","}
	}

	if unicode.IsDigit(l.current) {
		return l.readNumber()
	}

	if unicode.IsLetter(l.current) {
		return l.readIdentifierOrCell()
	}

	// Return as string by default if nothing else matches
	literal := l.readString()
	return Token{Type: STRING, Literal: literal}
}

// readNumber reads an integer or float token
func (l *Lexer) readNumber() Token {
	startPos := l.pos - 1
	isFloat := false

	for unicode.IsDigit(l.current) || l.current == '.' {
		if l.current == '.' {
			isFloat = true
		}
		l.readChar()
	}

	if isFloat {
		return Token{Type: FLOAT, Literal: l.input[startPos : l.pos-1]}
	}
	return Token{Type: INTEGER, Literal: l.input[startPos : l.pos-1]}
}

// readIdentifierOrCell reads either a cell address or a function identifier
func (l *Lexer) readIdentifierOrCell() Token {
	startPos := l.pos - 1

	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) {
		l.readChar()
	}

	literal := l.input[startPos : l.pos-1]

	// If the next character is '(', it's a function call
	if l.current == '(' {
		return Token{Type: FUNCTION, Literal: literal}
	}

	// Otherwise, it's a cell address
	return Token{Type: CELL_ADDRESS, Literal: literal}
}

// readString reads a string (sequence of letters/numbers without quotes)
func (l *Lexer) readString() string {
	startPos := l.pos - 1

	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) {
		l.readChar()
	}

	return l.input[startPos : l.pos-1]
}

// skipWhitespace skips any whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.current) {
		l.readChar()
	}
}
