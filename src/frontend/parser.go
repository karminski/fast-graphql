// parser.go


package frontend

import (
    // "regexp"
    "strings"
    // "strconv"
    "fmt"
    // "os"
    "strconv"
    "errors"

    "github.com/davecgh/go-spew/spew"

)


/**
 * Parse Document
 * Document ::= <Ignored> Definition+ <Ignored>
 */

func parseDocument(lexer *Lexer) (*Document, error) {
    fmt.Println("\n\n\033[33m////////////////////////////////////////// Parser Start ///////////////////////////////////////\033[0m\n")
    fmt.Println("parseDocument\n")
    var definitions []Definition
    var err error
    if definitions, err = parseDefinitions(lexer); err != nil {
        return nil, err
    }
    return &Document{
        LastLineNum:       lexer.GetLineNum(),
        Definitions:       definitions,
        // ReturnExpressions: parseReturnExpressions(lexer),
    }, nil
}

func isDocumentEnd(tokenType int) bool {
    if tokenType == TOKEN_EOF {
        return true
    }
    return false
}

/**
 * Parse Name
 * Name ::= #"[_A-Za-z][_0-9A-Za-z]*"
 */

func parseName(lexer *Lexer) (*Name, error) {
    lineNum, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    for _, b := range []rune(token) {
        if (b == '_' || 
            b >= 'a' && b <= 'z' ||
            b >= 'A' && b <= 'Z' ||
            b >= '0' && b <= '9' ){
            continue
        } else {
            err := fmt.Sprintf("line %d: unexpected symbol near '%v', it is not a GraphQL name expression", lineNum, token)
            return nil, errors.New(err)
        }
    }
    return &Name{lineNum, token}, nil
}

/**
 * Parse Definition
 * Definition ::= TypeSystemDefinition | OperationDefinition | FragmentDefinition
 */
func parseDefinitions(lexer *Lexer) ([]Definition, error) {
    var definitions []Definition
    for !isDocumentEnd(lexer.LookAhead()) {
        var definition Definition
        var err error
        if definition, err = parseDefinition(lexer); err != nil {
            return nil, err
        }
        fmt.Printf("- definitions ------\n%v\n", definition)
        definitions = append(definitions, definition)
    }   
    return definitions, nil
}

func parseDefinition(lexer *Lexer) (Definition, error) {
    // parse OperationType and OperationName
    switch lexer.LookAhead() {
    // type system definitation
    case TOKEN_QUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_DUOQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_TRIQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_HEXQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_IDENTIFIER:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_TYPE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_INTERFACE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_UNION:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_SCHEMA:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_ENUM:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_INPUT:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_DIRECTIVE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_EXTEND:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_SCALAR:
        return parseTypeSystemDefinition(lexer)
    // operation definitation
    case TOKEN_LEFT_BRACE:
        return parseOperationDefinition(lexer)
    case TOKEN_QUERY:
        return parseOperationDefinition(lexer)
    case TOKEN_MUTATION:
        return parseOperationDefinition(lexer)
    case TOKEN_SUBSCRIPTION:
        return parseOperationDefinition(lexer)
    // fragment definitation
    case TOKEN_FRAGMENT:
        return parseFragmentDefinition(lexer)
    default:
        // return parseQueryShorthandDefinition(lexer)
        return nil, nil
    }
}

func parseTypeSystemDefinition(lexer *Lexer) (Definition, error) {
    nextToken := lexer.LookAhead()
    var description StringValue 
    var err error
    // parse description
    if nextToken == TOKEN_QUOTE || nextToken == TOKEN_DUOQUOTE ||  nextToken == TOKEN_TRIQUOTE || nextToken == TOKEN_HEXQUOTE {
        if description, err = parseDescription(lexer); err != nil {
            return nil, err
        }
    }
    // type definition body
    switch lexer.LookAhead() {
    case TOKEN_QUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_DUOQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_TRIQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_HEXQUOTE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_IDENTIFIER:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_TYPE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_INTERFACE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_UNION:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_SCHEMA:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_ENUM:
        var enumTypeDefinition *EnumTypeDefinition
        var err error
        if enumTypeDefinition, err = parseEnumTypeDefinition(lexer); err != nil {
            return nil, err
        }
        enumTypeDefinition.Description = description
        return enumTypeDefinition, err
    case TOKEN_INPUT:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_DIRECTIVE:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_EXTEND:
        return parseTypeSystemDefinition(lexer)
    case TOKEN_SCALAR:
        return parseTypeSystemDefinition(lexer)
    default:
        return nil, nil
    }
    // append description
    
}

