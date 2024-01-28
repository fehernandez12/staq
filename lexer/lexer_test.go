package lexer

import (
	"fmt"
	"testing"

	"staq/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10.0;
	let add = fn(x, y) {
	x + y;
	};
	let result = add(five, ten);
	let someFloat = 3.1415;
	let someHex = 0xdeadbeef;
	let someOct = 0o1234567;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10.0"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "someFloat"},
		{token.ASSIGN, "="},
		{token.NUMBER, "3.1415"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "someHex"},
		{token.ASSIGN, "="},
		{token.NUMBER, "0xdeadbeef"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "someOct"},
		{token.ASSIGN, "="},
		{token.NUMBER, "0o1234567"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Printf("index: %d, type: %v\n", i, tok.Type)

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
