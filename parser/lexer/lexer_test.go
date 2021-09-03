package lexer

import (
	"github.com/kanersps/loop/parser/tokens"
	"testing"
)

func TestLexer_FindToken(tester *testing.T) {
	input := `
	var test = 1;
	var testtwo = 14;


	var multiply = func(a, b) {
		return a * b;
	}
`

	tests := []struct {
		expectedType  tokens.TokenType
		expectedValue string
	}{
		// First variable assignment
		{tokens.VariableDeclaration, "var"},
		{tokens.Identifier, "test"},
		{tokens.Equals, ""},
		{tokens.Number, "1"},
		{tokens.SemiColon, ""},

		// Second variable assignment
		{tokens.VariableDeclaration, "var"},
		{tokens.Identifier, "testtwo"},
		{tokens.Equals, ""},
		{tokens.Number, "14"},
		{tokens.SemiColon, ""},

		// First multiply function
		{tokens.VariableDeclaration, "var"},
		{tokens.Identifier, "multiply"},
		{tokens.Equals, ""},
		{tokens.Function, "func"},
		{tokens.LeftParentheses, ""},
		{tokens.Identifier, "a"},
		{tokens.Comma, ""},
		{tokens.Identifier, "b"},
		{tokens.RightParentheses, ""},
		{tokens.LeftBrace, ""},
	}

	l := Create(input)

	for i, test := range tests {
		token := l.FindToken()

		if token.TokenType != test.expectedType {
			tester.Fatalf("test (%d/%d) failed - wrong token: expected=%v, got=%v", i, len(tests), test.expectedType, token.TokenType)
		}

		if token.Value != test.expectedValue {
			tester.Fatalf("test (%d/%d) failed - wrong value: expected=%q, got=%q", i, len(tests), test.expectedValue, token.Value)
		}
	}
}
