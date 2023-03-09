package object

import "fmt"

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", rv.Inspect(), key)}
}

func (rv *ReturnValue) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", rv.Inspect(), key)}
}
