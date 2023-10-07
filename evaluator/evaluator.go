package evaluator

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/yushyn-andriy/firefly/ast"
	"github.com/yushyn-andriy/firefly/lexer"
	"github.com/yushyn-andriy/firefly/object"
	"github.com/yushyn-andriy/firefly/parser"
)

var (
	NULL  = object.NULL
	TRUE  = object.TRUE
	FALSE = object.FALSE
)

func readFullFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ImportLiteral:
		name := node.Name.Value

		// find better way to do this
		input, err := readFullFile("./lib/" + name + ".fl")
		if err != nil {
			log.Fatal(err)
		}
		l := lexer.New(string(input))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(os.Stderr, p.Errors())
		}

		moduleEnv := object.NewEnvironment()
		Eval(program, moduleEnv)

		// fmt.Println(moduleEnv.Inspect())

		m := object.NewModule(node.Name, moduleEnv)
		env.Set(name, m)
		return m
	case *ast.SelectorExpr:
		return evalSelectorExpression(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		name := node.Name
		function := object.NewFunction(name, params, env, body)
		if name != nil {
			env.Set(name.String(), function)
		}
		return function

	case *ast.ClassLiteral:
		body := node.Body
		name := node.Name
		if name == nil {
			return nil
		}
		cls := object.NewClass(name, body, env)

		extendedEnv := extendClassEnv(cls, nil)
		cls.Env = extendedEnv

		err := evalBlockStatement(body, extendedEnv)
		if isError(err) {
			return err
		}
		if name != nil {
			env.Set(name.String(), cls)
		}
		return cls
	case *ast.ForStatement:
		loop := &object.ForLoop{
			Init: node.Init.(*ast.AssignStatement),
			Cond: node.Cond,
			Post: node.Post.(*ast.AssignStatement),
			Body: node.Body,
		}
		return runForLoop(loop, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.StringLiteral:
		return object.NewString(node.Value)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return object.NewArray(elements)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		// this means that this is assign statement
		if node.Right != nil {
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			obj := evalAssignIndexStatement(left, right, index)
			if isError(obj) {
				return obj
			}
		} else {
			return evalIndexExpression(left, index)
		}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func evalSelectorExpression(node *ast.SelectorExpr, env *object.Environment) object.Object {
	obj := Eval(node.Expression, env)
	if node.Value != nil {
		value := Eval(node.Value, env)
		return obj.SetAttr(node.Selector.Value, value)
	}

	res := obj.GetAttr(node.Selector.Value)
	switch res := res.(type) {
	case *object.Function:
		res.Self = obj
		return res
	}
	return res
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalAssignIndexStatement(left, right, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalAssignArrayIndexStatement(left, right, index)
	case left.Type() == object.HASH_OBJ:
		return evalAssignHashIndexStatement(left, right, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalAssignHashIndexStatement(hash, value, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	hashObject.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: value}
	return NULL
}

func evalAssignArrayIndexStatement(array, value, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements)) - 1

	if idx < 0 || idx > max {
		return NULL
	}
	arrayObject.Elements[idx] = value
	return NULL
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return newError("key does not exists: %s", index.Inspect())
	}
	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements)) - 1

	if idx < 0 || idx > max {
		return newError("index out of range: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(obj, index object.Object) object.Object {
	stringObject := obj.(*object.String)
	idx := index.(*object.Integer).Value

	runes := []rune(stringObject.Value)

	max := int64(len(runes)) - 1
	if idx < 0 || idx > max {
		return newError("index out of range: %d", idx)
	}
	return object.NewString(string(runes[idx]))
}

func runForLoop(loop object.Object, env *object.Environment) object.Object {
	switch loop := loop.(type) {
	case *object.ForLoop:
		loop.Env = env
		extendedEnv := extendForLoopEnv(loop, []object.Object{})
		initVal := Eval(loop.Init, extendedEnv)
		if isError(initVal) {
			return initVal
		}

		condExpr := Eval(loop.Cond, extendedEnv)
		if isError(condExpr) {
			return condExpr
		}

		for isTruthy(condExpr) {
			evaluated := Eval(loop.Body, extendedEnv)
			if isError(evaluated) {
				return evaluated
			}

			post := Eval(loop.Post, extendedEnv)
			if isError(post) {
				return post
			}

			condExpr = Eval(loop.Cond, extendedEnv)
			if isError(condExpr) {
				return condExpr
			}
		}

		return nil
	default:
		return newError("not a ForLoop: %s", loop.Type())
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Class:
		obj := fn.NewInstance(args...)
		init := obj.GetAttr(object.MAGIC_METHOD_INIT)
		switch r := init.(type) {
		case *object.Error:
		case *object.Function:
			r.Self = obj
			applyFunction(r, args)
		}
		return obj
	case *object.Builtin:
		var res object.Object
		// If fn.Self is not nil that means
		// that this is a function realization to an object
		// and we need pass it to a function call as a first argument
		if fn.Self != nil {
			extended := []object.Object{fn.Self}
			extended = append(extended, args...)
			res = fn.Fn(fn.Env, extended...)
		} else {
			res = fn.Fn(fn.Env, args...)
		}
		switch res := res.(type) {
		case *object.Function:
			if len(args) != len(res.Parameters) {
				return newError("TypeError: expected %d arguments got %d", len(args), len(res.Parameters))
			}
			extendedEnv := extendFunctionEnv(res, args)
			evaluated := Eval(res.Body, extendedEnv)
			return unwrapReturnValue(evaluated)
		default:
			return res
		}
	default:
		return newError("not a function: %s", fn.Type())
	}

}

func extendForLoopEnv(
	loop *object.ForLoop,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(loop.Env)
	return env
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	if fn.Self != nil {
		env.Set("self", fn.Self)
	}
	return env
}

func extendClassEnv(
	cls *object.Class,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(cls.Env)
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		builtin.Env = env
		return builtin
	}

	if namedFunc, ok := env.Get(node.Value); ok {
		return namedFunc
	}

	return newError("identifier not found: " + node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringIntegerInfix(operator, left, right)
	case right.Type() == object.INTEGER_OBJ && left.Type() == object.STRING_OBJ:
		return evalStringIntegerInfix(operator, right, left)

	case operator == "==" && left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=" && left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return nativeBoolToBooleanObject(left != right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case operator == "and":
		return nativeBoolToBooleanObject(object.TRUE == left && object.TRUE == right)

	case operator == "or":
		return nativeBoolToBooleanObject(object.TRUE == left || object.TRUE == right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return object.NewString(leftVal + rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringIntegerInfix(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "*":
		return object.NewString(strings.Repeat(rightVal, int(leftVal)))
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}
func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// deprecated
// func evalStatements(stmts []ast.Statement) object.Object {
//	var result object.Object

//	for _, statement := range stmts {
//		result = Eval(statement)

//		if returnValue, ok := result.(*object.ReturnValue); ok {
//			return returnValue.Value
//		}
//	}

//	return result

// }

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
