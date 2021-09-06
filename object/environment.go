package object

// TODO: fix memory leak with assignments to existing objects
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return e.store[name]
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}
