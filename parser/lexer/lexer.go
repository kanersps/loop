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
	var returnToken tokens.Token

	l.SkipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadCharacter()
			returnToken = tokens.Token{TokenType: tokens.EqualsInfix, Value: string(ch) + string(l.ch)}
		} else {
			returnToken = tokens.Token{TokenType: tokens.Equals, Value: string(l.ch)}
		}
	case ';':
		returnToken = tokens.Token{TokenType: tokens.SemiColon, Value: string(l.ch)}
	case '(':
		returnToken = tokens.Token{TokenType: tokens.LeftParentheses, Value: string(l.ch)}
	case ')':
		returnToken = tokens.Token{TokenType: tokens.RightParentheses, Value: string(l.ch)}
	case ',':
		returnToken = tokens.Token{TokenType: tokens.Comma, Value: string(l.ch)}
	case '+':
		returnToken = tokens.Token{TokenType: tokens.Plus, Value: string(l.ch)}
	case '{':
		returnToken = tokens.Token{TokenType: tokens.LeftBrace, Value: string(l.ch)}
	case '}':
		returnToken = tokens.Token{TokenType: tokens.RightBrace, Value: string(l.ch)}
	case '[':
		returnToken = tokens.Token{TokenType: tokens.LeftBracket, Value: string(l.ch)}
	case ']':
		returnToken = tokens.Token{TokenType: tokens.RightBracket, Value: string(l.ch)}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadCharacter()
			returnToken = tokens.Token{TokenType: tokens.NotEquals, Value: string(ch) + string(l.ch)}
		} else {
			returnToken = tokens.Token{TokenType: tokens.Bang, Value: string(l.ch)}
		}
	case '*':
		returnToken = tokens.Token{TokenType: tokens.Asterisk, Value: string(l.ch)}
	case '/':
		returnToken = tokens.Token{TokenType: tokens.Slash, Value: string(l.ch)}
	case '<':
		returnToken = tokens.Token{TokenType: tokens.LessThan, Value: string(l.ch)}
	case '>':
		returnToken = tokens.Token{TokenType: tokens.GreaterThan, Value: string(l.ch)}
	case '-':
		returnToken = tokens.Token{TokenType: tokens.Minus, Value: string(l.ch)}
	case '"':
		returnToken = tokens.Token{
			TokenType: tokens.String,
			Value:     l.readString(),
		}
	case 0:
		returnToken = tokens.Token{TokenType: tokens.EOF, Value: string(l.ch)}
	default:
		if isLetter(l.ch) {
			returnToken.Value = l.ReadIdentifier()
			returnToken.TokenType = tokens.FindKeyword(returnToken.Value)

			return returnToken
		} else if isDigit(l.ch) {
			returnToken.TokenType = tokens.Number
			returnToken.Value = l.readNumber()

			return returnToken
		} else {
			returnToken = tokens.Token{TokenType: tokens.Unknown, Value: string(l.ch)}
		}
	}

	l.ReadCharacter()

	return returnToken
}

func (l *Lexer) ReadIdentifier() string {
	position := l.position

	// TODO: Add support for numbers in identifiers
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.ReadCharacter()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
