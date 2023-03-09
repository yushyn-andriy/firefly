package evaluator

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yushyn-andriy/firefly/ast"
	"github.com/yushyn-andriy/firefly/object"
	"github.com/yushyn-andriy/firefly/token"
)

func init() {
	registerBuiltin("len", blen)
	registerBuiltin("first", bfirst)

	registerBuiltin("print", bprint)
	registerBuiltin("println", bprintln)
	registerBuiltin("eprint", beprint)
	registerBuiltin("eprintln", beprintln)

	registerBuiltin("exit", bexit)
	registerBuiltin("locals", blocals)
	registerBuiltin("type", btype)
	registerBuiltin("builtins", bbuiltins)
	registerBuiltin("getattr", bgetattr)
	registerBuiltin("setattr", bsetattr)
	registerBuiltin("nclass", bnclass)
}

func registerBuiltin(
	name string,
	f func(env *object.Environment, args ...object.Object) object.Object,
) {
	builtins[name] = &object.Builtin{Fn: f, Env: nil}
}

var builtins = map[string]*object.Builtin{}

func bnclass(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	var name string
	switch arg := args[0].(type) {
	case *object.String:
		name = arg.Value
	default:
		return newError("argument to `nclass` not supported, got %s",
			args[0].Type())
	}

	return object.NewClass(&ast.Identifier{Token: token.Token{
		Type:    token.IDENT,
		Literal: name,
	}, Value: name}, nil, env)
}

func blen(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return arg.Len()
	case *object.Array:
		return arg.Len()
	case *object.Instance:
		r := arg.Len()
		switch r := r.(type) {
		case *object.Function:
			r.Self = &args[0]
			return applyFunction(r, args)
		}
		return arg.Len()
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}

}

func bsetattr(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got=%d, want=3",
			len(args))
	}
	obj := args[0]
	keyObj := args[1]
	value := args[2]

	key, ok := keyObj.(*object.String)
	if !ok {
		return newError("TypeError: getattr(): attribute name must be string")
	}
	return obj.SetAttr(key.Value, value)
}

func bgetattr(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	obj := args[0]
	keyObj := args[1]

	key, ok := keyObj.(*object.String)
	if !ok {
		return newError("TypeError: getattr(): attribute name must be string")
	}

	return obj.GetAttr(key.Value)
}

func bbuiltins(env *object.Environment, args ...object.Object) object.Object {
	arr := object.Array{}
	for k, _ := range builtins {
		arr.Elements = append(arr.Elements, &object.String{Value: k})
	}
	return &arr
}

func btype(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	return &object.ObjType{Value: string(args[0].Type())}
}

func blocals(env *object.Environment, args ...object.Object) object.Object {
	if env != nil {
		return env.ToHash()
	}
	return NULL
}

func bfirst(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to 'first' must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return NULL
}

func bprint(env *object.Environment, args ...object.Object) object.Object {
	var out bytes.Buffer
	for index, arg := range args {
		out.WriteString(arg.Inspect())
		if index+1 == len(args) {
			continue
		}
		out.WriteString(" ")
	}
	fmt.Print(out.String())
	return NULL
}

func bprintln(env *object.Environment, args ...object.Object) object.Object {
	var out bytes.Buffer
	for index, arg := range args {
		out.WriteString(arg.Inspect())
		if index+1 == len(args) {
			continue
		}
		out.WriteString(" ")
	}
	out.WriteString("\n")
	fmt.Print(out.String())
	return NULL
}

func beprint(env *object.Environment, args ...object.Object) object.Object {
	var out bytes.Buffer
	for index, arg := range args {
		out.WriteString(arg.Inspect())
		if index+1 == len(args) {
			continue
		}
		out.WriteString(" ")
	}
	fmt.Fprint(os.Stderr, out.String())
	return NULL
}

func beprintln(env *object.Environment, args ...object.Object) object.Object {
	var out bytes.Buffer
	for index, arg := range args {
		out.WriteString(arg.Inspect())
		if index+1 == len(args) {
			continue
		}
		out.WriteString(" ")
	}
	out.WriteString("\n")
	fmt.Fprint(os.Stderr, out.String())
	return NULL
}

func bexit(env *object.Environment, args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		os.Exit(0)
	case 1:
		if args[0].Type() != object.INTEGER_OBJ {
			return newError("argument to 'exit' must be INT, got %s",
				args[0].Type())

		}
		status := args[0].(*object.Integer).Value
		os.Exit(int(status))
	default:
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	return NULL
}
