// frontend.go
// 

package frontend

import (

)

func Compile(query string) (*Document, error) {
	var document *Document
	var err error
    lexer := NewLexer(query)
    if document, err = parseDocument(lexer); err != nil {
    	return nil, err
    }
    // set EOF for document end
    lexer.NextTokenIs(TOKEN_EOF) 
    return document, nil
}