/**
 * parseDescription
 * @Reference: http://spec.graphql.org/draft/#StringValue
 * Description:
 *     StringValue
 */
var parseDescription = parseStringValue



// parseTypeSystemDefinition
// parseObjectTypeDefinition
// parseInterfaceTypeDefinition
// parseUnionTypeDefinition
// parseSchemaDefinition

/**
 * EnumTypeDefinition
 * EnumDefinition ::= Description? <"enum"> <Ignored> TypeName <Ignored> Directives? <Ignored> <"{"> EnumValuesDefinition <"}"> <Ignored>
 * EnumValuesDefinition ::= EnumValueDefinition+
 * EnumValueDefinition ::= Description? <Ignored> EnumValue <Ignored> Directives? 
 * EnumValue ::= Name
 */
func parseEnumTypeDefinition(lexer *Lexer) (*EnumTypeDefinition, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func parseEnumTypeDefinition  \033[0m\n")

    var enumTypeDefinition EnumTypeDefinition
    var typeName Name
    var err error
    // start
    lexer.NextTokenIs(TOKEN_ENUM)
    _, typeName.Value = lexer.NextTokenIs(TOKEN_IDENTIFIER)
    enumTypeDefinition.Name = &typeName
    if enumTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
        return nil, err
    }
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // enum fields
    for ;; {
        if lexer.LookAhead() == TOKEN_RIGHT_BRACE {
            break
        }
        var enumValueDefinition *EnumValueDefinition
        var err error
        if enumValueDefinition, err = parseEnumValueDefinition(lexer); err != nil {
            return nil, err
        }
        enumTypeDefinition.Values = append(enumTypeDefinition.Values, enumValueDefinition)
    }

    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return &enumTypeDefinition, nil
}

func parseEnumValueDefinition(lexer *Lexer) (*EnumValueDefinition, error) {
    var enumValueDefinition EnumValueDefinition
    nextToken := lexer.LookAhead()
    var err error
    // parse description
    if nextToken == TOKEN_QUOTE || nextToken == TOKEN_DUOQUOTE ||  nextToken == TOKEN_TRIQUOTE || nextToken == TOKEN_HEXQUOTE {
        enumValueDefinition.Description, err = parseDescription(lexer)
    }
    // EnumValue
    if enumValueDefinition.Value, err = parseName(lexer); err != nil {
        return nil, err
    }
    if enumValueDefinition.Directives, err = parseDirectives(lexer); err != nil {
        return nil, err
    }
    return &enumValueDefinition, nil
}

// parseInputObjectTypeDefinition
// parseDirectiveDefinition
// parseTypeExtensionDefinition
// parseScalarTypeDefinition


/**
 * Parse OperationDefinition & OperationType
 * OperationDefinition ::= <Ignored> OperationType? <Ignored> OperationName? <Ignored> VariableDefinitions? <Ignored> Directives? SelectionSet
 * OperationType ::= Query | Mutation | Subscription  
 */

func parseOperationDefinition(lexer *Lexer) (*OperationDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseOperationDefinition  \033[0m\n")

    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    var lineNum               int
    var operationType         *OperationType
    var operationName         *OperationName
    var variableDefinitions []*VariableDefinition
    var directives          []*Directive
    var selectionSet          *SelectionSet
    var err error

    lineNum = lexer.GetLineNum()

    // anonymous operation
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        goto SHORT_QUERY_OPERATION
    }

    // named operation
    operationType       = parseOperationType(lexer)
    fmt.Printf("\033[33mOperationType: %v \033[0m\n", operationType)
    if operationName, err = parseOperationName(lexer); err != nil {
        return nil, err
    }
    fmt.Printf("\033[33mOperationName: %v \033[0m\n", operationName.Name)
    if variableDefinitions, err = parseVariableDefinitions(lexer); err != nil {
        return nil, err
    }
    fmt.Printf("\033[33mVariableDefinitions: %v \033[0m\n", variableDefinitions)
    spewo.Dump(variableDefinitions)

    if directives, err = parseDirectives(lexer); err != nil {
        return nil, err
    }
    fmt.Printf("\033[33mDirectives: %v \033[0m\n", directives)


    SHORT_QUERY_OPERATION:
        fmt.Printf("\033[33mParse SHORT_QUERY_OPERATION:  \033[0m\n")
        lexer.NextTokenIs(TOKEN_LEFT_BRACE)
        if selectionSet, err = parseSelectionSet(lexer); err != nil {
            return nil, err
        }
        lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
        fmt.Printf("\033[33mSelectionSet: %v \033[0m\n", selectionSet)

    // build OperationDefinition
    return &OperationDefinition{
        lineNum,
        operationType,
        operationName,
        variableDefinitions,
        directives,
        selectionSet,
    }, nil
}


