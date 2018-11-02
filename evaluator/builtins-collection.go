package evaluator

import (
	"github.com/kenshindeveloper/april/object"
)

var bcollection = map[string]*object.Builtin{

	//***************************************************************************************
	//*************************************** LIST ******************************************
	//***************************************************************************************

	"front": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.LIST_OBJ {
				switch length := len(args); length {
				case 0:
					return newError("wrong number of arguments. got'%d', want='1'", len(args))
				default:
					return newError("argument to 'front' must be LIST, got='%s'", args[0].Type())
				}
			}

			l := args[0].(*object.List)
			if len(l.Elements) > 0 {
				return l.Elements[0]
			}
			return NIL
		},
	},

	"back": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.LIST_OBJ {
				switch length := len(args); length {
				case 0:
					return newError("wrong number of arguments. got'%d', want='1'", len(args))
				default:
					return newError("argument to 'back' must be LIST, got='%s'", args[0].Type())
				}
			}

			l := args[0].(*object.List)
			length := len(l.Elements)
			if length > 0 {
				return l.Elements[length-1]
			}
			return NIL
		},
	},

	"pop": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 || len(args) >= 3 || args[0].Type() != object.LIST_OBJ {
				switch length := len(args); length {
				case 0:
					return newError("wrong number of arguments. got'%d', want='1 or 2'", len(args))
				default:
					return newError("argument to 'pop' must be LIST, got='%s'", args[0].Type())
				}
			}

			l := args[0].(*object.List)
			length := len(l.Elements)
			if length == 0 {
				return NIL
			}

			switch len(args) {
			case 1:
				obj := l.Elements[length-1]
				l.Elements = append(l.Elements[:len(l.Elements)-1], l.Elements[len(l.Elements):]...)
				return obj
			default:
				objInteger, ok := args[1].(*object.Integer)
				if !ok {
					return newError("expression is not interger.")
				}

				if objInteger.Value >= int64(length) || objInteger.Value < 0 {
					return newError("index out range.")
				}

				obj := l.Elements[objInteger.Value]
				l.Elements = append(l.Elements[:objInteger.Value], l.Elements[objInteger.Value+1:]...)
				return obj
			}
		},
	},

	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.LIST_OBJ {
				if len(args) != 2 {
					return newError("wrong number of arguments. got'%d', want='2'", len(args))
				}
				return newError("argument to 'push' must be LIST, got='%s'", args[0].Type())
			}
			l := args[0].(*object.List)
			l.Elements = append(l.Elements, args[1])
			return l
		},
	},

	"index": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.LIST_OBJ {
				if len(args) != 2 {
					return newError("wrong number of arguments. got'%d', want='2'", len(args))
				}
				return newError("argument to 'push' must be LIST, got='%s'", args[0].Type())
			}

			_, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", args[1].Type())
			}

			l := args[0].(*object.List)
			index := 0
			for _, element := range l.Elements {
				switch args[1].Type() {
				case object.INTEGER_OBJ:
					if element.Type() == object.INTEGER_OBJ && element.(*object.Integer).Value == args[1].(*object.Integer).Value {
						return &object.Integer{Value: int64(index)}
					}
				case object.DOUBLE_OBJ:
					if element.Type() == object.DOUBLE_OBJ && element.(*object.Double).Value == args[1].(*object.Double).Value {
						return &object.Integer{Value: int64(index)}
					}
				case object.STRING_OBJ:
					if element.Type() == object.STRING_OBJ && element.(*object.String).Value == args[1].(*object.String).Value {
						return &object.Integer{Value: int64(index)}
					}
				case object.BOOLEAN_OBJ:
					if element.Type() == object.BOOLEAN_OBJ && element.(*object.Boolean).Value == args[1].(*object.Boolean).Value {
						return &object.Integer{Value: int64(index)}
					}
				}

				index++
			}
			return NIL
		},
	},

	"range": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 || len(args) > 2 {
				return newError("wrong number of arguments. got'%d', want='1 or 2'", len(args))
			}

			intObj, ok := args[0].(*object.Integer)
			if !ok {
				return newError("variable must be a integer.")
			}

			switch len(args) {
			case 1:
				if intObj.Value < 0 {
					return newError("integer must be >= 0.")
				}
				list := &object.List{Elements: []object.Object{}}
				for i := 0; int64(i) < intObj.Value; i++ {
					list.Elements = append(list.Elements, &object.Integer{Value: int64(i)})
				}
				return list

			default:
				intObj1, ok := args[1].(*object.Integer)
				if !ok {
					return newError("variable must be a integer.")
				}

				if intObj.Value > intObj1.Value {
					return newError("left variable must be >= right variable.")
				}

				list := &object.List{Elements: []object.Object{}}
				for i := intObj.Value; int64(i) < intObj1.Value; i++ {
					list.Elements = append(list.Elements, &object.Integer{Value: int64(i)})
				}
				return list
			}
		},
	},

	//***************************************************************************************
	//*************************************** MAPS ******************************************
	//***************************************************************************************

	"delete": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='2'", len(args))
			}

			hash, ok := args[0].(*object.Hash)
			if !ok {
				return newError("is not a map.")
			}

			key, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", args[1].Type())
			}

			pairs, ok := hash.Pairs[key.HashKey()]
			if ok {
				delete(hash.Pairs, key.HashKey())
			} else {
				return newError("key error: %s", args[1].Type())
			}

			return pairs.Value
		},
	},

	"find": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='2'", len(args))
			}

			hash, ok := args[0].(*object.Hash)
			if !ok {
				return newError("is not a map.")
			}

			key, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", args[1].Type())
			}

			if _, ok := hash.Pairs[key.HashKey()]; ok {
				return TRUE
			}

			return FALSE
		},
	},
}
