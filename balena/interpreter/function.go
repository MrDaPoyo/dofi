package balena

import (
	parser "github.com/mrdapoyo/dofi/balena/parser"
	env "github.com/mrdapoyo/dofi/balena/env"
)

type BalenaFunction struct {
	declaration *parser.FunctionStmt
	closure     *env.Environment
}

func NewBalenaFunction(declaration *parser.FunctionStmt, closure *env.Environment) *BalenaFunction {
	return &BalenaFunction{declaration: declaration, closure: closure}
}

func (lf *BalenaFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := env.NewEnvironment(interpreter.Globals)
	for i, param := range lf.declaration.Params {
		environment.Set(param.Lexeme, arguments[i])
	}
	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(Return); ok {
				panic(returnValue)
			} else {
				panic(r)
			}
		}
	}()
	interpreter.executeBlock(lf.declaration.Body, environment)
	return nil
}

func (lf *BalenaFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lf *BalenaFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}

