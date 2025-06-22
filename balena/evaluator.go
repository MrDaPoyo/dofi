package main

var (
	EVAL_NULL  = &ObjectNull{}
	EVAL_TRUE  = &ObjectBoolean{Value: true}
	EVAL_FALSE = &ObjectBoolean{Value: false}
)

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalProgram(node)
	case *ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *IntegerLiteral:
		return &ObjectInteger{Value: node.Value}
	case *Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatement(node)
	case *IfExpression:
		return evalIfExpression(node)
	case *ReturnStatement:
		val := Eval(node.ReturnValue)
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
		return EVAL_NULL
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return EVAL_NULL
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
	default:
		return EVAL_NULL
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
		return EVAL_NULL
	}
}

func evalProgram(program *Program) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatement(block *BlockStatement) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}
