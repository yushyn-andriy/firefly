package object

import "fmt"

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
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

func (e *Environment) ToHash() *Hash {
	pairs := map[HashKey]HashPair{}
	for k, v := range e.store {
		hashKey := &String{Value: k}

		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: hashKey, Value: v}
	}
	hash := &Hash{Pairs: pairs}

	return hash
}

func (e *Environment) Type() ObjectType { return HASH_OBJ }
func (e *Environment) Inspect() string {
	return e.ToHash().Inspect()
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	// e.store[name] = val

	curr := e
	var found bool
	for curr != nil {
		if _, ok := curr.store[name]; ok {
			found = true
			break
		}
		curr = curr.outer
	}
	if found {
		curr.store[name] = val
	} else {
		e.store[name] = val
	}

	return val
}

func (e *Environment) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", e.Inspect(), key)}
}
func (e *Environment) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", e.Inspect(), key)}
}
