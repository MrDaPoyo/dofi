package balena

import (
	env "github.com/mrdapoyo/dofi/balena/env"
	parser "github.com/mrdapoyo/dofi/balena/parser"
)

type BalenaFunction struct {
	declaration *parser.FunctionStmt
	closure     *env.Environment
}

func NewBalenaFunction(declaration *parser.FunctionStmt, closure *env.Environment) *BalenaFunction {
	return &BalenaFunction{declaration: declaration, closure: closure}
}

func (lf *BalenaFunction) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}) {
	environment := env.NewEnclosedEnvironment(lf.closure)
	for i, param := range lf.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(Return); ok {
				result = returnValue.Value
			} else {
				panic(r)
			}
		}
	}()

	interpreter.executeBlock(lf.declaration.Body, environment)
	return
}

func (lf *BalenaFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lf *BalenaFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
