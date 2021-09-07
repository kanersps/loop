package evaluator

import (
	"fmt"
	"github.com/kanersps/loop/ast"
	"github.com/kanersps/loop/object"
	"github.com/kanersps/loop/object/builtins"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)

		if isError(left) {
			return left
		}

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}

		return &object.Return{Value: value}
	case *ast.WhileLiteral:
		return evalWhileExpression(node, env)
	case *ast.VariableStatement:
		value := Eval(node.Value, env)

		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
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
		return &object.String{Value: node.Value}
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
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}
	return nil
}

func throwError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return throwError("HASHMAP KEY IS INCORRECT TYPE. got=%s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object) object.Object {
	if left.Type() == object.ARRAY {
		if index.Type() != object.INTEGER {
			return throwError("INVALID INDEX. expected=INTEGER. got=%s", index.Type())
		}

		return evalArrayIndexExpression(left, index)
	}

	if left.Type() == object.HASH {
		return evalHashIndexExpression(left, index)
	}

	return throwError("ATTEMPTED INDEXING INVALID TYPE %s", left.Type())
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	idx := index.(object.Hashable).HashKey()

	pair, ok := hashObj.Pairs[idx]
	if !ok {
		return builtins.NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value

	return arrayObj.Elements[idx]
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Func(args...)
	default:
		return throwError("UNKNOWN-FUNCTION: %s", fn.Type())
	}
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramId, param := range fn.Parameters {
		env.Set(param.Value, args[paramId])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}

	return obj
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, exp := range expressions {
		evaluated := Eval(exp, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		results = append(results, evaluated)
	}

	return results
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(node.Value)
	if !ok {
		if builtin, ok := builtins.Functions[node.Value]; ok {
			return builtin
		}

		return throwError("UNKNOWN-IDENTIFIER: %s", node.Value)
	}

	return value
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if condition == builtins.TRUE {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return builtins.NULL
	}
}

func evalWhileExpression(node *ast.WhileLiteral, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	var lastEvaluation object.Object

	for condition == builtins.TRUE {
		lastEvaluation = Eval(node.Body, env)
		condition = Eval(node.Condition, env)
	}

	return lastEvaluation
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if operator == "+" && left.Type() == object.STRING && right.Type() == object.STRING {
		return &object.String{Value: left.Inspect() + right.Inspect()}
	}

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	}

	if left.Type() != right.Type() {
		return throwError("TYPE-MISMATCH: %s %s %s", left.Type(), operator, right.Type())
	}

	return throwError("UNKNOWN-OPERATOR: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "<":
		return nativeBoolToBooleanObject(lv < rv)
	case ">":
		return nativeBoolToBooleanObject(lv > rv)
	case "==":
		return nativeBoolToBooleanObject(lv == rv)
	case "!=":
		return nativeBoolToBooleanObject(lv != rv)
	}

	return throwError("UNKNOWN-OPERATOR: %s %s %s", left.Type(), operator, right.Type())
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return builtins.NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case builtins.TRUE:
		return builtins.FALSE
	case builtins.FALSE:
		return builtins.TRUE
	case builtins.NULL:
		return builtins.TRUE
	default:
		return builtins.FALSE
	}
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return throwError("UNKNOWN-OPERATOR: -%s", right.Type())
	}

	return &object.Integer{Value: -right.(*object.Integer).Value}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return builtins.TRUE
	}
	return builtins.FALSE
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result.(type) {
		case *object.Return:
			return result.(*object.Return).Value
		case *object.Error:
			return result
		}
	}

	return result
}
