package balena

import (
	"math"
	"sync"
	"time"
)

type GoBuiltin struct {
	fn    func(args ...interface{}) interface{}
	arity int
}

func (b *GoBuiltin) Arity() int { return b.arity }

func (b *GoBuiltin) Call(_ *Interpreter, arguments []interface{}) interface{} {
	return b.fn(arguments...)
}

func (b *GoBuiltin) String() string { return "<builtin fn>" }

var (
	builtinMu       sync.RWMutex
	builtinRegistry = map[string]BalenaCallable{}
)

func RegisterBuiltin(name string, arity int, fn func(args ...interface{}) interface{}) {
	builtinMu.Lock()
	defer builtinMu.Unlock()
	builtinRegistry[name] = &GoBuiltin{fn: fn, arity: arity}
}

func loadBuiltins(i *Interpreter) {
	builtinMu.RLock()
	defer builtinMu.RUnlock()
	for name, callable := range builtinRegistry {
		i.Globals.Define(name, callable)
	}
}

func init() {
	RegisterBuiltin("sumSlice", -1, func(args ...interface{}) interface{} {
		slice := args[0].([]float64)
		total := 0.0
		for _, v := range slice {
			total += v
		}
		return total
	})
	RegisterBuiltin("print", 1, func(args ...interface{}) interface{} {
		for _, arg := range args {
			if str, ok := arg.(string); ok {
				print(str)
			} else {
				print(arg)
			}
		}
		return nil
	})
	RegisterBuiltin("clock", 0, func(args ...interface{}) interface{} {
		return float64(time.Now().UnixNano()) / 1e9
	})
	RegisterBuiltin("len", 1, func(args ...interface{}) interface{} {
		if str, ok := args[0].(string); ok {
			return float64(len(str))
		} else if slice, ok := args[0].([]interface{}); ok {
			return float64(len(slice))
		}
		panic(&RuntimeError{
			Message: "len() expects a string or a slice.",
		})
	})
	RegisterBuiltin("sin", 1, func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic(&RuntimeError{
				Message: "sin() expects exactly one argument.",
			})
		}
		if num, ok := args[0].(float64); ok {
			return math.Sin(num)
		}
		panic(&RuntimeError{
			Message: "sin() expects a number.",
		})
	})
	RegisterBuiltin("cos", 1, func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic(&RuntimeError{
				Message: "cos() expects exactly one argument.",
			})
		}
		if num, ok := args[0].(float64); ok {
			return math.Cos(num)
		}
		panic(&RuntimeError{
			Message: "cos() expects a number.",
		})
	})
	RegisterBuiltin("tan", 1, func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic(&RuntimeError{
				Message: "tan() expects exactly one argument.",
			})
		}
		if num, ok := args[0].(float64); ok {
			return math.Tan(num)
		}
		panic(&RuntimeError{
			Message: "tan() expects a number.",
		})
	})
	RegisterBuiltin("tan", 1, func(args ...interface{}) interface{} {
		if len(args) != 1 {
			panic(&RuntimeError{
				Message: "tan() expects exactly one argument.",
			})
		}
		if num, ok := args[0].(float64); ok {
			return math.Tan(num)
		}
		panic(&RuntimeError{
			Message: "tan() expects a number.",
		})
	})
	RegisterBuiltin("array", -1, func(args ...interface{}) interface{} {
		if len(args) == 0 {
			return []interface{}{}
		}
		return args
	})
}
