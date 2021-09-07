package builtins

import (
	"fmt"
	"github.com/kanersps/loop/models"
	"log"
	"net/http"
)

type HttpEndpoint struct {
	Endpoint string
	Handler  *models.Function
}

type HttpHandler struct {
	Env      *models.Environment
	Endpoint string
	Handler  *models.Function
}

func (h *HttpHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	returns := ApplyFunction(h.Handler, []models.Object{}, h.Env)

	body := ""

	str, ok := returns.(*models.String)

	if !ok {
		body = returns.(*models.Integer).Inspect()
	} else {
		body = str.Inspect()
	}

	fmt.Fprint(w, body)
}

type applyFunction func(fn models.Object, args []models.Object, env *models.Environment) models.Object

var ApplyFunction applyFunction

func SetApplyFunction(a applyFunction) {
	ApplyFunction = a
}

var Functions = map[string]*models.Builtin{
	"len": {
		Func: func(env *models.Environment, args ...models.Object) models.Object {
			if len(args) <= 0 || len(args) >= 2 {
				return &models.Error{Message: fmt.Sprintf("WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `len`. expected=1. got=%d", len(args))}
			}

			arg, ok := args[0].(*models.String)

			if !ok {
				return &models.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `len`. got=%v. expected=STRING", args[0].Type())}
			}

			return &models.Integer{Value: int64(len(arg.Value))}
		},
	},
	"append": {
		Func: func(env *models.Environment, args ...models.Object) models.Object {
			if len(args) <= 0 || len(args) == 1 {
				return &models.Error{Message: fmt.Sprintf("WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `append`. expected=2. got=%d", len(args))}
			}

			array, ok := args[0].(*models.Array)

			if !ok {
				return &models.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `append` (argument 0). expected=ARRAY. got=%v", args[0])}
			}

			return &models.Array{Elements: append(array.Elements, args[1:]...)}
		},
	},
	"print": {
		Func: func(env *models.Environment, args ...models.Object) models.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}

			return models.NULL
		},
	},
	"println": {
		Func: func(env *models.Environment, args ...models.Object) models.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return models.NULL
		},
	},
	"webserver": {
		Func: func(env *models.Environment, args ...models.Object) models.Object {
			if len(args) == 0 || len(args) >= 3 {
				return &models.Error{Message: fmt.Sprintf("WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `len`. expected=2. got=%d", len(args))}
			}

			port, ok := args[0].(*models.Integer)
			addr := fmt.Sprintf(":%s", port.Inspect())

			if !ok {
				return &models.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `webserver` (argument 0). expected=INTEGER. got=%v", args[0].Type())}
			}

			config, ok := args[1].(*models.Hash)

			if !ok {
				return &models.Error{Message: fmt.Sprintf("ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `webserver` (argument 1). expected=HASH. got=%v", args[1].Type())}
			}

			for _, v := range config.Pairs {
				fmt.Println(v.Key.Inspect())
				fn := v.Value.(*models.Function)

				handler := &HttpHandler{
					Env:      env,
					Endpoint: v.Key.Inspect(),
					Handler:  fn,
				}

				http.HandleFunc(v.Key.Inspect(), handler.handleRequest)
			}

			log.Fatal(http.ListenAndServe(addr, nil))

			return models.NULL
		},
	},
}
