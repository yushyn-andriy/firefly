package object

import "fmt"

type Instance struct {
	class *Class
	dict  map[string]Object
}

func (i *Instance) Type() ObjectType {
	return INSTANCE
}
func (i *Instance) Len() Object {
	flen, ok := i.class.dict[MEGIC_METHOD_LEN]
	if !ok {
		return newError("TypeError: object of type '%s' has no len()", i.class.name)
	}
	return flen
}
func (i *Instance) Inspect() string { return fmt.Sprintf("<'%s' object>", i.class.name) }
func (i *Instance) SetAttr(key string, value Object) Object {
	i.dict[key] = value
	return NULL
}
func (i *Instance) GetAttr(key string) Object {
	v, ok := i.dict[key]
	if !ok {
		v, ok := i.class.dict[key]
		if !ok {
			return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", i.Inspect(), key)}
		}
		return v
	}
	return v
}
