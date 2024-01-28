package lexer

import (
	"errors"
	"staq/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken is the main function of the lexer.
// It returns the next token in the input string.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.ADDASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SUBASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DIVASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '/' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INTDIV, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MULASSIGN, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EXP, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SHL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GEQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SHR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '?':
		if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NULLCOAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BITAND, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BITOR, l.ch)
		}
	case '^':
		tok = newToken(token.BITXOR, l.ch)
	case '~':
		tok = newToken(token.BITNOT, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:

		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tokLiteral, tokType, err := l.readNumber()
			if err != nil {
				tok = newToken(token.ILLEGAL, l.ch)
			}
			tok.Literal = tokLiteral
			tok.Type = tokType
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number from the input string, including floats.
// Fails is the number is malformed (i.e. 3.1415.92),
// in that case it returns an error.
func (l *Lexer) readNumber() (string, token.TokenType, error) {
	var tokType token.TokenType
	position := l.position
	dotCount := 0
	tokType = token.INT

	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			dotCount++
			if dotCount > 1 {
				tokType = token.ILLEGAL
				// Gather the entire literal
				for isDigit(l.ch) || l.ch == '.' {
					l.readChar()
				}
				// Return the literal, the token type and an error
				return l.input[position:l.position], tokType, errors.New("malformed number")
			}
			tokType = token.FLOAT
		}
		l.readChar()
	}

	return l.input[position:l.position], tokType, nil
}

// peekChar returns the next character in the input string without advancing the
// lexer's position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
