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
