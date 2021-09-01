package tokens

type TokenType int64

type Token struct {
	TokenType TokenType
	Value string
}

const(
	Unknown TokenType = 0
	Number TokenType = 1
	Operator TokenType = 2
	VariableDeclaration TokenType = 3
	VariableIdentifier TokenType = 4
	Equals TokenType = 5
	Print TokenType = 6
)
