package ast

import (
	"github.com/kanersps/loop/parser/tokens"
)

type Node interface {
	TokenValue() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionStatement struct {
	Token      tokens.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) TokenValue() string { return es.Token.Value }

type Program struct {
	Statements []Statement
}

type VariableStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (vs *VariableStatement) statementNode()     {}
func (vs *VariableStatement) TokenValue() string { return vs.Token.Value }

type Identifier struct {
	Token tokens.Token // the token.IDENT token Value string
	Value string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) statementNode()     {}
func (i *Identifier) TokenValue() string { return i.Token.Value }

func (p *Program) TokenValue() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenValue()
	} else {
		return ""
	}
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) TokenValue() string { return il.Token.Value }
func (il *IntegerLiteral) String() string     { return il.Token.Value }
