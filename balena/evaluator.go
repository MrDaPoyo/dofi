package main

func Eval(node Node) Object {
	switch node := node.(type) {
	// Statements
	case *Program:
		return evalStatements(node.Statements)
	case *ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *IntegerLiteral:
		return &ObjectInteger{Value: node.Value}
	}
	return nil
}

func evalStatements(stmts []Statement) Object {
	var result Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}
