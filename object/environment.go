package object

type Environment struct {
	store  map[string]Object
	global map[string]Object
	outer  *Environment
	Scope  bool
}

func NewEncloseEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]Object), global: nil, outer: outer, Scope: false}
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), global: make(map[string]Object), outer: nil, Scope: true}
}

func NewEnvironmentFn(env *Environment) *Environment {
	return &Environment{store: make(map[string]Object), global: env.global, outer: nil, Scope: false}
}

func (e *Environment) Get(name string) (Object, bool) {

	obj, ok := e.global[name]
	if !ok {
		obj, ok = e.store[name]
		if !ok && e.outer != nil {
			obj, ok = e.outer.Get(name)
		}
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {

	_, ok := e.global[name]
	if !ok {
		_, ok = e.store[name]
		if ok {
			e.store[name] = value
		} else if e.outer != nil {
			e.outer.Set(name, value)
		}
	} else {
		e.global[name] = value
	}

	return value
	// e.store[name] = value
	// return value
}

func (e *Environment) Save(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Environment) SaveGlobal(name string, value Object) Object {
	e.global[name] = value
	return value
}
