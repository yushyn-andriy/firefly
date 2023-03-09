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
	flen, ok := i.class.dict[MAGIC_METHOD_LEN]
	if !ok {
		return i.class.GetAttr(MAGIC_METHOD_LEN)
	}
	return flen
}

func (i *Instance) Inspect() string { return fmt.Sprintf("<'%s' object>", i.class.Name.Value) }

func (i *Instance) SetAttr(key string, value Object) Object {
	i.dict[key] = value
	return NULL
}

func (i *Instance) GetAttr(key string) Object {
	v, ok := i.dict[key]
	if !ok {
		return i.class.GetAttr(key)
	}
	return v
}
