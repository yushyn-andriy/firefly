package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yushyn-andriy/firefly/ast"
)

type Function struct {
	dict       map[string]Object
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Self       Object
}

func NewFunction(name *ast.Identifier, params []*ast.Identifier, env *Environment, body *ast.BlockStatement) *Function {
	f := new(Function)
	f.dict = make(map[string]Object)
	f.Name = name
	f.Parameters = params
	f.Body = body
	f.Env = env
	return f
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	if f.Name != nil {
		out.WriteString(" ")
		out.WriteString(f.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) SetAttr(key string, value Object) Object {
	f.dict[key] = value
	return NULL
}
func (f *Function) GetAttr(key string) Object {
	v, ok := f.dict[key]
	if !ok {
		return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", f.Inspect(), key)}
	}
	return v
}
