package object

import (
	"fmt"
	"hash/fnv"
	"strings"
)

type String struct {
	dict  map[string]Object
	Value string
}

func NewString(value string) *String {
	s := new(String)
	s.Value = value
	s.dict = make(map[string]Object)

	s.initialize()

	return s
}

func (s *String) initialize() {
	s.dict["reverse"] = &Builtin{
		Fn:   strReverse,
		Env:  nil,
		Self: s,
		Doc: `reverse(self)
reverse the string and return new string object 
`,
	}
	s.dict["upper"] = &Builtin{
		Fn:   strUpper,
		Env:  nil,
		Self: s,
		Doc: `upper(self)
upper case the string and return new string object 
`,
	}
	s.dict["split"] = &Builtin{
		Fn:   strSplit,
		Env:  nil,
		Self: s,
		Doc: `split(self)
split the string by separator and return new array of string objects.
The default separator is a whitespace.  
`,
	}
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
	if val, ok := s.dict[key]; ok {
		return val
	}
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", s.Inspect(), key)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func strReverse(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	self, _ := args[0].(*String)

	tmp := []rune(self.Value)
	s := make([]rune, len(tmp))
	for i := len(tmp) - 1; i >= 0; i-- {
		s = append(s, tmp[i])
	}

	r := NewString(string(s))

	return r

}

func strUpper(env *Environment, args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	s := args[0].(*String)
	r := NewString(strings.ToUpper(s.Value))

	return r
}

func strSplit(env *Environment, args ...Object) Object {
	var sep = " "
	if len(args) < 1 || len(args) > 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if len(args) == 2 {
		sep = args[1].(*String).Value
	}

	s := args[0].(*String).Value
	array := strings.Split(s, sep)

	elements := make([]Object, len(array))
	for i, element := range array {
		elements[i] = NewString(element)
	}

	return NewArray(elements)
}
