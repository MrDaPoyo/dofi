package balena

type BalenaCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	String() string
}