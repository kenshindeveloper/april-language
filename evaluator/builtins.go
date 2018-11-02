package evaluator

import (
	"os"
	"strconv"

	"github.com/kenshindeveloper/april/object"
)

var builtins = map[string]*object.Builtin{

	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.List:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported, got='%s'", args[0].Type())
			}
		},
	},

	"str": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			switch args[0].(type) {
			case *object.String:
				return &object.String{Value: args[0].(*object.String).Value}
			case *object.Integer:
				return &object.String{Value: strconv.FormatInt(args[0].(*object.Integer).Value, 10)}
			case *object.Double:
				return &object.String{Value: strconv.FormatFloat(args[0].(*object.Double).Value, 'g', 10, 64)}
			case *object.Boolean:
				return &object.String{Value: strconv.FormatBool(args[0].(*object.Boolean).Value)}
			default:
				return newError("function 'str' not supported to '%s'", args[0].Type())
			}
		},
	},

	"int": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			switch args[0].(type) {
			case *object.Integer:
				return &object.Integer{Value: args[0].(*object.Integer).Value}
			case *object.Double:
				return &object.Integer{Value: int64(args[0].(*object.Double).Value)}
			case *object.String:
				value, err := strconv.ParseInt(args[0].(*object.String).Value, 10, 64)
				if err != nil {
					return newError("error to convert '%s' to integer.", args[0].(*object.String).Value)
				}

				return &object.Integer{Value: value}

			default:
				return newError("function 'int' not supported to '%s'", args[0].Type())
			}
		},
	},

	"double": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			switch args[0].(type) {
			case *object.Double:
				return &object.Double{Value: args[0].(*object.Double).Value}

			case *object.Integer:
				return &object.Double{Value: float64(args[0].(*object.Integer).Value)}

			case *object.String:
				value, err := strconv.ParseFloat(args[0].(*object.String).Value, 64)
				if err != nil {
					return newError("error to convert '%s' to double.", args[0].(*object.String).Value)
				}

				return &object.Double{Value: value}

			default:
				return newError("function 'double' not supported to '%s'", args[0].Type())
			}
		},
	},

	"type": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}
			return &object.String{Value: object.GetType(args[0].Type())}
		},
	},

	"exit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got'%d', want='0'", len(args))
			}
			os.Exit(0)

			return NIL
		},
	},
}

func joinBuiltins(mapa map[string]*object.Builtin) {

	for key, value := range mapa {
		builtins[key] = value
	}
}

func InitBuiltins() {

	joinBuiltins(bcollection)
	joinBuiltins(bfile)
	joinBuiltins(binouts)
	joinBuiltins(bmath)
	joinBuiltins(bstring)
	joinBuiltins(btime)
}
