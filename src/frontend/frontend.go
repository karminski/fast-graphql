// frontend.go
// 

package frontend

import (
    "fast-graphql/src/graphql"

    "github.com/davecgh/go-spew/spew"

)

// generate query hash, documents, Variables
func Compile(request *graphql.Request) (*Document, error) {
	var document *Document
	var err error

    // get query hash
    request.GenerateQueryHash()

	// get ast cache
	if cachedDoc, ok := loadAST(request); ok {
		return &cachedDoc, nil
	}

    // Arguments Scanner
    var ctx  *ContextWithArguments
    ctx, err = ScanArguments(request)
    GenerateRequestVariables(request, ctx)
    ArgumentsSubstitution(request, ctx)
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    spewo.Dump(request)


	// parse
    lexer := NewLexer(request.Query)
    if document, err = parseDocument(lexer); err != nil {
    	return nil, err
    }

    // document end
    lexer.NextTokenIs(TOKEN_EOF) 
    
    // set ast cache
    saveAST(request, *document) 

    return document, nil
}


