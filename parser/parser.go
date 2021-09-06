package parser

import (
	"bufio"
	"fmt"
	"github.com/kanersps/loop/ast"
	"github.com/kanersps/loop/parser/lexer"
	"github.com/kanersps/loop/parser/tokens"
	"os"
	"strconv"
)

var variables = map[string]int{}
var precedences = map[tokens.TokenType]int{
	tokens.EqualsInfix:     EQUALS,
	tokens.NotEquals:       EQUALS,
	tokens.LessThan:        LESSGREATER,
	tokens.GreaterThan:     LESSGREATER,
	tokens.Plus:            SUM,
	tokens.Minus:           SUM,
	tokens.Slash:           PRODUCT,
	tokens.Asterisk:        PRODUCT,
	tokens.LeftParentheses: CALL,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.TokenType]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.TokenType]; ok {
		return p
	}

	return LOWEST
}

func parse(value string, lineNumber int) *lexer.Lexer {
	l := lexer.Create(value)
	return l
}

func ParseValue(value string) *lexer.Lexer {
	return parse(value, 1)
}

func ParseFile(path string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) != 0 {
			parse(scanner.Text(), lineNumber)

			lineNumber++
		}
	}
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func Create(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Prefix parsers
	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.registerPrefix(tokens.Identifier, p.parseIdentifier)
	p.registerPrefix(tokens.Bang, p.parsePrefixExpression)
	p.registerPrefix(tokens.Minus, p.parsePrefixExpression)
	p.registerPrefix(tokens.Number, p.parseIntegerLiteral)
	p.registerPrefix(tokens.True, p.parseBoolean)
	p.registerPrefix(tokens.False, p.parseBoolean)
	p.registerPrefix(tokens.LeftParentheses, p.parseGroupedExpression)
	p.registerPrefix(tokens.If, p.parseIfExpression)
	p.registerPrefix(tokens.Function, p.parseFunctionLiteral)
	p.registerPrefix(tokens.String, p.parseStringLiteral)
	p.registerPrefix(tokens.While, p.parseWhileLiteral)

	// Infix parsers
	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)
	p.registerInfix(tokens.Plus, p.parseInfixExpression)
	p.registerInfix(tokens.Minus, p.parseInfixExpression)
	p.registerInfix(tokens.Slash, p.parseInfixExpression)
	p.registerInfix(tokens.Asterisk, p.parseInfixExpression)
	p.registerInfix(tokens.EqualsInfix, p.parseInfixExpression)
	p.registerInfix(tokens.NotEquals, p.parseInfixExpression)
	p.registerInfix(tokens.LessThan, p.parseInfixExpression)
	p.registerInfix(tokens.GreaterThan, p.parseInfixExpression)
	p.registerInfix(tokens.LeftParentheses, p.parseCallExpression)

	p.ExtractToken()
	p.ExtractToken()

	return p
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(tokens.RightParentheses) {
		p.ExtractToken()
		return args
	}
	p.ExtractToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(tokens.Comma) {
		p.ExtractToken()
		p.ExtractToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(tokens.RightParentheses) {
		return nil
	}
	return args
}

func (p *Parser) parseWhileLiteral() ast.Expression {
	while := &ast.WhileLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(tokens.LeftParentheses) {
		return nil
	}
	p.ExtractToken()
	while.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(tokens.RightParentheses) {
		return nil
	}
	if !p.expectPeek(tokens.LeftBrace) {
		return nil
	}

	while.Body = p.parseBlockStatement()

	return while
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(tokens.LeftParentheses) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(tokens.LeftBrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(tokens.RightParentheses) {
		p.ExtractToken()
		return identifiers
	}
	p.ExtractToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(tokens.Comma) {
		p.ExtractToken()
		p.ExtractToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(tokens.RightParentheses) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(tokens.LeftParentheses) {
		return nil
	}
	p.ExtractToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(tokens.RightParentheses) {
		return nil
	}
	if !p.expectPeek(tokens.LeftBrace) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(tokens.Else) {
		p.ExtractToken()

		if !p.expectPeek(tokens.LeftBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.ExtractToken()
	for !p.curTokenIs(tokens.RightBrace) && !p.curTokenIs(tokens.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.ExtractToken()
	}
	return block
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(tokens.True),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
	}
	p.ExtractToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.ExtractToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) FindError(t tokens.TokenType) {
	//expected := reflect.ValueOf(&book).Elem()

	msg := fmt.Sprintf("Expected %v, got %v instead",
		t, p.peekToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ExtractToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.FindToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.TokenType != tokens.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.ExtractToken()
	}
	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.TokenType]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.TokenType)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(tokens.SemiColon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.TokenType]

		if infix == nil {
			return leftExp
		}

		p.ExtractToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case tokens.VariableDeclaration:
		return p.parseVarStatement()
	case tokens.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(tokens.SemiColon) {
		p.ExtractToken()
	}
	return stmt
}

func (p *Parser) parseVarStatement() ast.Statement {
	stmt := &ast.VariableStatement{Token: p.curToken}
	if !p.expectPeek(tokens.Identifier) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
	if !p.expectPeek(tokens.Equals) {
		return nil
	}

	p.ExtractToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SemiColon) {
		p.ExtractToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.ExtractToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SemiColon) {
		p.ExtractToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %v found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) curTokenIs(t tokens.TokenType) bool {
	return p.curToken.TokenType == t
}
func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) expectPeek(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.ExtractToken()
		return true
	} else {
		p.FindError(t)
		return false
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.ExtractToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(tokens.RightParentheses) {
		return nil
	}
	return exp
}

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
