// query_variables_parser.go
package frontend

import (
)

type QueryVariables map[string]interface{}


// QueryVariables ::= Ignored "{" Ignored QueryVariable+ Ignored "}" Ignored
// QueryVariable  ::= Ignored VariableName Ignored ":" Ignored VariableValue Ignored
// VariableName   ::= StringValue
// VariableValue  ::= IntValue | FloatValue | StringValue | BooleanValue | NullValue 

func ParseQueryVariables(lexer *Lexer) (QueryVariables, error) {
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
	var name  StringValue
	var value interface{}
	var err   error

	// json name "abc"
	if name, err = parseStringValue(lexer); err != nil {
		return name.Value, value, err
	}

	// ":"
    lexer.NextTokenIs(TOKEN_COLON)

    // value
    if value, err = parseStringValue(lexer); err != nil {
    	return name.Value, value, err
    }

    return name.Value, value, nil
}
