package evaluator

import (
	"github.com/kanersps/loop/object"
	"github.com/kanersps/loop/parser"
	"github.com/kanersps/loop/parser/lexer"
	"testing"
)

func TestEval_IntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-12", -12},
		{"5 + 5", 10},
		{"5 + 5 * 2", 15},
		{"(5 + 5) * 2", 20},
		{"(5 + 5) / 2", 5},
		{"10 - 20", -10},
		{"(10 + 10) / (1 * 2) + 5", 15},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEval_BooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"3 > 2", true},
		{"3 == 2", false},
		{"3 != 2", true},
		{"(1 + 1) == 2", true},
		{"true == true", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEval_BangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEval_ConditionalExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { return 10; }", 10},
		{"if (true) { 10; }", 10},
		{"if (1) { return 10; }", nil},
		{"if (false) { return 10; }", nil},
		{"if (1 == 1) { return 10; }", 10},
		{"if (false) { return 10; } else { return 20; }", 20},
		{"if (true) { return 30; } else { return 20; }", 30},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEval_ReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
	if (true) {
		if(true) {
			return 10;
		}
		
		return 5;
	}
`, 10},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		evals, ok := evaluated.(*object.Integer)

		if !ok {
			t.Errorf("No return object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if evals.Value != tc.expected {
			t.Errorf("evaluated.Value does not equal %v. Got=%d\n", tc.expected, evals.Value)
		}
	}
}

func TestEval_ErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true - 18;", "TYPE-MISMATCH: BOOLEAN - INTEGER"},
		{"18 + false; 12;", "TYPE-MISMATCH: INTEGER + BOOLEAN"},
		{"-false;", "UNKNOWN-OPERATOR: -BOOLEAN"},
		{"false + false;", "UNKNOWN-OPERATOR: BOOLEAN + BOOLEAN"},
		{"if(true) { true + false; }", "UNKNOWN-OPERATOR: BOOLEAN + BOOLEAN"},
		{"test", "UNKNOWN-IDENTIFIER: test"},
		{`"Test" - "Test"`, "UNKNOWN-OPERATOR: STRING - STRING"},
		{`while(true - 2) { return ""; }`, "TYPE-MISMATCH: BOOLEAN - INTEGER"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("No error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tc.expected {
			t.Errorf("wrong error returned. expected=%q, got=%q", tc.expected, evaluated.(*object.Error).Message)
		}
	}
}

func TestEval_Variables(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var test = 10; test;", 10},
		{"var test = 10 * 2; test;", 20},
		{"var test = 10 * 2; var testb = test; testb", 20},
		{"var test = 10 * 2; var testb = test; var testc = testb * 2 + 10; testc", 50},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestEval_Functions(t *testing.T) {
	input := "func(x) { x * 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("Object is not a function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Function has wrong amount of parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x * 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Function body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestEval_FunctionExecution(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var triple = func(x) { x * 3 }; triple(1);", 3},
	}

	for _, tc := range tests {
		testIntegerObject(t, testEval(tc.input), tc.expected)
	}
}

func TestEval_Strings(t *testing.T) {
	input := `"Testing two"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("Object is not string. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Testing two" {
		t.Fatalf("String has incorrect value expected=%s. got=%s", "Testing two", str.Value)
	}
}

func TestEval_StringConcatenation(t *testing.T) {
	input := `"Testing" + " " + "two"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("Object is not string. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Testing two" {
		t.Fatalf("String has incorrect value expected=%s. got=%s", "Testing two", str.Value)
	}
}

func TestEval_WhileLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`var executed = 0; while(executed < 5) { var executed = executed + 1; }; executed`, 5},
		{`var executed = 0; while(executed < 20) { var executed = executed + 1; }; executed`, 20},
		{`var execute = true; while(execute) { var execute = false; }; execute`, false},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		integer, ok := evaluated.(*object.Integer)

		if ok {
			if integer.Value != int64(tc.expected.(int)) {
				t.Fatalf("While loop should have executed %d times. got=%d", tc.expected, evaluated.(*object.Integer).Value)
			}
		} else {
			boolean, ok := evaluated.(*object.Boolean)

			if !ok {
				t.Fatalf("While loop returned incorrect type. expected=%+v. got=%+v", object.BOOLEAN, evaluated.Type())
			} else {
				if boolean.Value != tc.expected {
					t.Fatalf("While loop should have returned %t. got=%t", tc.expected, boolean.Value)
				}
			}
		}
	}
}

func TestEval_BuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("test")`, 4},
		{`len(1)`, "ARGUMENT INVALID TYPE TO BUILT-IN FUNCTION `len`. got=INTEGER. expected=STRING"},
		{`len("1", "2")`, "WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `len`. expected=1. got=2"},
		{`len()`, "WRONG NUMBER OF ARGUMENTS TO BUILT-IN FUNCTION `len`. expected=1. got=0"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			err, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not object.Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if err.Message != expected {
				t.Errorf("Wrong error received. expected=%q. got=%q", expected, err.Message)
			}
		}
	}
}

func TestEval_Arrays(t *testing.T) {
	input := `[1, 2 + 2, "three"]`

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("evaluated.Type IS NOT object.Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("Array does not contain correct amount of elements. expected=%d. got=%d", 3, len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
}

func TestEval_IndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[20, 1, 30][1]`, 1},
		{`[20, 1, 30][2]`, 30},
		{`[20, 1, 30][0]`, 20},
		{`[3 * 2, 1, 30][0]`, 6},
		{`var test = 0; [1, 2, 3][test]`, 1},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)

		if !ok {
			testNullObject(t, evaluated)
		} else {
			testIntegerObject(t, evaluated, int64(integer))
		}
	}
}

func TestEval_Hashes(t *testing.T) {
	input := `
	var key = "hash_key"
	{
		key: 10,
		2: 20,
		true: 0,
		"test": 10 / 2
	}
`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "hash_key"}).HashKey(): 10,
		(&object.Integer{Value: 2}).HashKey():         20,
		TRUE.HashKey():                                0,
		(&object.String{Value: "test"}).HashKey():     5,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d. expected=%d", len(result.Pairs), len(expected))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestEval_HashIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.Create(input)
	p := parser.Create(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}
