package lexer

import (
	"testing"

	"staq/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOtherTokens(t *testing.T) {
	input := `!-/*5;
	10 < 20 > 10;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.LT, "<"},
		{token.INT, "20"},
		{token.GT, ">"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := `if (5 < 10) {
		return true;
	} else {
		return false;
	}`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestCompositeOperators(t *testing.T) {
	input := `10 == 10;
	10 != 9;
	10 >= 9;
	10 <= 9;
	10++;
	10--;
	10 ** 2;
	10 // 2;
	10 % 2;
	10 << 2;
	10 >> 2;
	10 += 2;
	10 -= 2;
	10 *= 2;
	10 /= 2;
	10 ?? 2;
	10 || 2;
	10 && 2;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.GEQ, ">="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.LEQ, "<="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.EXP, "**"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.INTDIV, "//"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.MOD, "%"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.SHL, "<<"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.SHR, ">>"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.ADDASSIGN, "+="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.SUBASSIGN, "-="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.MULASSIGN, "*="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.DIVASSIGN, "/="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NULLCOAL, "??"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.OR, "||"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.AND, "&&"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestBinaryOperators(t *testing.T) {
	input := `10 & 2;
	10 | 2;
	10 ^ 2;
	~10;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.BITAND, "&"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.BITOR, "|"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.BITXOR, "^"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.BITNOT, "~"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestInvalidNumber(t *testing.T) {
	input := `3.1415.1;`

	l := New(input)

	tok := l.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Fatalf("tokentype wrong. expected=%q, got=%q",
			token.ILLEGAL, tok.Type)
	}

	if tok.Literal != "3.1415.1" {
		t.Fatalf("literal wrong. expected=%q, got=%q",
			"3.1415.1", tok.Literal)
	}
}
