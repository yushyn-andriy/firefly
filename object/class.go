package object

import (
	"fmt"

	"github.com/yushyn-andriy/firefly/ast"
)

type Class struct {
	Name *ast.Identifier
	Body *ast.BlockStatement
	Env  *Environment

	dict map[string]Object
}

func NewClass(
	name *ast.Identifier,
	body *ast.BlockStatement,
	env *Environment,
) *Class {
	cls := new(Class)

	cls.Name = name
	cls.Body = body
	cls.Env = env
	cls.dict = make(map[string]Object)

	cls.initialize()

	return cls
}

func (c *Class) initialize() {
	// initialize attributes
	c.Env.Set(MAGIC_ATTR_NAME, NewString(c.Name.Value))
	// TODO: add doc string parsing
	c.Env.Set(MAGIC_ATTR_DOC, NewString(""))
}

func (c *Class) Type() ObjectType {
	return CLASS
}
func (c *Class) Inspect() string { return fmt.Sprintf("<class '%s'>", c.Name.Value) }
func (c *Class) SetAttr(key string, value Object) Object {
	c.dict[key] = value
	return NULL
}
func (c *Class) GetAttr(key string) Object {
	v, ok := c.dict[key]
	if !ok {
		v, ok := c.Env.Get(key)
		if !ok {
			return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", c.Inspect(), key)}
		}
		return v
	}
	return v
}

func (c *Class) NewInstance(args ...Object) *Instance {
	self := new(Instance)
	self.class = c
	self.dict = make(map[string]Object)
	return self
}
