package object

import (
	"fmt"

	"github.com/yushyn-andriy/firefly/ast"
)

type Module struct {
	Name *ast.StringLiteral
	Env  *Environment
	dict map[string]Object
}

func NewModule(name *ast.StringLiteral, env *Environment) *Module {
	m := new(Module)
	m.Name = name
	m.Env = env
	m.dict = make(map[string]Object)
	return m
}

func (m *Module) Type() ObjectType {
	return CLASS
}
func (m *Module) Inspect() string { return fmt.Sprintf("<class '%s'>", m.Name.Value) }
func (m *Module) SetAttr(key string, value Object) Object {
	m.dict[key] = value
	return NULL
}
func (m *Module) GetAttr(key string) Object {
	v, ok := m.dict[key]
	if !ok {
		v, ok := m.Env.Get(key)
		if !ok {
			return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", m.Inspect(), key)}
		}
		return v
	}
	return v
}
