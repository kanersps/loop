package lexer

import (
	"github.com/kanersps/loop/parser/tokens"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func Create(value string) *Lexer {
	lexer := &Lexer{input: value}
	lexer.ReadCharacter()
	return lexer
}

func (l *Lexer) ReadCharacter() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) FindToken() tokens.Token {
	var token tokens.TokenType
	var returnToken tokens.Token

	l.SkipWhitespace()

	switch l.ch {
	case '=':
		token = tokens.Equals
	case ';':
		token = tokens.SemiColon
	case '(':
		token = tokens.LeftParentheses
	case ')':
		token = tokens.RightParentheses
	case ',':
		token = tokens.Comma
	case '+':
		token = tokens.Plus
	case '{':
		token = tokens.LeftBrace
	case '}':
		token = tokens.RightBrace
	case '!':
		token = tokens.Bang
	case '*':
		token = tokens.Asterisk
	case '/':
		token = tokens.Slash
	case '<':
		token = tokens.LessThan
	case '>':
		token = tokens.GreaterThan
	case '-':
		token = tokens.Minus
	case 0:
		token = tokens.EOF
	default:
		if isLetter(l.ch) {
			returnToken.Value = l.ReadIdentifier() // Skips to (
			returnToken.TokenType = tokens.FindKeyword(returnToken.Value)

			return returnToken
		} else if isDigit(l.ch) {
			returnToken.TokenType = tokens.Number
			returnToken.Value = l.readNumber()

			return returnToken
		} else {
			token = tokens.Unknown
		}
	}

	returnToken.TokenType = token

	l.ReadCharacter() // Skips again

	return returnToken
}

func (l *Lexer) ReadIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.ReadCharacter()
	}

	return l.input[position:l.position]
}

func (l *Lexer) SkipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadCharacter()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.ReadCharacter()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