func parseOperationType(lexer *Lexer) *OperationType {
    fmt.Printf("\033[31m[INTO] func parseOperationType  \033[0m\n")

    var operation int
    switch lexer.LookAhead() {
    case TOKEN_QUERY: // operation "query"
        operation = TOKEN_QUERY
    case TOKEN_MUTATION: // operation "mutation"
        operation = TOKEN_MUTATION
    case TOKEN_SUBSCRIPTION: // operation "subscription"
        operation = TOKEN_SUBSCRIPTION
    default:
        operation = TOKEN_QUERY // anonymous operation must be query operation?
    }
    lexer.GetNextToken() // skip after lookAhead
    return &OperationType{lexer.GetLineNum(), tokenNameMap[operation], operation}
}

func parseOperationName(lexer *Lexer) (*OperationName, error) {
    fmt.Printf("\033[31m[INTO] func parseOperationName  \033[0m\n")
    var operationName OperationName
    var err  error
    if operationName.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    return &operationName, nil
}

/* 
    Parse VariableDefinitions 
    VariableDefinitions ::= <"("> VariableDefinition+ <")">
    VariableDefinition ::= <Ignored> VariableName <":"> <Ignored> Type <Ignored> DefaultValue? <Ignored>
*/

func parseVariableDefinitions(lexer *Lexer) ([]*VariableDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinitions  \033[0m\n")

    var VariableDefinitions []*VariableDefinition
    lexer.NextTokenIs(TOKEN_LEFT_PAREN) // start with (
    // parse variable def until token is ")"
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        var variableDefinition *VariableDefinition
        var err error
        if variableDefinition, err = parseVariableDefinition(lexer); err != nil {
            return nil, err
        }
        VariableDefinitions = append(VariableDefinitions, variableDefinition)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return VariableDefinitions, nil
}

func parseVariableDefinition(lexer *Lexer) (*VariableDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinition  \033[0m\n")

    var variableDefinition VariableDefinition
    var err error

    lexer.NextTokenIs(TOKEN_VAR_PREFIX)
    variableDefinition.LineNum = lexer.GetLineNum()
    if variableDefinition.Variable, err = parseName(lexer); err != nil {
        return nil, err
    }
    lexer.NextTokenIs(TOKEN_COLON)
    if variableDefinition.Type, err = parseType(lexer); err != nil {
        return nil, err
    }
    if lexer.LookAhead() == TOKEN_EQUAL {
        if variableDefinition.DefaultValue, err = parseDefaultValue(lexer); err != nil {
            return nil, err
        }
    }
    return &variableDefinition, nil
}


/**
 * Parse Type
 * Type ::= TypeName | ListType | NonNullType
 * TypeName ::= Name
 * ListType ::= <"["> Type <"]">
 * NonNullType ::= TypeName <"!"> | ListType <"!">
 */

func parseType(lexer *Lexer) (Type, error) {
    fmt.Printf("\033[31m[INTO] func parseType  \033[0m\n")

    var typeRet Type
    var err     error

    // parse type
    switch lexer.LookAhead() {
    case TOKEN_IDENTIFIER:   // named type
        if typeRet, err = parseTypeName(lexer); err != nil {
            return nil, err
        }
    case TOKEN_LEFT_BRACKET: // list type, start with "["
        if typeRet, err = parseListType(lexer); err != nil {
            return nil, err
        }
    }

    // for NonNullType at suffix
    if lexer.LookAhead() == TOKEN_NOT_NULL {
        if typeRet, err = parseNonNullType(lexer, typeRet); err != nil {
            return nil, err
        }   
    }
    return typeRet, nil
}

func parseTypeName(lexer *Lexer) (*NamedType, error) {
    fmt.Printf("\033[31m[INTO] func parseTypeName  \033[0m\n")
    var namedType NamedType
    var err       error
    if namedType.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    return &namedType, nil
}

func parseListType(lexer *Lexer) (*ListType, error) {
    fmt.Printf("\033[31m[INTO] func parseListType  \033[0m\n")
    var listType ListType
    var err      error
    lexer.NextTokenIs(TOKEN_LEFT_BRACKET) // start with "]"
    if listType.Type, err = parseType(lexer); err != nil {
        return nil, err
    } 
    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET) // and end with "]"
    return &listType, nil
}

func parseNonNullType(lexer *Lexer, previousType Type) (*NonNullType, error) {
    fmt.Printf("\033[31m[INTO] func parseNonNullType  \033[0m\n")

    return &NonNullType{lexer.GetLineNum(), previousType}, nil
}

/**
 * Parse DefaultValue
 * DefaultValue ::= <"="> <Ignored> Value
 */

func parseDefaultValue(lexer *Lexer) (*DefaultValue, error) {
    fmt.Printf("\033[31m[INTO] func parseDefaultValue  \033[0m\n")
    var value Value 
    var err error
    if value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return &DefaultValue{lexer.GetLineNum(), value}, nil
}


/**
 * Parse Value
 * Value        ::= Variable | IntValue | FloatValue | StringValue | BooleanValue | NullValue | EnumValue | ListValue | ObjectValue
 * IntValue     ::= #"[\+\-0-9]+"
 * FloatValue   ::= #"[\+\-0-9]+\.[0-9]"
 * StringValue  ::= <"\"\"\""><"\"\"\""> | <"\""><"\""> | <"\"">StringCharacter<"\""> | <"\"\"\"">BlockStringCharacter<"\"\"\""> |
 * BooleanValue ::= <"true"> | <"false">
 * NullValue    ::= <"null">
 * EnumValue    ::= #"(?!(true|false|null))[_A-Za-z][_0-9A-Za-z]*" # Name but not "true" or "false" or "null"
 * ListValue    ::= <"["> <"]"> | <"["> Value+ <"]">
 * ObjectValue  ::= <"{"> <"}"> | <"{"> ObjectField+ <"}">
 */

func parseValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseValue  \033[0m\n")

    var value Value
    var err error
    switch lexer.LookAhead() {
    case TOKEN_VAR_PREFIX: // VariableName, start with "$"
        if value, err = parseVariableValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_NUMBER:     // number, include IntValue, FloatValue
        if value, err = parseNumberValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_QUOTE:      // string
        if value, err = parseStringValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_TRUE:
        if value, err = parseBooleanValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_FALSE:
        if value, err = parseBooleanValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_NULL:
        if value, err = parseNullValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_IDENTIFIER:
        if value, err = parseEnumValue(lexer); err != nil {
            return nil, err
        }
    case TOKEN_LEFT_BRACKET:
        if value, err = parseListValue(lexer); err != nil {
            return nil, err
        }
    default:
        return nil, nil
    }
    return value, nil
}

func isFloat(token string) bool {
    i := strings.Index(token, ".")
    if i < 0 {
        return false
    }
    return true
}

func parseVariableValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseVariableValue  \033[0m\n")

    lexer.NextTokenIs(TOKEN_VAR_PREFIX) // start with $
    _, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    return VariableValue{lexer.GetLineNum(), token}, nil

}

func parseNumberValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseNumberValue  \033[0m\n")

    _, token := lexer.NextTokenIs(TOKEN_NUMBER)
    if isFloat(token) {
        num, _ := strconv.ParseFloat(token, 64)
        return FloatValue{lexer.GetLineNum(), num}, nil
    } else {
        num, _ := strconv.Atoi(token)
        return IntValue{lexer.GetLineNum(), num}, nil
    }
    return nil, nil
}


/**
 * parseStringValue
 * @Reference: http://spec.graphql.org/draft/#StringValue
 * StringValue
 *     ""
 *     "StringCharacterlist"
 *     """BlockStringCharacterlistopt"""
 */
func parseStringValue(lexer *Lexer) (StringValue, error) {
    lineNum := lexer.GetLineNum()
    if lexer.LookAhead() == TOKEN_HEXQUOTE {
        lexer.NextTokenIs(TOKEN_HEXQUOTE)
        return StringValue{lineNum, ""}, nil
    }
    if lexer.LookAhead() == TOKEN_DUOQUOTE {
        lexer.NextTokenIs(TOKEN_DUOQUOTE)
        return StringValue{lineNum, ""}, nil
    }
    if lexer.LookAhead() == TOKEN_TRIQUOTE {
        lexer.NextTokenIs(TOKEN_TRIQUOTE)
        str := lexer.scanBeforeToken(tokenNameMap[TOKEN_TRIQUOTE])
        lexer.NextTokenIs(TOKEN_TRIQUOTE)
        return StringValue{lineNum, str}, nil
    }
    if lexer.LookAhead() == TOKEN_QUOTE {
        lexer.NextTokenIs(TOKEN_QUOTE)
        str := lexer.scanBeforeToken(tokenNameMap[TOKEN_QUOTE])
        lexer.NextTokenIs(TOKEN_QUOTE)
        return StringValue{lineNum, str}, nil
    }
    err := "not a StringValue"
    return StringValue{lineNum, ""}, errors.New(err)
}


func parseBooleanValue(lexer *Lexer) (BooleanValue, error) {
    fmt.Printf("\033[31m[INTO] func parseBooleanValue  \033[0m\n")

    tokenType := lexer.LookAhead()
    if tokenType == TOKEN_TRUE {
        lexer.NextTokenIs(TOKEN_TRUE)
        return BooleanValue{lexer.GetLineNum(), true}, nil
    }
    lexer.NextTokenIs(TOKEN_FALSE)
    return BooleanValue{lexer.GetLineNum(), false}, nil
}


func parseNullValue(lexer *Lexer) (NullValue, error) {
    return NullValue{lexer.GetLineNum()}, nil
}


func parseEnumValue(lexer *Lexer) (EnumValue, error) {
    lineNum, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    canNotBe := map[string]bool{tokenNameMap[TOKEN_TRUE]: true, tokenNameMap[TOKEN_FALSE]: true, tokenNameMap[TOKEN_NULL]: true}
    if _, ok := canNotBe[token]; ok {
        err := fmt.Sprintf("line %d: unexpected symbol near '%v', enum value can not be 'true' or 'false' or 'null'.", lineNum, token)
        panic(err)
    }
    for _, b := range []rune(token) {
        if (b == '_' || 
            b >= 'a' && b <= 'z' ||
            b >= 'A' && b <= 'Z' ||
            b >= '0' && b <= '9' ){
            continue
        } else {
            err := fmt.Sprintf("line %d: unexpected symbol near '%v', it is not a GraphQL enum expression", lineNum, token)
            panic(err)
        }
    }
    return EnumValue{lineNum, token}, nil
}

func parseListValue(lexer *Lexer) (ListValue, error) {
    fmt.Printf("\033[31m[INTO] func parseListValue  \033[0m\n")

    var listValue ListValue
    lexer.NextTokenIs(TOKEN_LEFT_BRACKET) // start with [
    for lexer.LookAhead() != TOKEN_RIGHT_BRACKET {
        var value Value 
        var err error
        if value, err = parseValue(lexer); err != nil {
            return listValue, err
        }
        listValue.Value = append(listValue.Value, value)
    }

    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET) // end with ]
    return listValue, nil
}

