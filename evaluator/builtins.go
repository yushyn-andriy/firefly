package evaluator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/yushyn-andriy/firefly/ast"
	"github.com/yushyn-andriy/firefly/object"
	"github.com/yushyn-andriy/firefly/token"
)

var (
	stdout io.Writer
	stderr io.Writer
)

func init() {
	stdout = os.Stdout
	stderr = os.Stderr
}

func SetStd(sout, serr io.Writer) {
	stdout = sout
	stderr = serr
}

func init() {
	registerBuiltin("len", blen)
	registerBuiltin("first", bfirst)

	registerBuiltin("print", bprint)
	registerBuiltin("utskrift", bprint)

	registerBuiltin("println", bprintln)
	registerBuiltin("utskriftln", bprintln)

	registerBuiltin("printf", bprintf)

	registerBuiltin("eprint", beprint)
	registerBuiltin("futskrift", beprintln)

	registerBuiltin("exit", bexit)
	registerBuiltin("locals", blocals)
	registerBuiltin("type", btype)
	registerBuiltin("builtins", bbuiltins)
	registerBuiltin("getattr", bgetattr)
	registerBuiltin("setattr", bsetattr)
	registerBuiltin("new", bNewClass)
	registerBuiltin("help", bhelp)

	registerBuiltin("pow", bPow)
	registerBuiltin("file", bNewFile)
	registerBuiltin("input", bInput)
	registerBuiltin("system", bSystem)

	registerBuiltin("int", bInt)
	registerBuiltin("float", bFloat)
	registerBuiltin("string", bString)
}

func registerBuiltin(
	name string,
	f func(env *object.Environment, args ...object.Object) object.Object,
) {
	builtins[name] = &object.Builtin{Fn: f, Env: nil}
}

var builtins = map[string]*object.Builtin{}

func bFloat(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	obj := args[0]
	switch arg := obj.(type) {
	case *object.String:
		number, err := strconv.ParseFloat(arg.Value, 64)
		if err != nil {
			return newError("%s", err)
		}
		return object.NewFloat(number)
	case *object.Integer:
		return object.NewFloat(float64(arg.Value))
	default:
		return newError("invalid object type %T", arg)
	}
}

func bString(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	obj := args[0]
	switch arg := obj.(type) {
	case *object.Float:
		return object.NewString(fmt.Sprintf("%f", arg.Value))
	case *object.Integer:
		return object.NewString(fmt.Sprintf("%d", arg.Value))
	default:
		return newError("invalid object type %T", arg)
	}
}

func bInt(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	obj := args[0]
	switch arg := obj.(type) {
	case *object.String:
		number, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return newError("%s", err)
		}
		return object.NewInteger(number)
	case *object.Float:
		return object.NewInteger(int64(arg.Value))
	default:
		return newError("invalid object type %T", arg)
	}
}

func bSystem(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, minimum=1",
			len(args))
	}

	name, _ := args[0].(*object.String)
	args = args[1:]

	strArguments := []string{}
	for _, arg := range args {
		s := arg.(*object.String).Value
		strArguments = append(strArguments, s)
	}

	command := exec.Command(name.Value, strArguments...)

	out, err := command.Output()
	if err != nil {
		return newError("could not run command: %s", err)
	}

	return object.NewString(string(out))
}

func bInput(env *object.Environment, args ...object.Object) object.Object {
	if len(args) > 1 {
		return newError("wrong number of arguments. got=%d, want=0",
			len(args))
	}

	if len(args) == 1 {
		fmt.Fprint(stdout, args[0].Inspect())
	}

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.ReplaceAll(line, "\n", "")

	return object.NewString(string(line))
}

func bNewFile(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	var path, mode string
	switch arg := args[0].(type) {
	case *object.String:
		path = arg.Value
	default:
		return newError("argument to `nclass` not supported, got %s",
			args[0].Type())
	}
	switch arg := args[1].(type) {
	case *object.String:
		mode = arg.Value
	default:
		return newError("argument to `nclass` not supported, got %s",
			args[0].Type())
	}

	return object.NewFile(path, mode)
}

func bNewClass(env *object.Environment, args ...object.Object) object.Object {
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

func bPow(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	a := args[0]
	b := args[1]
	if a.Type() != object.FLOAT_OBJ || b.Type() != object.FLOAT_OBJ {
		return newError("both arguments must be FLOAT  type got %s, %s",
			a.Type(), b.Type())
	}

	x := a.(*object.Float)
	y := b.(*object.Float)

	return &object.Float{Value: math.Pow(x.Value, y.Value)}
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
			r.Self = args[0]
			return applyFunction(r, args)
		}
		return arg.Len()
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}

}

func bhelp(env *object.Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *object.Builtin:
		return object.NewString(arg.Doc)
	case *object.Array:
		return arg.Len()
	default:
		return arg.GetAttr("__doc__")
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
	arr := object.NewArray(nil)
	for k, _ := range builtins {
		arr.Elements = append(arr.Elements, object.NewString(k))
	}
	return arr
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

func bprintf(env *object.Environment, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, minimum=1",
			len(args))
	}

	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to 'first' must be %s, got %s",
			object.STRING_OBJ, args[0].Type())
	}

	format := args[0].(*object.String).Value
	arguments := []any{}
	for _, arg := range args[1:] {
		arguments = append(arguments, arg.Inspect())
	}

	format = strings.ReplaceAll(format, "\\n", "\n")

	_, err := fmt.Fprintf(stdout, format, arguments...)
	if err != nil {
		return newError("%s", err)
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
	fmt.Fprint(stdout, out.String())
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
	fmt.Fprint(stdout, out.String())
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
	fmt.Fprint(stderr, out.String())
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
	fmt.Fprint(stderr, out.String())
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
