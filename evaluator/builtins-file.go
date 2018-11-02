package evaluator

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/kenshindeveloper/april/object"
)

var bfile = map[string]*object.Builtin{

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

			if stream.FILE == nil {
				return newError("variable file is equal to null.")
			}

			// w := bufio.NewWriter(stream.FILE)
			// w.WriteString(fmt.Sprintf(str.Value))
			// w.Flush()
			io.WriteString(stream.FILE, str.Value)

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

			return &object.String{Value: string(data)}
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
}
