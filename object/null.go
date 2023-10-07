package object

import "fmt"

type Null struct {
}

func NewNull() *Null {
	return new(Null)
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", n.Inspect(), key)}
}
func (n *Null) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", n.Inspect(), key)}
}
