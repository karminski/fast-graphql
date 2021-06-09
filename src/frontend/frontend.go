// frontend.go
// 

package frontend

import (
    "fast-graphql/src/graphql"

    "os"
    "github.com/davecgh/go-spew/spew"

)

// generate query hash, documents, Variables
func Compile(request graphql.Request) (*Document, error) {
	var document *Document
	var err error

    // get query hash
    request.QueryHash = GetQueryHash(request.Query)

	// get ast cache
	if cachedDoc, ok := loadAST(request.QueryHash); ok {
		return &cachedDoc, nil
	}

    // Arguments Scanner
    var ctx  *ContextWithArguments
    ctx, err = ScanArguments(request.Query)
    generateRequestVariables(request, ctx)
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    spewo.Dump(ctx)
    spewo.Dump(request)

    if true {
        os.Exit(1)
    }

	// parse
    lexer := NewLexer(request.Query)
    if document, err = parseDocument(lexer); err != nil {
    	return nil, err
    }

    // document end
    lexer.NextTokenIs(TOKEN_EOF) 
    
    // set ast cache
    saveAST(request.QueryHash, *document) 

    return document, nil
}




