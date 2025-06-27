package balena

import (
	"fmt"

	environment "github.com/mrdapoyo/dofi/balena/env"
	parser "github.com/mrdapoyo/dofi/balena/parser"
	"github.com/mrdapoyo/dofi/balena/token"
)

type ExprVisitor interface {
	VisitLiteralExpr(expr *parser.LiteralExpr) interface{}
	VisitGroupingExpr(expr *parser.GroupingExpr) interface{}
	VisitUnaryExpr(expr *parser.UnaryExpr) interface{}
	VisitBinaryExpr(expr *parser.BinaryExpr) interface{}
	VisitVariableExpr(expr *parser.VariableExpr) interface{}
	VisitAssignExpr(expr *parser.AssignExpr) interface{}
	VisitLogicalExpr(expr *parser.LogicalExpr) interface{}
	VisitCallExpr(expr *parser.CallExpr) interface{}
}

type Interpreter struct {
	environment *environment.Environment
}

func (i *Interpreter) VisitExpressionStmt(stmt *parser.ExpressionStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitFunctionStmt(stmt *parser.FunctionStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitIfStmt(stmt *parser.IfStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitVarStmt(stmt *parser.VarStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Set(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitReturnStmt(stmt *parser.ReturnStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitWhileStmt(stmt *parser.WhileStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.LiteralExpr) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.GroupingExpr) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.UnaryExpr) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case token.BANG:
		if right == nil {
			return true
		}
		if val, ok := right.(bool); ok {
			return !val
		}
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *parser.VariableExpr) interface{} {
	variable, _ := i.environment.Get(expr.Name.Lexeme)
	return variable
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.BinaryExpr) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case token.PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
	case token.GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *parser.AssignExpr) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitLogicalExpr(expr *parser.LogicalExpr) interface{} {
	// TODO: Implement logical expressions
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *parser.CallExpr) interface{} {
	// TODO: Implement function calls
	return nil
}

func (i *Interpreter) evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(statements []parser.Stmt) {
	for _, statement := range statements {
		statement.Accept(i)
	}
}

func (i *Interpreter) Interpret(expression parser.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				LoxRuntimeError(runtimeErr)
			} else {
				panic(r)
			}
		}
	}()
	value := i.evaluate(expression)
	println(stringify(value))
}

func (i *Interpreter) VisitBlockStmt(stmt *parser.BlockStmt) {
	previous := i.environment
	i.environment = environment.NewEnvironment(previous)
	defer func() { i.environment = previous }()
	i.execute(stmt.Statements)
}

func (i *Interpreter) executeBlock(statements []parser.Stmt, env *environment.Environment) {
	previous := i.environment
	i.environment = env
	defer func() { i.environment = previous }()
	i.execute(statements)
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	switch v := object.(type) {
	case float64:
		s := fmt.Sprintf("%g", v)
		return s
	default:
		return fmt.Sprintf("%v", object)
	}
}

func LoxRuntimeError(err *RuntimeError) {
	println(err.Error())
}

func checkNumberOperands(operator token.Token, left, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(&RuntimeError{
		Token:   operator,
		Message: "Operands must be two numbers or two strings.",
	})
}

func checkNumberOperand(operator token.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(&RuntimeError{
		Token:   operator,
		Message: "Operand must be a number.",
	})
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}
