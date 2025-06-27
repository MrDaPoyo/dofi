package balena

import (
	parser "github.com/mrdapoyo/dofi/balena/parser"
	token "github.com/mrdapoyo/dofi/balena/token"
)

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *parser.BlockStmt) {
	r.BeginScope()
	r.Resolve(stmt.Statements)
	r.EndScope()
}

func (r *Resolver) Resolve(statements []parser.Stmt) {
	for _, statement := range statements {
		r.ResolveStmt(statement)
	}
}

func (r *Resolver) ResolveStmt(stmt parser.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) ResolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) BeginScope() {
	if r.scopes == nil {
		r.scopes = []map[string]bool{}
	}
	scope := make(map[string]bool)
	r.scopes = append(r.scopes, scope)
}

func (r *Resolver) EndScope() {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitExpressionStmt(stmt *parser.ExpressionStmt) {
	r.ResolveExpr(stmt.Expression)
}

func (r *Resolver) VisitVarStmt(stmt *parser.VarStmt) {
	r.Declare(stmt.Name)
	if stmt.Initializer != nil {
		r.ResolveExpr(stmt.Initializer)
	}
	r.Define(stmt.Name)
}

func (r *Resolver) Declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = false
}

func (r *Resolver) Define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) VisitVariableExpr(expr *parser.VariableExpr) interface{} {
	return r.lookUpVariable(expr.Name, expr)
}

func (r *Resolver) lookUpVariable(name token.Token, expr parser.Expr) interface{} {
	r.resolveLocal(expr, name)
	return nil
}

func (r *Resolver) resolveLocal(expr parser.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitAssignExpr(expr *parser.AssignExpr) interface{} {
	r.ResolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *parser.FunctionStmt) {
	r.Declare(stmt.Name)
	r.Define(stmt.Name)

	r.resolveFunction(stmt)
}

func (r *Resolver) VisitIfStmt(stmt *parser.IfStmt) {
	r.ResolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.ResolveStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) resolveFunction(stmt *parser.FunctionStmt) {
	r.BeginScope()
	for _, param := range stmt.Params {
		r.Declare(param)
		r.Define(param)
	}
	r.Resolve(stmt.Body)
	r.EndScope()
}

func (r *Resolver) VisitPrintStmt(stmt *parser.PrintStmt) {
	r.ResolveExpr(stmt.Expression)
}

func (r *Resolver) VisitCallExpr(expr *parser.CallExpr) interface{} {
	r.ResolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.ResolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *parser.GroupingExpr) interface{} {
	r.ResolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *parser.LiteralExpr) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *parser.LogicalExpr) interface{} {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *parser.UnaryExpr) interface{} {
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *parser.BinaryExpr) interface{} {
	r.ResolveExpr(expr.Left)
	r.ResolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *parser.ReturnStmt) {
	if stmt.Value != nil {
		r.ResolveExpr(stmt.Value)
	}
}

func (r *Resolver) VisitWhileStmt(stmt *parser.WhileStmt) {
	r.ResolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.Body)
}
