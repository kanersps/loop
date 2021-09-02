package tokens

type TokenType int64

type Token struct {
	TokenType TokenType
	Value     string
}

var keywords = map[string]TokenType{
	"var":   VariableDeclaration,
	"print": Print,
}

func FindKeyword(keyword string) TokenType {
	if token, ok := keywords[keyword]; ok {
		return token
	}

	return Identifier
}

const (
	Unknown             TokenType = 0
	Number              TokenType = 1
	Operator            TokenType = 2
	VariableDeclaration TokenType = 3
	Identifier          TokenType = 4
	Equals              TokenType = 5
	Print               TokenType = 6
	SemiColon           TokenType = 7
	LeftParentheses     TokenType = 8
	RightParentheses    TokenType = 9
	Comma               TokenType = 10
	Plus                TokenType = 11
	LeftBrace           TokenType = 12
	RightBrace          TokenType = 13
	EOF                 TokenType = 14
)