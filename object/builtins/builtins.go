package builtins

import (
	"fmt"
	"github.com/kanersps/loop/object"
)

var Functions = map[string]*object.Builtin{
	"len": {
		Func: func(args ...object.Object) object.Object {
			if len(args) <= 0 || len(args) >= 2 {
				return &object.Error{Message: fmt.Sprintf("WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `len`. expected=1. got=%d", len(args))}
			}

			arg, ok := args[0].(*object.String)

			if !ok {
				return &object.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `len`. got=%v. epxected=STRING", args[0].Type())}
			}

			return &object.Integer{Value: int64(len(arg.Value))}
		},
	},
	"append": {
		Func: func(args ...object.Object) object.Object {
			if len(args) <= 0 || len(args) == 1 {
				return &object.Error{Message: fmt.Sprintf("WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `append`. expected=2. got=%d", len(args))}
			}

			array, ok := args[0].(*object.Array)

			if !ok {
				return &object.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `append` (argument 0). expected=ARRAY. got=%v", args[0])}
			}

			return &object.Array{Elements: append(array.Elements, args[1:]...)}
		},
	},
}
