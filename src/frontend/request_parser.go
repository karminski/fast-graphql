// request_parser.go
package frontend

import (
    "fast-graphql/src/graphql"
)

const REQUEST_FIELD_QUERY          = "query"
const REQUEST_FIELD_VARIABLES      = "variables"
const REQUEST_FIELD_OPERATION_NAME = "operationName"

// Request Sample & EBNF see (./DOCUMENTS/Request-Parser.md)

func parseRequest(lexer *Lexer, request *graphql.Request) (error) {
    
    // start with "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)

    // parse request field 
    if err := parseRequestField(lexer, request); err != nil {
        return err
    }

    // end with "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return nil

}

func parseRequestField(lexer *Lexer, request *graphql.Request) (error) {
    var name string  
    var err  error
    // end with "}"
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        // parse field "xxx"
        if name, err = parseStringValueSimple(lexer); err != nil {
            return err
        }
        // ":"
        lexer.NextTokenIs(TOKEN_COLON)
        // query
        if name == REQUEST_FIELD_QUERY {
            if request.Query, err = parseStringValueSimple(lexer); err != nil {
                return err
            }
            continue
        } else if name == REQUEST_FIELD_OPERATION_NAME {
            if request.OperationName, err = parseStringValueSimple(lexer); err != nil {
                return err
            }
            continue
        } else if name == REQUEST_FIELD_VARIABLES {
            if request.Variables, err = parseQueryVariables(lexer); err != nil {
                return err
            }
            continue
        }
        // mismatch
        return err
    }
    return nil
}