func parseObjectValue(lexer *Lexer) (ObjectValue, error) {
    fmt.Printf("\033[31m[INTO] func parseObjectValue  \033[0m\n")

    var objectValue ObjectValue
    lexer.NextTokenIs(TOKEN_LEFT_BRACE) // start with {
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        var objectField *ObjectField
        var err          error
        if objectField, err = parseObjectField(lexer); err != nil {
            return objectValue, err
        }
        objectValue.Value = append(objectValue.Value, objectField)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE) // end with }
    return objectValue, nil
}

func parseObjectField(lexer *Lexer) (*ObjectField, error) {
    var objectField ObjectField
    var err error
    if objectField.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    lexer.NextTokenIs(TOKEN_COLON)
    if objectField.Value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return &objectField, nil
}

// 
// func parseFloatValue(lexer *Lexer) *FloatValue {
//     return nil
// }
// 
// func ListValue(lexer *Lexer) *ListValue {
//     return nil
// }
// 
// func OneOrMoreValue(lexer *Lexer) *OneOrMoreValue {
//     return nil
// }
// 



// func parseStringValue(lexer *Lexer) StringValue {
//     fmt.Printf("\033[31m[INTO] func parseStringValue  \033[0m\n")
// 
//     lexer.NextTokenIs(TOKEN_QUOTE)
//     // in quote, all token except TOKEN_QUOTE are string (TOKEN_IDENTIFIER)
//     var strBuf strings.Builder
//     for ;; {
//         tokenType := lexer.LookAhead()
//         if tokenType != TOKEN_QUOTE {
//             _, token := lexer.NextTokenIs(tokenType)
//             strBuf.WriteString(token)
//         } else {
//             break;
//         }
//     }
//     lexer.NextTokenIs(TOKEN_QUOTE)
//     return StringValue{lexer.GetLineNum(), strBuf.String()}
// }

