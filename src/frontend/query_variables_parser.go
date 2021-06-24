// query_variables_parser.go
package frontend

import (
	"fmt"
)

type QueryVariables map[string]interface{}


// QueryVariables ::= Ignored "{" Ignored QueryVariable+ Ignored "}" Ignored
// QueryVariable  ::= Ignored VariableName Ignored ":" Ignored VariableValue Ignored
// VariableName   ::= StringValue
// VariableValue  ::= IntValue | FloatValue | StringValue | BooleanValue | NullValue 

func parseQueryVariables(lexer *Lexer) (QueryVariables, error) {
    fmt.Printf("parseQueryVariables() -> lexer.document: %s\n", lexer.document)

	queryVariables := make(QueryVariables)
	// start with "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)

    // loop until "}"
    var name string
    var value interface{}
    var err error
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
    	if name, value, err = parseQueryVariable(lexer); err != nil {
    		return nil, err
    	}
    	queryVariables[name] = value
    }

    // end with "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return queryVariables, nil
}

func parseQueryVariable(lexer *Lexer) (string, interface{}, error) {
	var name  string
	var value interface{}
	var err   error

	// json field name "xxx"
	if name, err = parseStringValueSimple(lexer); err != nil {
		return name, value, err
	}

	// ":"
    lexer.NextTokenIs(TOKEN_COLON)

    // json value
    if value, err = parseValue(lexer); err != nil {
    	return name, value, err
    }

    return name, value, nil
}
