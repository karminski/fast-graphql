// frontend.go
// 

package frontend

import (

)

func Compile(query string) *Document {
    lexer := NewLexer(query)
    document   := parseDocument(lexer)
    lexer.NextTokenIs(TOKEN_EOF) // set EOF for ast end
    return document
}

