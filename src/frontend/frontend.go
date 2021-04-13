// frontend.go
// 

package frontend

import (
	"crypto/md5"
)

func Compile(query string) (*Document, error) {
	var document *Document
	var err error

	// get ast cache
	queryHash := md5.Sum([]byte(query))
	if cachedDoc, ok := loadAST(queryHash); ok {
		return &cachedDoc, nil
	}

	// parse
    lexer := NewLexer(query)
    if document, err = parseDocument(lexer); err != nil {
    	return nil, err
    }

    // set EOF for document end
    lexer.NextTokenIs(TOKEN_EOF) 
    
    // set ast cache
    saveAST(queryHash, *document) 

    return document, nil
}




