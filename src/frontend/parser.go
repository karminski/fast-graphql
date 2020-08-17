// parser.go


package frontend

import (
	"regexp"
	"strings"
	"strconv"
)


/* parse document phrase */

func parseDocument(lexer *Lexer) *Block {
	return &Document{
		LastLineNum: 	   lexer.GetLineNum(),
		Definitions: 	   parseDefinitions(lexer),
		ReturnExpressions: parseReturnExpressions(lexer),
	}
}

func isDocumentEnd(tokenType int) bool {
	if tokenType == TOKEN_EOF {
		return true
	}
	return false
}

/* parse Definition phrase*/

func parseDefinitions(lexer *Lexer) []Definition {
	definitions := make([]Definition, 0, 8)
	for !isDocumentEnd(lexer.LookAhead()) {
		definition := parseDefinition(lexer)
		if _, ok := definition.(*EmptyDefinition); !ok {
			definitions = append(definitions, definition)
		}
	}	
	return definitions
}


func parseDefinition(lexer *Lexer) Definition {
	// parse OperationType and OperationName
	switch lexer.LookAhead() {
	case TOKEN_QUERY:
		return parseOperationTypeDefinition(lexer, TOKEN_QUERY)
	case TOKEN_MUTATION:
		return parseOperationTypeDefinition(lexer, TOKEN_MUTATION)
	default:
		return parseQueryShorthandDefinition(lexer)
	}
	// parse VariableDefinitions
	if lexer.LookAhead() == TOKEN_LEFT_PAREN {

	}
	
	// parse Directives
	// @todo:
	
	// parse SelectionSet
	// SelectionSet ::= <"{"> <Ignored> Selection+ <"}"> <Ignored>
	// Selection ::= Field <Ignored> | FragmentSpread <Ignored> | InlineFragment <Ignored>
	if lexer.LookAhead() == TOKEN_LEFT_BRACE { // start with "\{"
		lexer.NextTokenIs(TOKEN_LEFT_BRACE) // then skip "{" 

	}
}



// hold OperationType, include Query, Mutation and operation name
func parseOperationTypeDefinition(lexer *Lexer, operation int) *OperationTypeDefinition {
	lineNum, _ := lexer.NextTokenIs(operation) 
	// anonymous operation by default, check if it is named query
	token := "" 
	if lexer.LookAhead() == TOKEN_IDENTIFIER {
		lineNum, token = lexer.NextTokenIs(TOKEN_IDENTIFIER)
	}
	return &OperationTypeDefinition{lineNum, token, operation}
}

func parseQueryShorthandDefinition(lexer *Lexer) *OperationTypeDefinition {
	return &OperationTypeDefinition{lexer.GetLineNum(), "", TOKEN_QUERY}
}



/* parse expression phrase */

