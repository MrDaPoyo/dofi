package balena

import "math"

var builtins = map[string]*Builtin{
	"len": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				return &ObjectInteger{Value: int64(len(arg.Elements))}
			case *String:
				return &ObjectInteger{Value: int64(len(arg.Value))}
			default:
				return NewError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return NewError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return EVAL_NULL
		},
	},
	"last": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return NewError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return EVAL_NULL
		},
	},
	"rest": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return NewError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &Array{Elements: newElements}
			}
			return EVAL_NULL
		},
	},
	"cos": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != INTEGER_OBJ && args[0].Type() != FLOAT_OBJ {
				return NewError("argument to `cos` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
			if args[0].Type() == INTEGER_OBJ {
				intArg := args[0].(*ObjectInteger)
				return &ObjectFloat{Value: math.Cos(float64(intArg.Value))}
			}
			floatArg := args[0].(*ObjectFloat)
			return &ObjectFloat{Value: math.Cos(floatArg.Value)}
		},
	},
	"sin": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != INTEGER_OBJ && args[0].Type() != FLOAT_OBJ {
				return NewError("argument to `sin` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
			if args[0].Type() == INTEGER_OBJ {
				intArg := args[0].(*ObjectInteger)
				return &ObjectFloat{Value: math.Sin(float64(intArg.Value))}
			}
			floatArg := args[0].(*ObjectFloat)
			return &ObjectFloat{Value: math.Sin(floatArg.Value)}
		},
	},
	"floor": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *ObjectInteger:
				return &ObjectInteger{Value: int64(math.Floor(float64(arg.Value)))}
			case *ObjectFloat:
				return &ObjectInteger{Value: int64(math.Floor(arg.Value))}
			default:
				return NewError("argument to `floor` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
		},
	},
	"array": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 2 {
				if size, ok := args[0].(*ObjectInteger); ok {
					elements := make([]Object, size.Value)
					for i := range elements {
						elements[i] = args[1]
					}
					return &Array{Elements: elements}
				}
			}
			return &Array{Elements: args}
		},
	},
	"int": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() == NULL_OBJ {
				return &ObjectInteger{Value: 0}
			}
			switch arg := args[0].(type) {
			case *ObjectInteger:
				return arg
			case *ObjectFloat:
				return &ObjectInteger{Value: int64(arg.Value)}
			default:
				return NewError("argument to `int` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
		},
	},
	"float": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() == NULL_OBJ {
				return &ObjectFloat{Value: 0.0}
			}
			switch arg := args[0].(type) {
			case *ObjectInteger:
				return &ObjectFloat{Value: float64(arg.Value)}
			case *ObjectFloat:
				return arg
			default:
				return NewError("argument to `float` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
		},
	},
	"abs": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *ObjectInteger:
				v := arg.Value
				if v < 0 {
					v = -v
				}
				return &ObjectInteger{Value: v}
			case *ObjectFloat:
				v := arg.Value
				if v < 0 {
					v = -v
				}
				return &ObjectFloat{Value: v}
			default:
				return NewError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
}

var mathObject = map[string]*Builtin{
	"floor": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *ObjectInteger:
				return &ObjectInteger{Value: int64(math.Floor(float64(arg.Value)))}
			case *ObjectFloat:
				return &ObjectInteger{Value: int64(math.Floor(arg.Value))}
			default:
				return NewError("argument to `floor` must be INTEGER or FLOAT, got %s",
					args[0].Type())
			}
		},
	},
}
