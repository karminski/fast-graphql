// argumentsscanner.go

package backend

import (
    "fmt"
)

// arguments substitution
// 
// ScanArguments()
//  ↓
//  haveArguments()
//  ↓(Y)                        ↓(N)
//  ArgumentsSubstitution()     |
//  ↓                           |
//  loadASTFromCache()          ←
//  
type VariableScanner struct {

}



// ContextWithArguments Expression 
// ContextWithArguments ::= Ignored TargetArguments+ Ignored
// 
// TargetArguments Expression 
// TargetArguments      ::= "(" Ignored Argument+ | VariableDefinition+ Ignored ")" Ignored
// Argument             ::= Name Ignored ":" Ignored Value Ignored
// VariableDefinition   ::= Variable Ignored ":" Ignored Type Ignored DefaultValue? Ignored
// Variable             ::= "$" Name
// DefaultValue         ::= "=" Ignored Value


type ContextWithArguments struct {
    LastLineNum       int
    TargetArguments   []TargetArguments 
}

type TargetArguments struct {
    Arguments         []Argument
}


func ScanArguments(request Request) {
    var ctx  *ContextWithArguments
    var err   error
    var query string

    query = request.Query

    // parse
    lexer := NewLexer(query)
    if ctx, err = parseContextWithArguments(lexer); err != nil {
        return nil, err
    }
}


func parseContextWithArguments(lexer *Lexer) (*ContextWithArguments, error) {
    var ctx ContextWithArguments
    var err error

    // LastLineNum
    ctx.LastLineNum = lexer.GetLineNum()

    // parse TargetArguments
    
    for {
        nextToken := lexer.LookAhead()
        if isDocumentEnd(nextToken) {
            break
        }

        var TargetArguments TargetArgument
        var err            error
        if nextToken == TOKEN_LEFT_PAREN {
            if TargetArguments, err = parseTargetArguments(lexer); err != nil {
                return nil, err
            }
        }

        ctx.TargetArguments = append(ctx.TargetArguments, TargetArguments)
    }   
    return &ctx, nil
}


func parseTargetArguments(lexer *Lexer) (TargetArguments, error) {
    var TargetArguments TargetArguments
    var err             error

    // "("
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)

    // start with '$', it is VariableDefinition, skip
    if lexer.LookAhead() == TOKEN_VAR_PREFIX {
        return nil, nil
    }

    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        if argument, err = parseArgument(lexer); err != nil {
            return nil, err
        }
        TargetArguments.Arguments = append(TargetArguments.Arguments, argument)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return TargetArguments, nil   
}




//func ArgumentsSubstitution(request Request) {
//
//}