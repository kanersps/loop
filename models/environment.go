package models

type Environment struct {
	Store map[string]Object
	Outer *Environment
}

func (e *Environment) Set(name string, value Object) Object {
	e.Store[name] = value
	return e.Store[name]
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]

	if !ok && e.Outer != nil {
		obj, ok = e.Outer.Get(name)
	}

	return obj, ok
}
