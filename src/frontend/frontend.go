// frontend.go
// 

package frontend

import (
    "fast-graphql/src/graphql"

    "errors"

)

// generate query hash, documents, Variables
func Compile(requestStr string, request *graphql.Request) (*Document, error) {
    var document *Document
	var err       error

    // init request
    if err = initRequestObject(requestStr, request); err != nil {
        return nil, err
    }

    // run variables scanner & arguments substitution 
    runSubstitution(request)

	// check if AST already cached
	if cachedDoc, ok := loadAST(request); ok {
		return &cachedDoc, nil
	}

	// cache miss, parse     
    if document, err = parseQuery(request); err != nil {
        return nil, err
    }
    
    // set cache
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

// run full parser, parse GraphQL Query
func parseQuery(request *graphql.Request) (*Document, error) {
    var document *Document
    var query     string
    var err       error

    query = request.GetAvaliableQuery()

    // parse
    lexer := NewLexer(query)
    if document, err = parseDocument(lexer); err != nil {
        return nil, err
    }

    // document end
    lexer.NextTokenIs(TOKEN_EOF) 

    return document, nil
}

// if query variables are not set, may be it is in the query, run arguments scanner pick it out.
func runSubstitution(request *graphql.Request) {
    var err error
    if !request.IsVariablesAvaliable() {
        var ctx *ContextWithArguments
        if ctx, err = ScanArguments(request); err != nil {   
            return 
        }
        if ctx.IsTargetArgumentsAvaliable() {
            GenerateRequestVariables(request, ctx)
            ArgumentsSubstitution(request, ctx)
        }
    }

}

func AssertArgumentType(arg interface{}) (interface{}, error) {
    if val, ok := arg.(IntValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(FloatValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(StringValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(BooleanValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(NullValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(EnumValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(ListValue); ok {
        return val.Value, nil
    } else if val, ok := arg.(ObjectValue); ok {
        return val.Value, nil
    } else {
        err := "getFieldArgumentsMap(): Field.Arguments.Argument type assert failed, please check your GraphQL Field.Arguments.Argument syntax."
        return nil, errors.New(err)
    }
}