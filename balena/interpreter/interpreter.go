package balena

import (
	"fmt"

	"github.com/mrdapoyo/dofi/balena/token"
)

type ExprVisitor interface {
	Visit(expr Expr) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
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

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
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
		if leftVal, ok := left.(float64); ok {
			if rightVal, ok := right.(float64); ok {
				return leftVal + rightVal
			}
			if leftStr, ok := left.(string); ok {
				if rightStr, ok := right.(string); ok {
					return leftStr + rightStr
				}
			}
		}
		break
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

	// unreachable
	return nil
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
	return a == b || (a != nil && b != nil && a == b)
}

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Interpreter struct{}

func (i *Interpreter) Visit(expr Expr) interface{} {
	return nil
}

func (i *Interpreter) Interpret(expression Expr) {
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