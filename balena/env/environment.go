package balena

import (
	"fmt"
	token "github.com/mrdapoyo/dofi/balena/token"
)
// import "fmt"
//
// func NewEnclosedEnvironment(outer *Environment) *Environment {
// 	env := NewEnvironment()
// 	env.outer = outer
// 	return env
// }
//
// func RegisterExternalBinding(name string, fn func(args ...Object) Object) {
// 	if _, exists := builtins[name]; exists {
// 		panic("builtin already exists: " + name)
// 	}
// 	builtins[name] = &Builtin{Fn: fn}
// }
//
// func NewEnvironment() *Environment {
// 	s := make(map[string]Object)
// 	return &Environment{store: s, outer: nil}
// }
//
// type Environment struct {
// 	store map[string]Object
// 	outer *Environment
// 	modifiedVars map[string]struct{}
// }
//
// func (e *Environment) Get(name string) (Object, bool) {
// 	obj, ok := e.store[name]
// 	if !ok && e.outer != nil {
// 		obj, ok = e.outer.Get(name)
// 	}
// 	return obj, ok
// }
//
// func (e *Environment) Set(name string, val Object) Object {
// 	e.store[name] = val
// 	return val
// }
//
// func (e *Environment) SetUserData(key string, value Object) {
// 	if e.store == nil {
// 		e.store = make(map[string]Object)
// 	}
// 	e.store[key] = value
// }
//
// func (e *Environment) NewGoObject(val interface{}) Object {
// 	return &GoObject{Value: val}
// }
//
// func (e *Environment) CallFunction(name string, args ...Object) (Object, error) {
// 	obj, ok := e.Get(name)
// 	if ok {
// 		switch fn := obj.(type) {
// 		case *Function:
// 			result := CallFunction(e, name, args...)
// 			if errObj, ok := result.(*Error); ok {
// 				return nil, fmt.Errorf("%s", errObj.Message)
// 			}
// 			return result, nil
// 		case *Builtin:
// 			return fn.Fn(args...), nil
// 		default:
// 			return nil, fmt.Errorf("%s is not a function", name)
// 		}
// 	}
//
// 	builtin, exists := builtins[name]
// 	if !exists {
// 		return nil, fmt.Errorf("function not found: %s", name)
// 	}
// 	return builtin.Fn(args...), nil
// }
//
// func (e *Environment) GetGoObject(name string) (interface{}, error) {
// 	obj, ok := e.Get(name)
// 	if !ok {
// 		return nil, fmt.Errorf("object not found: %s", name)
// 	}
// 	goObj, ok := obj.(*GoObject)
// 	if !ok {
// 		return nil, fmt.Errorf("%s is not a Go object", name)
// 	}
// 	return goObj.Value, nil
// }

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnclosedEnvironment(Enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: Enclosing,
		Values:    make(map[string]interface{}),
	}
}

func NewEnvironment(previous ...*Environment) *Environment {
	env := &Environment{
		Values: make(map[string]interface{}),
	}
	if len(previous) > 0 {
		env.Enclosing = previous[0]
	}
	return env
}

func (env *Environment) Set(key string, value interface{}) {
	if env.Enclosing != nil {
		env.Enclosing.Set(key, value)
		return
	}
	env.Values[key] = value
}

func (env *Environment) Get(key string) (interface{}, bool) {
	if value, exists := env.Values[key]; exists {
		return value, true
	}
	if env.Enclosing != nil {
		return env.Enclosing.Get(key)
	}
	value, exists := env.Values[key]
	if !exists && env.Values != nil {
		panic("key not found: " + key)
	}

	return value, exists
}
func (env *Environment) Assign(name token.Token, value interface{}) error {
	if _, exists := env.Values[name.Lexeme]; exists {
		env.Values[name.Lexeme] = value
		return nil
	}
	return fmt.Errorf("undefined variable '%s'.", name.Lexeme)
}
