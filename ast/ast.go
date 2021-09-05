package ast

import (
	"bytes"
	"github.com/kanersps/loop/parser/tokens"
	"strings"
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

type InfixExpression struct {
	Token    tokens.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()    {}
func (oe *InfixExpression) TokenValue() string { return oe.Token.Value }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

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

type Boolean struct {
	Token tokens.Token
	Value bool
}

func (b *Boolean) expressionNode()    {}
func (b *Boolean) TokenValue() string { return b.Token.Value }
func (b *Boolean) String() string     { return b.Token.Value }

type IfExpression struct {
	Token       tokens.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()    {}
func (ie *IfExpression) TokenValue() string { return ie.Token.Value }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()     {}
func (bs *BlockStatement) TokenValue() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      tokens.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()    {}
func (fl *FunctionLiteral) TokenValue() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenValue())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     tokens.Token // The '(' token
	Function  Expression   // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()    {}
func (ce *CallExpression) TokenValue() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
