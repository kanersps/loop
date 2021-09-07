package helpers

import (
	"fmt"
	"github.com/kanersps/loop/ast"
	"github.com/kanersps/loop/models"
	"github.com/kanersps/loop/object"
	"github.com/kanersps/loop/object/builtins"
)

func Eval(node ast.Node, env *models.Environment) models.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &models.Integer{Value: node.Value}
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

		return &models.Return{Value: value}
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
		return &models.Function{
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

		return ApplyFunction(function, args, env)
	case *ast.StringLiteral:
		return &models.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &models.Array{Elements: elements}
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

func throwError(format string, a ...interface{}) *models.Error {
	return &models.Error{Message: fmt.Sprintf(format, a...)}
}

func evalHashLiteral(node *ast.HashLiteral, env *models.Environment) models.Object {
	pairs := make(map[models.HashKey]models.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(models.Hashable)
		if !ok {
			return throwError("HASHMAP KEY IS INCORRECT TYPE. got=%s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = models.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return &models.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index models.Object) models.Object {
	if left.Type() == models.ARRAY {
		if index.Type() != models.INTEGER {
			return throwError("INVALID INDEX. expected=INTEGER. got=%s", index.Type())
		}

		return evalArrayIndexExpression(left, index)
	}

	if left.Type() == models.HASH {
		return evalHashIndexExpression(left, index)
	}

	return throwError("ATTEMPTED INDEXING INVALID TYPE %s", left.Type())
}

func evalHashIndexExpression(hash, index models.Object) models.Object {
	hashObj := hash.(*models.Hash)

	idx := index.(models.Hashable).HashKey()

	pair, ok := hashObj.Pairs[idx]
	if !ok {
		return models.NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index models.Object) models.Object {
	arrayObj := array.(*models.Array)
	idx := index.(*models.Integer).Value

	return arrayObj.Elements[idx]
}

func ApplyFunction(fn models.Object, args []models.Object, env *models.Environment) models.Object {
	switch fn := fn.(type) {
	case *models.Function:
		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *models.Builtin:
		builtins.SetApplyFunction(ApplyFunction)
		return fn.Func(env, args...)
	default:
		return throwError("UNKNOWN-FUNCTION: %s", fn.Type())
	}
}

func extendedFunctionEnv(fn *models.Function, args []models.Object) *models.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramId, param := range fn.Parameters {
		env.Set(param.Value, args[paramId])
	}

	return env
}

func unwrapReturnValue(obj models.Object) models.Object {
	if returnValue, ok := obj.(*models.Return); ok {
		return returnValue.Value
	}

	return obj
}

func evalExpressions(expressions []ast.Expression, env *models.Environment) []models.Object {
	var results []models.Object

	for _, exp := range expressions {
		evaluated := Eval(exp, env)

		if isError(evaluated) {
			return []models.Object{evaluated}
		}

		results = append(results, evaluated)
	}

	return results
}

func evalIdentifier(node *ast.Identifier, env *models.Environment) models.Object {
	value, ok := env.Get(node.Value)
	if !ok {
		if builtin, ok := builtins.Functions[node.Value]; ok {
			return builtin
		}

		return throwError("UNKNOWN-IDENTIFIER: %s", node.Value)
	}

	return value
}

func evalIfExpression(node *ast.IfExpression, env *models.Environment) models.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if condition == models.TRUE {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return models.NULL
	}
}

func evalWhileExpression(node *ast.WhileLiteral, env *models.Environment) models.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	var lastEvaluation models.Object

	for condition == models.TRUE {
		lastEvaluation = Eval(node.Body, env)
		condition = Eval(node.Condition, env)
	}

	return lastEvaluation
}

func evalInfixExpression(operator string, left models.Object, right models.Object) models.Object {
	if left.Type() == models.INTEGER && right.Type() == models.INTEGER {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if operator == "+" && left.Type() == models.STRING && right.Type() == models.STRING {
		return &models.String{Value: left.Inspect() + right.Inspect()}
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

func evalIntegerInfixExpression(operator string, left, right models.Object) models.Object {
	lv := left.(*models.Integer).Value
	rv := right.(*models.Integer).Value

	switch operator {
	case "+":
		return &models.Integer{Value: lv + rv}
	case "*":
		return &models.Integer{Value: lv * rv}
	case "/":
		return &models.Integer{Value: lv / rv}
	case "-":
		return &models.Integer{Value: lv - rv}
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

func evalBlockStatement(block *ast.BlockStatement, env *models.Environment) models.Object {
	var result models.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == models.RETURN || rt == models.ERROR {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right models.Object) models.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return models.NULL
	}
}

func evalBangOperatorExpression(right models.Object) models.Object {
	switch right {
	case models.TRUE:
		return models.FALSE
	case models.FALSE:
		return models.TRUE
	case models.NULL:
		return models.TRUE
	default:
		return models.FALSE
	}
}

func evalMinusPrefixOperator(right models.Object) models.Object {
	if right.Type() != models.INTEGER {
		return throwError("UNKNOWN-OPERATOR: -%s", right.Type())
	}

	return &models.Integer{Value: -right.(*models.Integer).Value}
}

func isError(obj models.Object) bool {
	if obj != nil {
		return obj.Type() == models.ERROR
	}

	return false
}

func nativeBoolToBooleanObject(input bool) *models.Boolean {
	if input {
		return models.TRUE
	}
	return models.FALSE
}

func evalProgram(stmts []ast.Statement, env *models.Environment) models.Object {
	var result models.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result.(type) {
		case *models.Return:
			return result.(*models.Return).Value
		case *models.Error:
			return result
		}
	}

	return result
}
