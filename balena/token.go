package balena

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"   // Identifiers + literals
	IDENT   = "IDENT" // add, foobar, x, y, ...
	INT     = "INT"   // 1343456
	FLOAT   = "FLOAT" // 3.14159
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	PLUS_EQ  = "+="
	MINUS_EQ = "-="
	EQ       = "=="
	NOT_EQ   = "!="
	// Delimiters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]" // Keywords
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	WHILE     = "WHILE"
	FOR       = "FOR"

	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	STRING = "STRING"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
