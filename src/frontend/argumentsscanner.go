// argumentsscanner.go

package frontend

import (
    "fast-graphql/src/graphql"

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
    TargetArguments   []*TargetArguments 
}

type TargetArguments struct {
    Arguments         []*Argument
}


func ScanArguments(query string) (*ContextWithArguments, error) {
    var ctx  *ContextWithArguments
    var err   error

    // parse
    lexer := NewLexer(query)
    if ctx, err = parseContextWithArguments(lexer); err != nil {
        return nil, err
    }

    return ctx, nil
}


func parseContextWithArguments(lexer *Lexer) (*ContextWithArguments, error) {
    var ctx ContextWithArguments
    var err error

    // LastLineNum
    ctx.LastLineNum = lexer.GetLineNum()

    // parse TargetArguments
    
    for {
        nextToken := lexer.LookAhead()
        fmt.Printf("n: %d\n", nextToken)
        // end
        if isDocumentEnd(nextToken) {
            break
        }
        // skip content
        if nextToken != TOKEN_LEFT_PAREN {
            lexer.NextTokenIs(nextToken)
            continue
        }
        // got "(", start parse target arguments
        var TargetArguments *TargetArguments
        if TargetArguments, err = parseTargetArguments(lexer); err != nil {
            return nil, err
        }

        ctx.TargetArguments = append(ctx.TargetArguments, TargetArguments)
    }   
    return &ctx, nil
}


func parseTargetArguments(lexer *Lexer) (*TargetArguments, error) {
    var TargetArguments TargetArguments
    var err             error

    // "("
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)

    // start with '$', it is VariableDefinition, skip
    if lexer.LookAhead() == TOKEN_VAR_PREFIX {
        return nil, nil
    }

    var argument *Argument
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        if argument, err = parseArgument(lexer); err != nil {
            return nil, err
        }
        TargetArguments.Arguments = append(TargetArguments.Arguments, argument)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return &TargetArguments, nil   
}


func generateRequestVariables(request graphql.Request, ctx *ContextWithArguments) error {
    for _, ta := range ctx.TargetArguments {
        for _, a := range ta.Arguments {
            request.Variables[a.Name.Value] = a.Value
        }
    }
    return nil
}


//func ArgumentsSubstitution(request Request) {
//
//}