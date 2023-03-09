package object

import "fmt"

type Instance struct {
	class *Class
	dict  map[string]Object
}

func (i *Instance) Type() ObjectType {
	return INSTANCE
}

func (i *Instance) Inspect() string { return fmt.Sprintf("<'%s' object>", i.class.name) }
func (i *Instance) SetAttr(key string, value Object) Object {
	i.dict[key] = value
	return NULL
}
func (i *Instance) GetAttr(key string) Object {
	v, ok := i.dict[key]
	if !ok {
		return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", i.Inspect(), key)}
	}
	return v
}
