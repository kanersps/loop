package tokenizer

import (
	"fmt"
	"loop/src/parser/tokenizer/tokens"
	"os"
	"strconv"
	"strings"
)

var previousToken = tokens.Token{TokenType: tokens.Unknown, Value: ""}

func Tokenize(value string, line int) []tokens.Token {
	tokens := []tokens.Token{}

	keywords := strings.Split(value, " ")

	for _, keyword := range keywords {
		token := GetToken(keyword, line)
		previousToken = token

		if token.TokenType == 0 {
			fmt.Println(fmt.Sprintf("Error #001: Unknown keyword \"%s\" on line %d", keyword, line))
			os.Exit(1)
		}
	}

	return tokens
}

func GetToken(keyword string, line int) (tokens.Token) {
	returnToken := tokens.Token{}

	if value, err := strconv.Atoi(keyword); err == nil {
		returnToken = tokens.Token{
			TokenType: tokens.Number,
			Value:     string(value),
		}
	} else if keyword == "var" {
		returnToken = tokens.Token{
			TokenType: tokens.VariableDeclaration,
			Value:     "",
		}
	} else if previousToken.TokenType == tokens.VariableDeclaration {
		return tokens.Token{
			TokenType: tokens.VariableIdentifier,
			Value: keyword,
		}
	} else if keyword == "=" {
		return tokens.Token{
			TokenType: tokens.Equals,
			Value:     "",
		}
	} else if keyword == "print" {
		return tokens.Token{
			TokenType: tokens.Print,
			Value:     "",
		}
	} else if strings.HasPrefix(keyword, "\"") {

	}

	return returnToken

}