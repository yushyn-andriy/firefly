package object

import (
	"fmt"
	"hash/fnv"
)

type String struct {
	dict  map[string]Object
	Value string
}

func NewString(value string) *String {
	s := new(String)
	s.Value = value
	s.dict = make(map[string]Object)
	return s
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", s.Inspect(), key)}
}
func (s *String) Len() Object {
	return &Integer{Value: int64(len(s.Value))}
}

func (s *String) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", s.Inspect(), key)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
