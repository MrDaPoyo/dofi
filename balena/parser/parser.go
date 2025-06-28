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
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*VariableExpr); ok {
			name := variable.Name
			return &AssignExpr{
				Name:  name,
				Value: value,
			}
		}

		panic(p.error(equals, "Invalid assignment target."))
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
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

	if p.match(token.FUN) {
		return p.function()
	}

	return p.statement()
}

func (p *Parser) function() Stmt {
	name := p.consume(token.IDENTIFIER, "Expect function name.")
	p.consume(token.LEFT_PAREN, "Expect '(' after function name.")
	params := []token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				panic(p.error(p.peek(), "Can't have more than 255 parameters."))
			}
			params = append(params, p.consume(token.IDENTIFIER, "Expect parameter name."))
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(token.LEFT_BRACE, "Expect '{' before function body.")
	body := p.blockStatement()
	block, ok := body.(*BlockStmt)
	if !ok {
		panic("Expected block statement for function body")
	}
	return &FunctionStmt{Name: name, Params: params, Body: block.Statements}
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

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := []Expr{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				panic(p.error(p.peek(), "Can't have more than 255 arguments."))
			}
			arguments = append(arguments, p.expression())
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	return &CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
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

	if p.match(token.LEFT_BRACKET) {
		// Array literal: [expr, expr, ...]
		elements := []Expr{}
		if !p.check(token.RIGHT_BRACKET) {
			for {
				elements = append(elements, p.expression())
				if !p.match(token.COMMA) {
					break
				}
			}
		}
		p.consume(token.RIGHT_BRACKET, "Expect ']' after array elements.")
		return &ArrayExpr{Elements: elements}
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
	VisitArrayExpr(expr *ArrayExpr) interface{}
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

type ArrayExpr struct {
	Elements []Expr
}

func (e *ArrayExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitArrayExpr(e)
}

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

func (s *FunctionStmt) Statements() []Stmt {
	return s.Body
}

type ReturnStmt struct {
	Keyword token.Token
	Value   Expr
}

func (s *ReturnStmt) Accept(visitor StmtVisitor) {
	visitor.VisitReturnStmt(s)
}

func (p *Parser) statement() Stmt {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return p.blockStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.RETURN) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return &PrintStmt{Expression: value}
}

func (p *Parser) blockStatement() Stmt {
	statements := []Stmt{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return &BlockStmt{Statements: statements}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	// Create the while loop body that includes both the original body and the increment
	var whileBody Stmt = body
	if increment != nil {
		whileBody = &BlockStmt{
			Statements: []Stmt{
				body,
				&ExpressionStmt{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{Value: true}
	}
	body = &WhileStmt{
		Condition: condition,
		Body:      whileBody,
	}

	if initializer != nil {
		body = &BlockStmt{
			Statements: []Stmt{
				initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return &ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return &ExpressionStmt{Expression: expr}
}

// Optionally define BreakStmt and ContinueStmt types if not already present:
type BreakStmt struct{}

func (s *BreakStmt) Accept(visitor StmtVisitor) { /* implement if needed */ }

type ContinueStmt struct{}

func (s *ContinueStmt) Accept(visitor StmtVisitor) { /* implement if needed */ }
