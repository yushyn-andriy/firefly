package object

import "fmt"

type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Fn   BuiltinFunction
	Env  *Environment
	Self Object
	Doc  string
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", b.Inspect(), key)}
}
func (b *Builtin) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", b.Inspect(), key)}
}
