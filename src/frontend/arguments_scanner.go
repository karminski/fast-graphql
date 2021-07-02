// arguments_scanner.go

package frontend

import (
    "fast-graphql/src/graphql"
    "bytes"
    "errors"
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
// ContextWithArguments ::= Ignored OperationType Ignored OperationName Ignored TargetArguments+ Ignored
// 
// Type & Name Expression
// OperationType       ::= "query" | "mutation" | "subscription" 
// OperationName       ::= Name 
// 
// TargetArguments Expression 
// TargetArguments      ::= "(" Ignored Argument+ | VariableDefinition+ Ignored ")" Ignored
// TargetArgument       ::= Name Ignored ":" Ignored Value Ignored
// VariableDefinition   ::= Variable Ignored ":" Ignored Type Ignored DefaultValue? Ignored
// Variable             ::= "$" Name
// DefaultValue         ::= "=" Ignored Value


type ContextWithArguments struct {
    LastLineNum               int
    VariableDefinitionsPos    int
    TargetArguments        []*TargetArguments 
}

func (ctx *ContextWithArguments)IsTargetArgumentsAvaliable() bool {
    if len(ctx.TargetArguments) != 0 {
        return true
    }
    return false
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
    var ctx ContextWithArguments
    var err error

    // parse
    lexer := NewLexer(request.Query)
    // skip OperationType & OperationName 
    parseOperationType(lexer)
    parseName(lexer) 

    // pin variable definition position
    ctx.VariableDefinitionsPos = lexer.GetPos()

    // request does not provide query variables, but variable definitions express detected, 
    // may be syntax error, stop arguments scanner, return to normal parse prhase.
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        err = errors.New("ScanArguments(): request does not provide query variables, but variable definitions express detected.") 
        return nil, err
    }
    
    // parse left part
    if err = parseContextWithArguments(lexer, &ctx); err != nil {
        return nil, err
    }

    return &ctx, nil
}


func parseContextWithArguments(lexer *Lexer, ctx *ContextWithArguments) error {
    var err error

    // LastLineNum
    ctx.LastLineNum = lexer.GetLineNum()
    
    // parse TargetArguments
    for {
        nextToken := lexer.LookAhead()
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
            return err
        }

        ctx.TargetArguments = append(ctx.TargetArguments, TargetArguments)
    }   
    return nil
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
    for _, targ := range ctx.TargetArguments {
        for _, arg := range targ.Arguments {
            request.SubstitutedVariables[arg.Name.Value] = arg.Value
        }
    }
    return nil
}


func ArgumentsSubstitution(request *graphql.Request, ctx *ContextWithArguments) {
    var buffer                    bytes.Buffer
    var argumentsBuffer           bytes.Buffer
    var variableDefinitionsBuffer bytes.Buffer
    lastPos := 0
    // concat
    variableDefinitionsBuffer.WriteString("(")
    for _, targs := range ctx.TargetArguments {
        for _, targ := range targs.Arguments {
            if lastPos == 0 {
                argumentsBuffer.WriteString(request.Query[:targ.ValueStartPos])
            } else {
                argumentsBuffer.WriteString(request.Query[lastPos:targ.ValueStartPos])
            }
            argumentsBuffer.WriteString("$")
            argumentsBuffer.WriteString(targ.Name.Value)
            variableDefinitionsBuffer.WriteString("$")
            variableDefinitionsBuffer.WriteString(targ.Name.Value)
            variableDefinitionsBuffer.WriteString(":")
            variableDefinitionsBuffer.WriteString(targ.Name.Value)
            variableDefinitionsBuffer.WriteString(",")
            lastPos = targ.ValueEndPos
        }
    }
    variableDefinitionsBuffer.WriteString(")")
    // write back to query
    argumentsBuffer.WriteString(request.Query[lastPos:])
    request.SubstitutedQuery = argumentsBuffer.String()
    buffer.WriteString(request.SubstitutedQuery[:ctx.VariableDefinitionsPos])
    buffer.WriteString(variableDefinitionsBuffer.String())
    buffer.WriteString(request.SubstitutedQuery[ctx.VariableDefinitionsPos:])
    request.SubstitutedQuery = buffer.String()
    request.GenerateSubstitutedQueryHash()
}


