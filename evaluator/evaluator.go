package evaluator

import (
	"github.com/kanersps/loop/ast"
	"github.com/kanersps/loop/evaluator/helpers"
	"github.com/kanersps/loop/models"
)

func Eval(node ast.Node, env *models.Environment) models.Object {
	return helpers.Eval(node, env)
}
