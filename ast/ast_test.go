package ast

import (
	"github.com/kanersps/loop/parser/tokens"
	"testing"
)

func TestProgram_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: tokens.Token{TokenType: tokens.VariableDeclaration, Value: "var"},
				Name: &Identifier{
					Token: tokens.Token{TokenType: tokens.Identifier, Value: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tokens.Token{TokenType: tokens.Identifier, Value: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "var myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}

}
