package balena

import (
	"strconv"

	"github.com/mrdapoyo/dofi/balena/token"
)

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: make([]token.Token, 0),
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.Token{Type: token.EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	ch := rune(s.source[s.current])
	s.current++
	return ch
}

func (s *Scanner) addToken(t token.TokenType, lexeme string, literal interface{}) {
	if lexeme == "" {
		lexeme = s.source[s.start:s.current]
	}
	s.tokens = append(s.tokens, token.Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) scanToken() {
	switch r := s.advance(); r {
	case ' ', '\r', '\t':
		// get rekt
		break
	case '\n':
		s.line++
		break
	case '(':
		s.addToken(token.LEFT_PAREN, "", nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, "", nil)
	case '{':
		s.addToken(token.LEFT_BRACE, "", nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, "", nil)
	case ',':
		s.addToken(token.COMMA, "", nil)
	case '.':
		s.addToken(token.DOT, "", nil)
	case '-':
		s.addToken(token.MINUS, "", nil)
	case '+':
		s.addToken(token.PLUS, "", nil)
	case ';':
		s.addToken(token.SEMICOLON, "", nil)
	case '*':
		s.addToken(token.STAR, "", nil)
	case '%':
		// not defined maybe treat as STAR? for now ignore or add? we'll ignore break.
		s.addToken(token.STAR, "", nil)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, "", nil)
		} else {
			s.addToken(token.BANG, "", nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, "", nil)
		} else {
			s.addToken(token.EQUAL, "", nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, "", nil)
		} else {
			s.addToken(token.LESS, "", nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, "", nil)
		} else {
			s.addToken(token.GREATER, "", nil)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && rune(s.source[s.current]) != '\n' {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, "", nil)
		}
	case '"':
		s.newString()
	default:
		if s.isDigit(r) {
			s.newNumber()
		} else if s.isAlpha(r) {
			s.newIdentifier()
		} else {
			// balena.Error(s.line, "Unexpected character: "+string(r))
		}
	}
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) newIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	lexeme := s.source[s.start:s.current]
	typ, ok := keywords[lexeme]
	if !ok {
		typ = token.IDENTIFIER
	}
	s.addToken(typ, lexeme, lexeme)
}

func (s *Scanner) newString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		// balena.Error(s.line, "Unterminated string.")
		return
	}

	// the closing "
	s.advance()

	// just the juice of the fruit, not the quotes
	lexeme := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, lexeme, lexeme)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) newNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	// fractional part? maybe later
	literalStr := s.source[s.start:s.current]
	var numVal float64
	if parsed, err := strconv.ParseFloat(literalStr, 64); err == nil {
		numVal = parsed
	} else {
		numVal = 0
	}
	s.addToken(token.NUMBER, literalStr, numVal)
}

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"fn":     token.FUN,
	"for":    token.FOR,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}