// 
// func StringCharacter(lexer *Lexer) *StringCharacter {
//     return nil
// }
// 

// 
// func EnumValue(lexer *Lexer) *EnumValue {
//     return nil
// }
// 
// func ObjectValue(lexer *Lexer) *ObjectValue {
//     return nil
// }





/* 
    Parse Directives
    Directives ::= Directive+
    Directive ::= <"@"> Name Arguments? <Ignored>
 */

func parseDirectives(lexer *Lexer) ([]*Directive, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectives  \033[0m\n")

    var directives []*Directive
    var directive    *Directive
    var err           error
    if lexer.LookAhead() != TOKEN_AT {
        return directives, nil
    }
    for lexer.LookAhead() == TOKEN_AT { // Directive start with "@"
        if directive, err = parseDirective(lexer); err != nil {
            return nil, err
        }
        directives = append(directives, directive)
    }
    return directives, nil
}

func parseDirective(lexer *Lexer) (*Directive, error) {
    fmt.Printf("\033[31m[INTO] func parseDirective  \033[0m\n")
    var directive Directive
    var err error
    lexer.NextTokenIs(TOKEN_AT)
    directive.LineNum = lexer.GetLineNum()
    if directive.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if directive.Arguments, err = parseArguments(lexer); err != nil {
            return nil, err
        }
    }
    return &directive, nil
}

