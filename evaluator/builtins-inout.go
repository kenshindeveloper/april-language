package evaluator

import (
	"fmt"

	"github.com/kenshindeveloper/april/object"
)

var binouts = map[string]*object.Builtin{

	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NIL
		},
	},

	"printf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			str := args[0].Inspect()[1 : len(args[0].Inspect())-1]
			fmt.Printf(str)
			return NIL
		},
	},
}
