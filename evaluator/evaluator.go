package evaluator

import (
	"bytes"
	"fmt"
	"math"
	"staq/ast"
	"staq/object"
	"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func GetBuiltins() map[string]*object.Builtin {
	var builtins = map[string]*object.Builtin{
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch arg := args[0].(type) {
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"first": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return NULL
			},
		},
		"last": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}
				return NULL
			},
		},
		"rest": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arr.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return NULL
			},
		},
		"push": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			},
		},
		"pop": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arr.Elements[0:length-1])
					return &object.Array{Elements: newElements}
				}
				return NULL
			},
		},
		"forEach": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `forEach` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ {
					return newError("argument to `forEach` must be FUNCTION, got %s", args[1].Type())
				}
				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)
				for _, el := range arr.Elements {
					evaluated := Eval(fn.Body, extendFunctionEnv(fn, []object.Object{el}))
					if isError(evaluated) {
						return evaluated
					}
				}
				return NULL
			},
		},
		"map": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `map` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ {
					return newError("argument to `map` must be FUNCTION, got %s", args[1].Type())
				}
				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)
				var elements []object.Object
				for _, el := range arr.Elements {
					evaluated := Eval(fn.Body, extendFunctionEnv(fn, []object.Object{el}))
					if isError(evaluated) {
						return evaluated
					}
					elements = append(elements, evaluated)
				}
				return &object.Array{Elements: elements}
			},
		},
		"filter": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `filter` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ {
					return newError("argument to `filter` must be FUNCTION, got %s", args[1].Type())
				}
				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)
				var elements []object.Object
				for _, el := range arr.Elements {
					evaluated := Eval(fn.Body, extendFunctionEnv(fn, []object.Object{el}))
					if isError(evaluated) {
						return evaluated
					}
					if isTruthy(evaluated) {
						elements = append(elements, el)
					}
				}
				return &object.Array{Elements: elements}
			},
		},
		"reduce": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("wrong number of arguments. got=%d, want=3", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `reduce` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ {
					return newError("argument to `reduce` must be FUNCTION, got %s", args[1].Type())
				}
				if args[2].Type() != object.INTEGER && args[2].Type() != object.FLOAT {
					return newError("argument to `reduce` must be INTEGER or FLOAT, got %s", args[2].Type())
				}
				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)
				acc := args[2]
				for _, el := range arr.Elements {
					evaluated := Eval(fn.Body, extendFunctionEnv(fn, []object.Object{acc, el}))
					if isError(evaluated) {
						return evaluated
					}
					acc = evaluated
				}
				return acc
			},
		},
	}
	return builtins
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	// Expressions
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
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return evalString(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
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
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}
	return nil
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "~":
		return evalBitwiseNotPrefixOperatorExpression(right)
	case "++":
		return evalIncrementPrefixOperatorExpression(right)
	case "--":
		return evalDecrementPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT && right.Type() == object.FLOAT:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER && right.Type() == object.FLOAT:
		return evalFloatInfixExpression(operator, &object.Float{Value: float64(left.(*object.Integer).Value)}, right)
	case left.Type() == object.FLOAT && right.Type() == object.INTEGER:
		return evalFloatInfixExpression(operator, left, &object.Float{Value: float64(right.(*object.Integer).Value)})
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(strings.Compare(leftVal, rightVal) == 0)
	case "!=":
		return nativeBoolToBooleanObject(strings.Compare(leftVal, rightVal) != 0)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER && right.Type() != object.FLOAT {
		return newError("unknown operator: -%s", right.Type())
	}
	switch right.Type() {
	case object.INTEGER:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FLOAT:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	default:
		return NULL
	}
}

func evalBitwiseNotPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: ^value}
}

func evalIncrementPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER && right.Type() != object.FLOAT {
		return NULL
	}
	switch right.Type() {
	case object.INTEGER:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: value + 1}
	case object.FLOAT:
		value := right.(*object.Float).Value
		return &object.Float{Value: value + 1}
	default:
		return NULL
	}
}

func evalDecrementPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER && right.Type() != object.FLOAT {
		return NULL
	}
	switch right.Type() {
	case object.INTEGER:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: value - 1}
	case object.FLOAT:
		value := right.(*object.Float).Value
		return &object.Float{Value: value - 1}
	default:
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	builtins := GetBuiltins()
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

func evalString(node *ast.StringLiteral, env *object.Environment) object.Object {
	// Interpolate the values between brackets ({}) to their corresponding values from env
	// Example: "Hello, {name}!" -> "Hello, World!"
	// If the value is not found in env, create an error
	// Example: "Hello, {name}!" -> "Hello, ERROR: identifier not found: name!"
	var out bytes.Buffer
	var inBrackets bool
	var identifier string
	for _, c := range node.Value {
		if c == '{' {
			inBrackets = true
			continue
		}
		if c == '}' {
			inBrackets = false
			val, ok := env.Get(identifier)
			if !ok {
				return newError("identifier not found: %s", identifier)
			}
			out.WriteString(val.Inspect())
			identifier = ""
			continue
		}
		if inBrackets {
			identifier += string(c)
			continue
		}
		out.WriteString(string(c))
	}
	return &object.String{Value: out.String()}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
	case "//":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "**":
		return &object.Integer{Value: intPow(leftVal, rightVal)}
	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
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
	case "//":
		return &object.Float{Value: leftVal / rightVal}
	case "%":
		return NULL
	case "**":
		return &object.Float{Value: math.Pow(leftVal, rightVal)}
	case "<<":
		return NULL
	case ">>":
		return NULL
	case "&":
		return NULL
	case "|":
		return NULL
	case "^":
		return NULL
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalArrayIndexExpression returns the element at the given index of the array
// If the index is out of bounds, it returns NULL
// If the index is negative, it counts from the end of the array, starting from -1.
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 {
		idx = max + 1 + idx
	}
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
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

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func intShiftLeft(x, y int64) int64 {
	return x << y
}

func intPow(x, y int64) int64 {
	if y == 0 {
		return 1
	}
	if y == 1 {
		return x
	}
	result := x
	for i := int64(2); i <= y; i++ {
		result *= x
	}
	return result
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