/*
    Parse Arguments
    Arguments ::= <"("> <Ignored> Argument+ <")">
    Argument ::= ArgumentName <":"> <Ignored> ArgumentValue <Ignored>*
    ArgumentName ::= Name 
    ArgumentValue ::= Value | VariableName
 */

func parseArguments(lexer *Lexer) ([]*Argument, error) {
    fmt.Printf("\033[31m[INTO] func parseArguments  \033[0m\n")

    var arguments []*Argument
    var argument    *Argument
    var err          error 
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        if argument, err = parseArgument(lexer); err != nil {
            return nil, err
        }
        arguments = append(arguments, argument)
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return arguments, nil
}

func parseArgument(lexer *Lexer) (*Argument, error) {
    fmt.Printf("\033[31m[INTO] func parseArgument  \033[0m\n")
    var argument Argument
    var err      error
    argument.LineNum = lexer.GetLineNum()
    if argument.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    lexer.NextTokenIs(TOKEN_COLON)
    if argument.Value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return &argument, nil
}


/* 
    Parse SelectionSet & Selection 
    SelectionSet ::= <"{"> <Ignored> Selection+ <"}"> <Ignored>
    Selection ::= Field <Ignored> | FragmentSpread <Ignored> | InlineFragment <Ignored>
 */

func parseSelectionSet(lexer *Lexer) (*SelectionSet, error) {
    fmt.Printf("\033[31m[INTO] func parseSelectionSet  \033[0m\n")

    var selections []Selection
    lineNum := lexer.GetLineNum() 
    // parse variable def until token is "}"
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        var selectionInterface interface{}
        var err error 
        if selectionInterface, err = parseSelection(lexer); err != nil {
            return nil, err
        }
        selections = append(selections, selectionInterface.(Selection))
    }

    return &SelectionSet{lineNum, selections}, nil
}

func parseSelection(lexer *Lexer) (interface{}, error) {
    fmt.Printf("\033[31m[INTO] func parseSelection  \033[0m\n")

    switch lexer.LookAhead() {
    case TOKEN_DOTS:
        return parseFragment(lexer)
    default:
        return parseField(lexer)
    }
}

/* 
    Parse Field 
    Field ::= Alias? <Ignored> FieldName <Ignored> Arguments? <Ignored> Directives? SelectionSet?
    Alias ::= Name <":">
    FieldName ::= Name
 */

