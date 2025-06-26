package balena

import (
	"github.com/mrdapoyo/dofi/balena/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

type ParseError struct {
	Line     int
	Message  string
	Location string
}

func (e ParseError) Error() string {
	return e.Message
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() Expr {
	return p.equality()
}
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(token.FALSE) {
		return &LiteralExpr{Value: false}
	}
	if p.match(token.TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(token.NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{Expression: expr}
	}

	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) consume(t token.TokenType, message string) token.Token {
	if p.check(t) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) error(t token.Token, message string) ParseError {
	var Error = ParseError{
		Line:     t.Line,
		Message:  message,
		Location: t.Lexeme,
	}
	// balena.error(ParseError)
	return Error
}

func (p *Parser) Parse() Expr {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
			} else {
				panic(r)
			}
		}
	}()
	return p.expression()
}