package object

import (
	"bytes"
	"fmt"

	"github.com/yushyn-andriy/firefly/ast"
)

type ForLoop struct {
	Init *ast.AssignStatement
	Cond ast.Expression
	Post *ast.AssignStatement
	Body *ast.BlockStatement
	Env  *Environment
}

func (fl *ForLoop) Type() ObjectType { return FORLOOP_OBJ }
func (fl *ForLoop) Inspect() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString("(")
	out.WriteString(fl.Init.String())
	out.WriteString(fl.Cond.String())
	out.WriteString(";")
	out.WriteString(fl.Post.String())
	out.WriteString(")")
	out.WriteString("{")
	if fl.Body != nil {
		out.WriteString(fl.Body.String())
	}
	out.WriteString("}")

	return out.String()
}
func (fl *ForLoop) SetAttr(key string, value Object) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", fl.Inspect(), key)}
}

func (fl *ForLoop) GetAttr(key string) Object {
	return &Error{Message: fmt.Sprintf("AttributeError: '%s' object has no attribute  %s", fl.Inspect(), key)}
}
