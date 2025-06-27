package balena

import (
	"fmt"
	"time"

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
	Globals     *environment.Environment
	Environment *environment.Environment
	locals      map[parser.Expr]int
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		Globals:     environment.NewEnvironment(),
		Environment: environment.NewEnvironment(),
	}
	interpreter.Globals.Set("clock", &Clock{})
	interpreter.Environment = interpreter.Globals
	return interpreter
}

func (in *Interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if v, ok := val.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) VisitExpressionStmt(stmt *parser.ExpressionStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitFunctionStmt(stmt *parser.FunctionStmt) {
	function := NewBalenaFunction(stmt, i.Environment)
	i.Environment.Set(stmt.Name.Lexeme, function)
}

func (i *Interpreter) VisitIfStmt(stmt *parser.IfStmt) {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute([]parser.Stmt{stmt.ThenBranch})
	} else if stmt.ElseBranch != nil {
		i.execute([]parser.Stmt{stmt.ElseBranch})
	}
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) {
	panic("unimplemented")
}

func (i *Interpreter) VisitVarStmt(stmt *parser.VarStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.Environment.Set(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitReturnStmt(stmt *parser.ReturnStmt) {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(NewReturn(value))
}

func (i *Interpreter) VisitWhileStmt(stmt *parser.WhileStmt) {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute([]parser.Stmt{stmt.Body})
	}
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.LiteralExpr) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *parser.LogicalExpr) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
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
	variable, _ := i.Environment.Get(expr.Name.Lexeme)
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
	if i.locals != nil {
		if distance, ok := i.locals[expr]; ok {
			i.Environment.AssignAt(distance, expr.Name, value)
			return value
		}
	}
	i.Globals.Assign(expr.Name, value)
	i.Environment.Assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitCallExpr(expr *parser.CallExpr) interface{} {
	callee := i.evaluate(expr.Callee)
	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := callee.(BalenaCallable)
	if !ok {
		panic(&RuntimeError{
			Token:   expr.Paren,
			Message: "Can only call functions and classes.",
		})
	}
	if len(expr.Arguments) != function.Arity() {
		panic(&RuntimeError{
			Token:   expr.Paren,
			Message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(expr.Arguments)),
		})
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(statements []parser.Stmt) {
	for _, statement := range statements {
		statement.Accept(i)
	}
}

func (i *Interpreter) Resolve(expr parser.Expr, depth int) {
	if i.locals == nil {
		i.locals = make(map[parser.Expr]int)
	}
	i.locals[expr] = depth
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
	previous := i.Environment
	i.Environment = environment.NewEnvironment(previous)
	defer func() { i.Environment = previous }()
	i.execute(stmt.Statements)
}

func (i *Interpreter) executeBlock(statements []parser.Stmt, env *environment.Environment) {
	previous := i.Environment
	i.Environment = env
	defer func() { i.Environment = previous }()
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

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixNano()) / 1e9
}

func (c *Clock) String() string {
	return "<native fn>"
}
