package main

import "fmt"

var (
	EVAL_NULL  = &ObjectNull{}
	EVAL_TRUE  = &ObjectBoolean{Value: true}
	EVAL_FALSE = &ObjectBoolean{Value: false}
)

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalProgram(node)
	case *LetStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
	case *ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *IntegerLiteral:
		return &ObjectInteger{Value: node.Value}
	case *Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatement(node)
	case *IfExpression:
		return evalIfExpression(node)
	case *ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	}

	return nil
}

func evalStatements(stmts []Statement) Object {
	var result Object
	for _, statement := range stmts {
		result = Eval(statement)
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalIfExpression(ie *IfExpression) Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return EVAL_NULL
	}
}

func isTruthy(obj Object) bool {
	switch obj {
	case EVAL_NULL:
		return false
	case EVAL_TRUE:
		return true
	case EVAL_FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *ObjectBoolean {
	if input {
		return EVAL_TRUE
	}
	return EVAL_FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*ObjectInteger).Value
	return &ObjectInteger{Value: -value}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case EVAL_TRUE:
		return EVAL_FALSE
	case EVAL_FALSE:
		return EVAL_TRUE
	case EVAL_NULL:
		return EVAL_TRUE
	default:
		return EVAL_FALSE
	}
}

func evalInfixExpression(
	operator string,
	left, right Object,
) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right Object,
) Object {
	leftVal := left.(*ObjectInteger).Value
	rightVal := right.(*ObjectInteger).Value
	switch operator {
	case "+":
		return &ObjectInteger{Value: leftVal + rightVal}
	case "-":
		return &ObjectInteger{Value: leftVal - rightVal}
	case "*":
		return &ObjectInteger{Value: leftVal * rightVal}
	case "/":
		return &ObjectInteger{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalProgram(program *Program) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *BlockStatement) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
