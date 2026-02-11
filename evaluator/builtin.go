package evaluator

import (
	"fmt"
	"os"

	"monkey/object"
)

var builtins = map[string]object.BuiltinFnObject{
	"puts": {
		Name: "puts",
		Function: func(args ...object.Object) object.Object {
			rawArgs := []any{}
			for _, a := range args {
				if !isOneOfTypes(
					a,
					object.STRING,
					object.INT,
					object.BOOL,
					object.ARRAY,
					object.HASH,
				) {
					return &object.ErrorObject{
						Message: &object.StringObject{
							Value: fmt.Sprintf("cannot call 'puts' on %s", a.Type()),
						},
					}
				}
				rawArgs = append(rawArgs, a.Inspect())
			}
			fmt.Println(rawArgs...)
			return object.NULL_OBJECT
		},
	},
	"readFile": {
		Name: "readFile",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 || !isType(object.STRING, args[0]) {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'readFile' accepts only single STRING argument, but was %s: ",
							types(args),
						),
					},
				}
			}
			path := args[0].(*object.StringObject)

			file, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'readFile': %s:",
							err.Error(),
						),
					},
				}
			}
			return &object.StringObject{Value: string(file)}
		},
	},
	"writeFile": {
		Name: "writeFile",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 || !isType(object.STRING, args[0]) ||
				!isType(object.STRING, args[1]) {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'writeFile' accepts only path as STRING, and contents as STRING arguments, but was %s: ",
							types(args),
						),
					},
				}
			}
			path := args[0].(*object.StringObject)
			content := args[1].(*object.StringObject)

			err := os.WriteFile(path.Value, []byte(content.Value), 0x777)
			if err != nil {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf("'writeFile': %s:", err.Error()),
					},
				}
			}
			return object.TRUE_OBJECT
		},
	},
	"len": {
		Name: "len",
		Function: func(args ...object.Object) object.Object {
			if len(args) == 0 || len(args) > 1 {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'len' requires at least one argument, but had %d",
							len(args),
						),
					},
				}
			}

			switch arg := args[0].(type) {
			case *object.StringObject:
				return &object.IntObject{Value: int64(len(arg.Value))}
			case *object.ArrayObject:
				return &object.IntObject{Value: int64(len(arg.Items))}
			default:
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'len' accepts only STRING or ARRAY arguments, but was %s: ",
							args[0].Type(),
						),
					},
				}
			}
		},
	},
	"first": {
		Name: "first",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'first' requires exactly one argument, but had %d",
							len(args),
						),
					},
				}
			}
			arr, isArray := args[0].(*object.ArrayObject)
			if !isArray {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'first' accepts only ARRAY argument, but was %s",
							args[0].Type(),
						),
					},
				}
			}
			if len(arr.Items) == 0 {
				return object.NULL_OBJECT
			}
			return arr.Items[0]
		},
	},
	"last": {
		Name: "last",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'last' requires exactly one argument, but had %d",
							len(args),
						),
					},
				}
			}
			arr, isArray := args[0].(*object.ArrayObject)
			if !isArray {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'last' accepts only ARRAY argument, but was %s",
							args[0].Type(),
						),
					},
				}
			}
			if len(arr.Items) == 0 {
				return object.NULL_OBJECT
			}
			return arr.Items[len(arr.Items)-1]
		},
	},
	"rest": {
		Name: "rest",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'rest' requires exactly one argument, but had %d",
							len(args),
						),
					},
				}
			}
			arr, isArray := args[0].(*object.ArrayObject)
			if !isArray {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'rest' accepts only ARRAY argument, but was %s",
							args[0].Type(),
						),
					},
				}
			}
			if len(arr.Items) == 0 {
				return &object.ArrayObject{Items: []object.Object{}}
			}
			newItems := make([]object.Object, len(arr.Items)-1)
			copy(newItems, arr.Items[1:])
			return &object.ArrayObject{Items: newItems}
		},
	},
	"push": {
		Name: "push",
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'push' requires exactly two arguments, but had %d",
							len(args),
						),
					},
				}
			}
			arr, isArray := args[0].(*object.ArrayObject)
			if !isArray {
				return &object.ErrorObject{
					Message: &object.StringObject{
						Value: fmt.Sprintf(
							"'push' first argument must be ARRAY, but was %s",
							args[0].Type(),
						),
					},
				}
			}
			newItems := make([]object.Object, len(arr.Items), len(arr.Items)+1)
			copy(newItems, arr.Items)
			newItems = append(newItems, args[1])
			return &object.ArrayObject{Items: newItems}
		},
	},
}
