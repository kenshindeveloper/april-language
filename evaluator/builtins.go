package evaluator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kenshindeveloper/april/object"
)

func getFormat(str string) string {
	final := str
	for _, char := range []string{"n", "r", "t", "a"} {
		if strings.Contains(final, "\\"+char) {
			l := strings.Split(final, "\\"+char)
			switch char {
			case "a":
				final = strings.Join(l, "\a") // Attempts to alert the user through audible or visible notification.
			case "n":
				final = strings.Join(l, "\n") // Moves the printing position to the start of the next line.
			case "r":
				final = strings.Join(l, "\r") // Moves the printing position to the start of the current line.
			case "t":
				final = strings.Join(l, "\t") // Moves the printing position to the next tab position on the current line. If there are no more tab positions left on the line, the behaviour is undefined.

			}
		}
	}

	return final
}

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
	//***************************************************************************************
	//***************************************************************************************

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
	//***************************************************************************************
	//***************************************************************************************
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

	//***************************************************************************************
	//***************************************************************************************
	//***************************************************************************************
	"type": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}
			return &object.String{Value: object.GetType(args[0].Type())}
		},
	},

	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(getFormat(arg.Inspect()))
			}
			return NIL
		},
	},
	"printf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			str := args[0].Inspect()
			fmt.Printf(getFormat(str) + "\n")
			return NIL
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
	//***************************************************************************************
	//***************************************************************************************
	//***************************************************************************************
	"open": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			file, err := os.Open(path.Value)
			if err != nil {
				return newError("error to open file: '%s'", path.Value) //manejo de la variable NIL
			}

			return &object.Stream{FILE: file}
		},
	},

	"isExist": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type stream.")
			}

			_, err := os.Stat(path.Value)
			if err != nil {
				if os.IsNotExist(err) {
					return FALSE
				}
			}

			return TRUE
		},
	},

	"isOpen": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			stream, ok := args[0].(*object.Stream)
			if !ok {
				return newError("argument is not type stream.")
			}

			if stream.FILE == nil {
				return FALSE
			}

			return TRUE
		},
	},

	"create": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			file, err := os.Create(path.Value)
			if err != nil {
				return newError("error to create file: '%s'", path.Value)
			}

			return &object.Stream{FILE: file}
		},
	},
	"rename": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			originalName, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			newName, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			err := os.Rename(originalName.Value, newName.Value)
			if err != nil {
				return newError("file could not be renamed")
			}

			return TRUE
		},
	},

	"move": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			originalPath, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			newPath, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			err := os.Rename(originalPath.Value, newPath.Value)
			if err != nil {
				return newError("file could not be renamed")
			}

			return TRUE
		},
	},

	"write": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='2'", len(args))
			}

			stream, ok := args[0].(*object.Stream)
			if !ok {
				return newError("firs argument is not type stream.")
			}

			str, ok := args[1].(*object.String)
			if !ok {
				return newError("second argument is not type string.")
			}
			str.Value = getFormat(str.Value)

			if stream.FILE == nil {
				return newError("variable file is equal to null.")
			}

			w := bufio.NewWriter(stream.FILE)
			w.WriteString(str.Value)
			w.Flush()

			return NIL
		},
	},

	"read": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			stream, ok := args[0].(*object.Stream)
			if !ok {
				return newError("firs argument is not type stream.")
			}

			if stream.FILE == nil {
				return newError("variable file is equal to null.")
			}

			data, err := ioutil.ReadFile(stream.FILE.Name())
			if err != nil {
				return newError("error to read file.")
			}

			return &object.String{Value: getFormat(string(data))}
		},
	},

	"close": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			stream, ok := args[0].(*object.Stream)
			if !ok {
				return newError("argument is not type stream.")
			}

			if stream.FILE == nil {
				return newError("variable file is equal to null.")
			}

			stream.FILE.Close()
			stream.FILE = nil
			return NIL
		},
	},

	"remove": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument is not type string.")
			}

			err := os.Remove(str.Value)
			if err != nil {
				return newError("file could not be deleted.")
			}

			return TRUE
		},
	},
	//***************************************************************************************
	//***************************************************************************************
	//***************************************************************************************
	"GoCeil": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Ceil(x.Value)}
		},
	},

	"GoFoor": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Floor(x.Value)}
		},
	},

	"GoLog": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Log(x.Value)}
		},
	},

	"GoLog10": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Log10(x.Value)}
		},
	},

	"GoPow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			y, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Pow(x.Value, y.Value)}
		},
	},

	"GoPow10": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Pow10(int(x.Value))}
		},
	},

	"GoRound": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Round(x.Value)}
		},
	},

	"GoTrunc": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x, ok := args[0].(*object.Double)
			if !ok {
				return newError("argument is not type Double.")
			}

			return &object.Double{Value: math.Trunc(x.Value)}
		},
	},

	"GoMax": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x1, ok1 := args[0].(*object.Integer)
			y1, ok2 := args[1].(*object.Integer)
			if ok1 {
				if !ok2 {
					return newError("argument is not type Integer.")
				} else {
					return &object.Double{Value: math.Max(float64(x1.Value), float64(y1.Value))}
				}
			}

			x2, ok1 := args[0].(*object.Double)
			y2, ok2 := args[1].(*object.Double)
			if ok1 {
				if !ok2 {
					return newError("argument is not type Integer o Double.")
				} else {
					return &object.Double{Value: math.Max(x2.Value, y2.Value)}
				}
			} else {
				return newError("argument is not type Integer o Double.")
			}
		},
	},

	"GoMin": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			x1, ok1 := args[0].(*object.Integer)
			y1, ok2 := args[1].(*object.Integer)
			if ok1 {
				if !ok2 {
					return newError("argument is not type Integer.")
				} else {
					return &object.Double{Value: math.Max(float64(x1.Value), float64(y1.Value))}
				}
			}

			x2, ok1 := args[0].(*object.Double)
			y2, ok2 := args[1].(*object.Double)
			if ok1 {
				if !ok2 {
					return newError("argument is not type Integer o Double.")
				} else {
					return &object.Double{Value: math.Min(x2.Value, y2.Value)}
				}
			} else {
				return newError("argument is not type Integer o Double.")
			}
		},
	},
	//***************************************************************************************
	//***************************************************************************************
	//***************************************************************************************
	"GoNow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			now := time.Now()
			println(now.String())
			return NIL
		},
	},

	"GoYear": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			now := time.Now()
			Year := now.Year()
			fmt.Printf("%v", Year)
			return NIL
		},
	},

	"GoMonth": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			now := time.Now()
			month := now.Month()
			fmt.Printf("%v", month)
			return NIL
		},
	},

	"GoDay": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got'%d', want='1'", len(args))
			}

			now := time.Now()
			day := now.Day()
			fmt.Printf("%v", day)
			return NIL
		},
	},
}
