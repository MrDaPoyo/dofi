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

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return &VarStmt{
		Name:        name,
		Initializer: initializer,
	}
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

	if p.match(token.IDENTIFIER) {
		return &VariableExpr{Name: p.previous()}
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

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) interface{}
	VisitGroupingExpr(expr *GroupingExpr) interface{}
	VisitLiteralExpr(expr *LiteralExpr) interface{}
	VisitUnaryExpr(expr *UnaryExpr) interface{}
	VisitVariableExpr(expr *VariableExpr) interface{}
	VisitAssignExpr(expr *AssignExpr) interface{}
	VisitLogicalExpr(expr *LogicalExpr) interface{}
	VisitCallExpr(expr *CallExpr) interface{}
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name token.Token
}

func (e *VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  token.Token
	Value Expr
}

func (e *AssignExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(e)
}

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *LogicalExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(e)
}

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (e *CallExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(e)
}

// Stmt interface and statement types

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitPrintStmt(stmt *PrintStmt)
	VisitVarStmt(stmt *VarStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitWhileStmt(stmt *WhileStmt)
	VisitFunctionStmt(stmt *FunctionStmt)
	VisitReturnStmt(stmt *ReturnStmt)
}

type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(s)
}

type VarStmt struct {
	Name        token.Token
	Initializer Expr
}

func (s *VarStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(visitor StmtVisitor) {
	visitor.VisitBlockStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s *IfStmt) Accept(visitor StmtVisitor) {
	visitor.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) Accept(visitor StmtVisitor) {
	visitor.VisitWhileStmt(s)
}

type FunctionStmt struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (s *FunctionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitFunctionStmt(s)
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func (s *ReturnStmt) Accept(visitor StmtVisitor) {
	visitor.VisitReturnStmt(s)
}

func (p *Parser) statement() Stmt {
	// TODO: Implement statement parsing (print, if, while, block, etc.)
	return &ExpressionStmt{Expression: p.expression()}
}
