package ast

import (
	"bytes"
	"github.com/kanersps/loop/parser/tokens"
)

type Node interface {
	TokenValue() string
	String() string
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
	Token      tokens.Token
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

type ReturnStatement struct {
	Token       tokens.Token
	ReturnValue Expression
}

func (vs *ReturnStatement) statementNode()     {}
func (vs *ReturnStatement) TokenValue() string { return vs.Token.Value }

type Identifier struct {
	Token tokens.Token
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) TokenValue() string { return il.Token.Value }
func (il *IntegerLiteral) String() string     { return il.Token.Value }

func (v *VariableStatement) String() string {
	var out bytes.Buffer
	out.WriteString(v.TokenValue() + " ")
	out.WriteString(v.Name.String())
	out.WriteString(" = ")
	if v.Value != nil {
		out.WriteString(v.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenValue() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) String() string { return i.Value }

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) TokenValue() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
