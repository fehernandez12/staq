package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="
	GEQ       = ">="
	LEQ       = "<="
	INC       = "++"
	DEC       = "--"
	EXP       = "**"
	INTDIV    = "//"
	MOD       = "%"
	SHL       = "<<"
	SHR       = ">>"
	OR        = "||"
	AND       = "&&"
	ADDASSIGN = "+="
	SUBASSIGN = "-="
	MULASSIGN = "*="
	DIVASSIGN = "/="
	NULLCOAL  = "??"
	BITAND    = "&"
	BITOR     = "|"
	BITXOR    = "^"
	BITNOT    = "~"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Strings
	QUOTE = "\""

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks the keywords table to see whether the given identifier is
// a keyword. If it is, it returns the keyword's TokenType constant. If it is
// not, it returns the IDENT TokenType constant.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