func parseField(lexer *Lexer) (*Field, error) {
    fmt.Printf("\033[31m[INTO] func parseField  \033[0m\n")

    var alias          *Alias
    var fieldName      *FieldName
    var arguments    []*Argument
    var directives   []*Directive
    var selectionSet   *SelectionSet
    var err             error
    //  Alias & FieldName
    var name *Name
    if name, err = parseName(lexer); err != nil {
        return nil ,err
    }
    fmt.Printf("parseField.parseName() %v\n", name)
    lineNum := lexer.GetLineNum()
    nextToken := lexer.LookAhead()
    if nextToken == TOKEN_COLON { // suffix is ":"
        alias = &Alias{lineNum, name}
        if fieldName, err = parseFieldName(lexer); err != nil {
            return nil, err
        }
    } else {
        fieldName = &FieldName{lineNum, name}
    } 
    fmt.Printf("parseField.fieldName{lineNum: %v, name: %v} -> %v\n", fieldName.LineNum, fieldName.Name, fieldName)

    // Arguments
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if arguments, err = parseArguments(lexer); err != nil {
            return nil, err
        }
        fmt.Printf("\033[33marguments: %v \033[0m\n", arguments)
    }

    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        if directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
        fmt.Printf("\033[33mdirectives: %v \033[0m\n", directives)
    }

    // more SelectionSet
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        fmt.Printf("\033[33m into more SelectionSet: \033[0m\n")
        lexer.NextTokenIs(TOKEN_LEFT_BRACE)
        if selectionSet, err = parseSelectionSet(lexer); err != nil {
            return nil, err
        }
        lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
        fmt.Printf("\033[33m out more SelectionSet: \033[0m\n")
    }
    return &Field{lineNum, alias, fieldName, arguments, directives, selectionSet}, nil
}

func parseFieldName(lexer *Lexer) (*FieldName, error) {
    fmt.Printf("\033[31m[INTO] func parseFieldName  \033[0m\n")
    var fieldName FieldName
    var err       error
    fieldName.LineNum = lexer.GetLineNum()
    if fieldName.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    return &fieldName, nil
}


/* 
    Parse FragmentSpread 
    FragmentSpread ::= <"..."> FragmentName <Ignored> Directives?
    FragmentName ::= Name
 */

func parseFragment(lexer *Lexer) (interface{}, error) {
    return nil, nil
}

func parseFragmentSpread(lexer *Lexer) *FragmentSpread {
    return nil
}

func parseFragmentName(lexer *Lexer) *FragmentName {
    return nil
}

/* 
    Parse InlineFragment 
    InlineFragment ::= <"..."> <Ignored> TypeCondition? Directives? SelectionSet?
 */

func parseInlineFragment(lexer *Lexer) *InlineFragment {
    return nil
}


/*
    Parse TypeCondition
    TypeCondition ::= <"on"> <Ignored> TypeName <Ignored>
 */

func parseTypeCondition(lexer *Lexer) *TypeCondition {
    return nil
}


/*
    Parse Ignored
    Ignored ::= Ignore*
    Ignore ::= UnicodeBOM | WhiteSpace | LineTerminator | Comment | Comma
    UnicodeBOM ::= "\uFEFF"
    WhiteSpace ::= #"[\x{9}\x{20}]"   ### ASCII: \t | Space
    LineTerminator ::= #"\x{A}" | #"\x{D}\x{A}" | #"\x{D}"   ### ASCII: \n | \r\n | \r 
    Comment ::= "#" CommentChar* <LineTerminator>
    Comma ::= ","
    CommentChar ::= #"[\x{9}\x{20}-\uFFFF]"
 */


// func parseIgnored(lexer *Lexer) *Ignored {
// 
// }
// 
// func parseIgnore(lexer *Lexer) *Ignore {
// 
// }
// 
// func parseUnicodeBOM(lexer *Lexer) *UnicodeBOM {
// 
// }
// 
// func parseWhiteSpace(lexer *Lexer) *WhiteSpace {
// 
// }
// 
// func parseLineTerminator(lexer *Lexer) *LineTerminator {
// 
// }
// 
// func parseComment(lexer *Lexer) *Comment {
// 
// }
// 
// func parseComma(lexer *Lexer) *Comma {
// 
// }
// 
// func parseCommentChar(lexer *Lexer) *CommentChar {
// 
// }


/**
 * Parse FragmentDefinition
 * FragmentDefinition ::= <"fragment"> <Ignored> FragmentName <Ignored> TypeCondition Directives? SelectionSet
 */
func parseFragmentDefinition(lexer *Lexer) (*FragmentDefinition, error) {
    return nil, nil
}