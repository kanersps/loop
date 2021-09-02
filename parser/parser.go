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

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.registerPrefix(tokens.Identifier, p.parseIdentifier)
	p.registerPrefix(tokens.Number, p.parseIntegerLiteral)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
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
		p.nextToken()
	}
	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.TokenType]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case tokens.VariableDeclaration:
		return p.parseVarStatement()
	default:
		return p.parseExpressionStatement()
	}
}

const (
	_ int = iota
	LOWEST
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(tokens.SemiColon) {
		p.nextToken()
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
	// TODO: We're skipping the expressions until we // encounter a semicolon
	for !p.curTokenIs(tokens.SemiColon) {
		p.nextToken()
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
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
