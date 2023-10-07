package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Array struct {
	dict     map[string]Object
	Elements []Object
}

func NewArray(elements []Object) *Array {
	s := new(Array)
	s.Elements = elements
	s.dict = make(map[string]Object)

	s.initialize()

	return s
}

func (ao *Array) initialize() {}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (ao *Array) Len() Object {
	return &Integer{Value: int64(len(ao.Elements))}
}
func (ao *Array) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", ao.Inspect(), key)}
}
func (ao *Array) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", ao.Inspect(), key)}
}
