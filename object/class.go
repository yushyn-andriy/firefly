package object

import "fmt"

type Class struct {
	name string
	dict map[string]Object
}

func NewClass(name string) *Class {
	cls := new(Class)
	cls.dict = make(map[string]Object)
	cls.name = name
	return cls
}

func (c *Class) Type() ObjectType {
	return CLASS
}
func (c *Class) Inspect() string { return fmt.Sprintf("<class '%s'>", c.name) }
func (c *Class) SetAttr(key string, value Object) Object {
	c.dict[key] = value
	return NULL
}
func (c *Class) GetAttr(key string) Object {
	v, ok := c.dict[key]
	if !ok {
		return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", c.Inspect(), key)}
	}
	return v
}

func (c *Class) NewInstance(args ...Object) *Instance {
	self := new(Instance)
	self.class = c
	self.dict = make(map[string]Object)
	return self
}
