module github.com/kanersps/loop

go 1.17

replace github.com/kanersps/loop/parser => ./parser

replace github.com/kanersps/loop/parser/tokenizer => ./parser/tokenizer

replace github.com/kanersps/loop/parser/tokenizer/tokens => ./parser/tokenizer/tokens

require github.com/kanersps/loop/parser v0.0.0-00010101000000-000000000000

require (
	github.com/kanersps/loop/parser/tokenizer v0.0.0-00010101000000-000000000000 // indirect
	github.com/kanersps/loop/parser/tokenizer/tokens v0.0.0-00010101000000-000000000000 // indirect
)
