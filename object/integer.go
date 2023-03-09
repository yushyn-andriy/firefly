package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (i *Integer) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", i.Inspect(), key)}
}

func (i *Integer) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", i.Inspect(), key)}
}
