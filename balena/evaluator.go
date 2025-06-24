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
		return nil
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
	case *CStyleForStatement:
		return evalCStyleForStatement(node, env)
	case *GoObjectLiteral:
		return &GoObject{Value: node.Value}
	case *HashLiteral:
		return evalHashLiteral(node, env)

	case *AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *CompoundAssignmentStatement:
		current, ok := env.Get(node.Name.Value)
		if !ok {
			return NewError("identifier not found: %s", node.Name.Value)
		}

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		var operator string
		switch node.Operator {
		case "+=":
			operator = "+"
		case "-=":
			operator = "-"
		default:
			return NewError("unknown compound operator: %s", node.Operator)
		}

		result := evalInfixExpression(operator, current, val)
		if isError(result) {
			return result
		}

		env.Set(node.Name.Value, result)
		return result
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
	case left.Type() == HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return NewError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(hash, index Object) Object {
	hashObject := hash.(*Hash)
	key, ok := index.(Hashable)
	if !ok {
		return NewError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return EVAL_NULL
	}
	return pair.Value
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
	switch right.Type() {
	case INTEGER_OBJ:
		value := right.(*ObjectInteger).Value
		return &ObjectInteger{Value: -value}
	case FLOAT_OBJ:
		value := right.(*ObjectFloat).Value
		return &ObjectFloat{Value: -value}
	default:
		return NewError("unknown operator: -%s", right.Type())
	}
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
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case (left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ) ||
		(left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ):
		return evalMixedNumericInfixExpression(operator, left, right)
	case left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return NewError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return NewError("eval infix unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// Add this function after evalInfixExpression
func evalBooleanInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*ObjectBoolean).Value
	rightVal := right.(*ObjectBoolean).Value
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	default:
		return NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalMixedNumericInfixExpression(operator string, left, right Object) Object {
	var leftVal, rightVal float64

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*ObjectInteger).Value)
	} else {
		leftVal = left.(*ObjectFloat).Value
	}

	if right.Type() == INTEGER_OBJ {
		rightVal = float64(right.(*ObjectInteger).Value)
	} else {
		rightVal = right.(*ObjectFloat).Value
	}

	switch operator {
	case "+":
		return &ObjectFloat{Value: leftVal + rightVal}
	case "-":
		return &ObjectFloat{Value: leftVal - rightVal}
	case "*":
		return &ObjectFloat{Value: leftVal * rightVal}
	case "/":
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
		return NewError("unknown operator: %s", operator)
	}
}

func evalFloatInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*ObjectFloat).Value
	rightVal := right.(*ObjectFloat).Value

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
		return NewError("unknown operator: %s", operator)
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
		if rightVal == 0 {
			return NewError("division by zero")
		}
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

func evalCStyleForStatement(stmt *CStyleForStatement, env *Environment) Object {
	forEnv := NewEnclosedEnvironment(env)
	
	if stmt.Init != nil {
		initResult := Eval(stmt.Init, forEnv)
		if isError(initResult) {
			return initResult
		}
	}
	
	for {
		if stmt.Condition != nil {
			condition := Eval(stmt.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}
		
		result := Eval(stmt.Body, forEnv)
		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
		
		if stmt.Increment != nil {
			incResult := Eval(stmt.Increment, forEnv)
			if isError(incResult) {
				return incResult
			}
		}
	}
	
	return EVAL_NULL
}

func evalForStatement(stmt *ForStatement, env *Environment) Object {
	forEnv := NewEnclosedEnvironment(env)
	for param := range env.modifiedVars {
		if _, ok := forEnv.Get(param); ok {
			return NewError("variable already defined: %s", param)
		}
		forEnv.Set(param, EVAL_NULL)
	}

	startObj := Eval(stmt.Start, forEnv)
	if isError(startObj) {
		return startObj
	}
	endObj := Eval(stmt.End, forEnv)
	if isError(endObj) {
		return endObj
	}

	startInt, startIsInt := startObj.(*ObjectInteger)
	endInt, endIsInt := endObj.(*ObjectInteger)

	if startIsInt && endIsInt {
		step := int64(1)
		if startInt.Value > endInt.Value {
			step = -1
		}
		for i := startInt.Value; (step > 0 && i <= endInt.Value) || (step < 0 && i >= endInt.Value); i += step {
			forEnv.Set(stmt.Variable.Value, &ObjectInteger{Value: i})
			result := Eval(stmt.Body, forEnv)
			if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
				return result
			}
		}
		return EVAL_NULL
	}

	// fallback to float math if either bound is float
	var startF, endF float64
	if startIsInt {
		startF = float64(startInt.Value)
	} else if s, ok := startObj.(*ObjectFloat); ok {
		startF = s.Value
	} else {
		return NewError("for loop start must be int or float")
	}
	if endIsInt {
		endF = float64(endInt.Value)
	} else if e, ok := endObj.(*ObjectFloat); ok {
		endF = e.Value
	} else {
		return NewError("for loop end must be int or float")
	}
	step := 1.0
	if startF > endF {
		step = -1.0
	}
	epsilon := 1e-9
	for i := startF; (step > 0 && i <= endF+epsilon) || (step < 0 && i >= endF-epsilon); i += step {
		forEnv.Set(stmt.Variable.Value, &ObjectFloat{Value: i})
		result := Eval(stmt.Body, forEnv)
		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
	}
	return EVAL_NULL
}

func evalHashLiteral(
	node *HashLiteral,
	env *Environment,
) Object {
	pairs := make(map[HashKey]HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(Hashable)
		if !ok {
			return NewError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}
	return &Hash{Pairs: pairs}
}

func NewError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}