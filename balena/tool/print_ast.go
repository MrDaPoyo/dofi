package balena

import (
	"fmt"

	parser "github.com/mrdapoyo/dofi/balena/parser"
)

type AstPrinter struct {
	lastResult string
}

func (a *AstPrinter) VisitLogicalExpr(expr *parser.LogicalExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitVariableExpr(expr *parser.VariableExpr) interface{} {
	return expr.Name.Lexeme
}

func (a *AstPrinter) Print(expr parser.Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) VisitBinaryExpr(expr *parser.BinaryExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *parser.GroupingExpr) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *parser.LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitUnaryExpr(expr *parser.UnaryExpr) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...parser.Expr) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(a).(string)
	}
	result += ")"
	return result
}

func (a *AstPrinter) VisitAssignExpr(expr *parser.AssignExpr) interface{} {
	return a.parenthesize("assign "+expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitCallExpr(expr *parser.CallExpr) interface{} {
	result := a.parenthesize("call " + expr.Callee.Accept(a).(string))
	for _, arg := range expr.Arguments {
		result += " " + arg.Accept(a).(string)
	}
	result += ")"
	return result
}

func (a *AstPrinter) PrintStmt(stmt parser.Stmt) string {
	stmt.Accept(a)
	return a.lastResult
}

func (a *AstPrinter) VisitBlockStmt(stmt *parser.BlockStmt) {
	result := "(block"
	for _, s := range stmt.Statements {
		result += " " + a.PrintStmt(s)
	}
	result += ")"
	a.lastResult = result
}

func (a *AstPrinter) VisitExpressionStmt(stmt *parser.ExpressionStmt) {
	a.lastResult = a.parenthesize2(";", stmt.Expression)
}

func (a *AstPrinter) VisitFunctionStmt(stmt *parser.FunctionStmt) {
	result := "(fun " + stmt.Name.Lexeme + "("
	for i, param := range stmt.Params {
		if i > 0 {
			result += " "
		}
		result += param.Lexeme
	}
	result += ") "
	for _, body := range stmt.Body {
		result += a.PrintStmt(body)
	}
	result += ")"
	a.lastResult = result
}

func (a *AstPrinter) VisitIfStmt(stmt *parser.IfStmt) {
	if stmt.ElseBranch == nil {
		a.lastResult = a.parenthesize2("if", stmt.Condition, stmt.ThenBranch)
	} else {
		a.lastResult = a.parenthesize2("if-else", stmt.Condition, stmt.ThenBranch, stmt.ElseBranch)
	}
}

func (a *AstPrinter) VisitPrintStmt(stmt *parser.PrintStmt) {
	a.lastResult = a.parenthesize2("print", stmt.Expression)
}

func (a *AstPrinter) VisitReturnStmt(stmt *parser.ReturnStmt) {
	if stmt.Value == nil {
		a.lastResult = "(return)"
	} else {
		a.lastResult = a.parenthesize2("return", stmt.Value)
	}
}

func (a *AstPrinter) VisitVarStmt(stmt *parser.VarStmt) {
	if stmt.Initializer == nil {
		a.lastResult = a.parenthesize2("var", stmt.Name.Lexeme)
	} else {
		a.lastResult = a.parenthesize2("var", stmt.Name.Lexeme, "=", stmt.Initializer)
	}
}

func (a *AstPrinter) VisitWhileStmt(stmt *parser.WhileStmt) {
	a.lastResult = a.parenthesize2("while", stmt.Condition, stmt.Body)
}

func (a *AstPrinter) parenthesize2(name string, parts ...interface{}) string {
	result := "(" + name
	result += a.transform(parts...)
	result += ")"
	return result
}

func (a *AstPrinter) transform(parts ...interface{}) string {
	result := ""
	for _, part := range parts {
		switch v := part.(type) {
		case parser.Expr:
			result += " " + v.Accept(a).(string)
		case parser.Stmt:
			result += " " + a.PrintStmt(v)
		case string:
			result += " " + v
		case []parser.Expr:
			for _, e := range v {
				result += " " + e.Accept(a).(string)
			}
		case []parser.Stmt:
			for _, s := range v {
				result += " " + a.PrintStmt(s)
			}
		default:
			if v == nil {
				continue
			}
			result += fmt.Sprintf(" %v", v)
		}
	}
	return result
}
