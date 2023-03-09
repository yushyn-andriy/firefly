package object

import "fmt"

type ObjType struct {
	Value string
}

func (ot *ObjType) Type() ObjectType {
	return TYPE_OBJ
}

func (ot *ObjType) Inspect() string {
	return ot.Value
}
func (ot *ObjType) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", ot.Inspect(), key)}
}

func (ot *ObjType) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", ot.Inspect(), key)}
}
