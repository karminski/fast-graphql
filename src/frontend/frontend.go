// frontend.go
// 

package frontend

import (
	"fmt"
)

func Compile(document string) {
	lexer := NewLexer(document)
	ast   := parse(lexer)
	lexer.NextTokenIs(TOKEN_EOF) // set EOF
	fmt.Printf("%f\n", ast)
}