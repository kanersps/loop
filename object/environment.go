package object

import "github.com/kanersps/loop/models"

// TODO: fix memory leak with assignments to existing objects
func NewEnvironment() *models.Environment {
	s := make(map[string]models.Object)
	return &models.Environment{Store: s}
}

func NewEnclosedEnvironment(outer *models.Environment) *models.Environment {
	env := NewEnvironment()
	env.Outer = outer

	return env
}
