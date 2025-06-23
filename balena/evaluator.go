package balena

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

func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node, env)
	case *LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ExpressionStatement:
		return Eval(node.Expression, env)
	case *IntegerLiteral:
		return &ObjectInteger{Value: node.Value}
	case *FloatLiteral:
		return &ObjectFloat{Value: node.Value}
	case *Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *BlockStatement:
		return evalBlockStatement(node, env)
	case *IfExpression:
		return evalIfExpression(node, env)
	case *ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *Identifier:
		return evalIdentifier(node, env)
	case *StringLiteral:
		return &String{Value: node.Value}
	case *ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	case *IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}
	case *CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *WhileStatement:
		return evalWhileStatement(node, env)
	case *ForStatement:
		return evalForStatement(node, env)
	case *AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *CompoundAssignmentStatement:
		return evalCompoundAssignmentStatement(node, env)
	case *GoObjectLiteral:
		return &GoObject{Value: node.Value}
	}

	return nil
}

func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		if len(fn.Parameters) != len(args) {
			return NewError("wrong number of arguments. got=%d, want=%d",
				len(args), len(fn.Parameters))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return NewError("not a function: %s", fn.Type())
	}
}

func HasFunction(env *Environment, name string) bool {
	obj, ok := env.Get(name)
	if ok {
		_, isFunc := obj.(*Function)
		return isFunc
	}
	_, exists := builtins[name]
	return exists
}

func CallFunction(env *Environment, name string, args ...Object) Object {
	obj, ok := env.Get(name)
	if ok {
		fn, ok := obj.(*Function)
		if !ok {
			return NewError("%s is not a function", name)
		}
		return applyFunction(fn, args)
	}
	builtin, exists := builtins[name]
	if !exists {
		return NewError("function not found: %s", name)
	}
	return applyFunction(builtin, args)
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIfExpression(ie *IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return EVAL_NULL
	}
}

func evalExpressions(
	exps []Expression,
	env *Environment,
) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(
	node *Identifier,
	env *Environment,
) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return NewError("identifier not found: %s", node.Value)
}

func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return NewError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*ObjectInteger).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return EVAL_NULL
	}
	return arrayObject.Elements[idx]
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
		return NewError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return NewError("unknown operator: -%s", right.Type())
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
		if left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ {
			leftVal := left.(*ObjectInteger).Value
			return evalFloatInfixExpression(operator, &ObjectFloat{Value: float64(leftVal)}, right)
		}
		if left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ {
			rightVal := right.(*ObjectInteger).Value
			return evalFloatInfixExpression(operator, left, &ObjectFloat{Value: float64(rightVal)})
		}
		return NewError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	default:
		return NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right Object,
) Object {
	var leftVal, rightVal float64
	switch left := left.(type) {
	case *ObjectFloat:
		leftVal = left.Value
	case *ObjectInteger:
		leftVal = float64(left.Value)
	default:
		return NewError("left operand is not a number: %s", left.Type())
	}
	switch right := right.(type) {
	case *ObjectFloat:
		rightVal = right.Value
	case *ObjectInteger:
		rightVal = float64(right.Value)
	default:
		return NewError("right operand is not a number: %s", right.Type())
	}

	switch operator {
	case "+":
		return &ObjectFloat{Value: leftVal + rightVal}
	case "-":
		return &ObjectFloat{Value: leftVal - rightVal}
	case "*":
		return &ObjectFloat{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return NewError("division by zero")
		}
		return &ObjectFloat{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right Object,
) Object {
	if operator != "+" {
		return NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
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
		return NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalProgram(program *Program, env *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalWhileStatement(stmt *WhileStatement, env *Environment) Object {
	for {
		condition := Eval(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := Eval(stmt.Body, env)
		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
	}
	return EVAL_NULL
}

func evalForStatement(stmt *ForStatement, env *Environment) Object {
	start := Eval(stmt.Start, env)
	if isError(start) {
		return start
	}

	end := Eval(stmt.End, env)
	if isError(end) {
		return end
	}

	startInt, ok := start.(*ObjectInteger)
	if !ok {
		return NewError("for loop start value must be integer, got %s", start.Type())
	}

	endInt, ok := end.(*ObjectInteger)
	if !ok {
		return NewError("for loop end value must be integer, got %s", end.Type())
	}

	forEnv := NewEnclosedEnvironment(env)

	for i := startInt.Value; i <= endInt.Value; i++ {
		forEnv.Set(stmt.Variable.Value, &ObjectInteger{Value: i})

		result := Eval(stmt.Body, forEnv)
		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
	}

	return EVAL_NULL
}

func NewError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func evalCompoundAssignmentStatement(stmt *CompoundAssignmentStatement, env *Environment) Object {
	// Get the current value of the variable
	current := evalIdentifier(stmt.Name, env)
	if isError(current) {
		return current
	}

	// Evaluate the right-hand side expression
	val := Eval(stmt.Value, env)
	if isError(val) {
		return val
	}

	var result Object

	// Perform the compound operation
	switch stmt.Operator {
	case "+=":
		result = evalInfixExpression("+", current, val)
	case "-=":
		result = evalInfixExpression("-", current, val)
	default:
		return NewError("unknown compound assignment operator: %s", stmt.Operator)
	}

	if isError(result) {
		return result
	}

	// Set the new value
	env.Set(stmt.Name.Value, result)
	return result
}
