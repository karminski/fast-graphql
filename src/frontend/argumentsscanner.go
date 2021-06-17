// argumentsscanner.go

package frontend

import (
    "fast-graphql/src/graphql"
    "errors"
    "fmt"
    "bytes"

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
// TargetArgument       ::= Name Ignored ":" Ignored Value Ignored
// VariableDefinition   ::= Variable Ignored ":" Ignored Type Ignored DefaultValue? Ignored
// Variable             ::= "$" Name
// DefaultValue         ::= "=" Ignored Value


type ContextWithArguments struct {
    LastLineNum       int
    TargetArguments   []*TargetArguments 
}

type TargetArguments struct {
    Arguments         []*TargetArgument
}

type TargetArgument struct {
    LineNum        int
    Name          *Name
    Value          Value
    ValueStartPos  int // for Argument Substitution
    ValueEndPos    int // -
}


func ScanArguments(request *graphql.Request) (*ContextWithArguments, error) {
    var ctx  *ContextWithArguments
    var err   error

    // parse
    lexer := NewLexer(request.Query)
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

    var targetArgument *TargetArgument
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        if targetArgument, err = parseTargetArgument(lexer); err != nil {
            return nil, err
        }
        TargetArguments.Arguments = append(TargetArguments.Arguments, targetArgument)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return &TargetArguments, nil   
}

func parseTargetArgument(lexer *Lexer) (*TargetArgument, error) {
    var targetArgument TargetArgument
    var err            error

    // LineNum
    targetArgument.LineNum = lexer.GetLineNum()
    // Name
    if targetArgument.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    targetArgument.ValueStartPos = lexer.GetPos()
    // Value
    if targetArgument.Value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    targetArgument.ValueEndPos = lexer.GetPos()
    return &targetArgument, nil
}


func GenerateRequestVariables(request *graphql.Request, ctx *ContextWithArguments) error {
    request.InitSubstitutedVariables()
    for _, ta := range ctx.TargetArguments {
        for _, a := range ta.Arguments {
            if val, ok := a.Value.(IntValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(FloatValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(StringValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(BooleanValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(NullValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(EnumValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(ListValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else if val, ok := a.Value.(ObjectValue); ok {
                request.SubstitutedVariables[a.Name.Value] = val.Value
            } else {
                err := "generateRequestVariables(): ContextWithArguments.TargetArguments.Arguments type assert failed, please check your GraphQL Arguments syntax."
                return errors.New(err)
            }
        }
    }
    return nil
}


func ArgumentsSubstitution(request *graphql.Request, ctx *ContextWithArguments) {
    var buffer bytes.Buffer
    lastPos := 0
    for _, tas := range ctx.TargetArguments {
        for _, ta := range tas.Arguments {
            if lastPos == 0 {
                buffer.WriteString(request.Query[:ta.ValueStartPos])
            } else {
                buffer.WriteString(request.Query[lastPos:ta.ValueStartPos])
            }
            buffer.WriteString("$")
            buffer.WriteString(ta.Name.Value)
            lastPos = ta.ValueEndPos
        }
    }
    // write back to query
    buffer.WriteString(request.Query[lastPos:])
    request.SubstituteQuery = buffer.String()
    request.GenerateSubstituteQueryHash()
}

