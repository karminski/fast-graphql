// frontend.go
// 

package frontend

import (
    "fast-graphql/src/graphql"

    "github.com/davecgh/go-spew/spew"
    "fmt"

)

// generate query hash, documents, Variables
func Compile(requestStr string, request *graphql.Request) (*Document, error) {
	var document *Document
	var err       error

    // init request
    if err = initRequestObject(requestStr, request); err != nil {
        return nil, err
    }

    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    fmt.Printf("------------------------------\n")
    spewo.Dump(requestStr)
    spewo.Dump(request)

	// check if AST already cached
	if cachedDoc, ok := loadAST(request); ok {
		return &cachedDoc, nil
	}

    // Arguments Scanner
    var ctx  *ContextWithArguments
    ctx, err = ScanArguments(request)
    GenerateRequestVariables(request, ctx)
    ArgumentsSubstitution(request, ctx)
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

// init request object, parse request string to request object, see (./DOCUMENTS/Request-Parser.md)
func initRequestObject(requestStr string, request *graphql.Request) (error) {
    var err error
    // parse request
    lexer := NewLexer(requestStr)
    if err = parseRequest(lexer, request); err != nil {
        return err
    }
    // get query hash
    request.GenerateQueryHash()
    return nil
}
