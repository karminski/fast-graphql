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

func parseName(lexer *Lexer) *Name {
    lineNum, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    for _, b := range []rune(token) {
        if (b == '_' || 
            b >= 'a' && b <= 'z' ||
            b >= 'A' && b <= 'Z' ||
            b >= '0' && b <= '9' ){
            continue
        } else {
            err := fmt.Sprintf("line %d: unexpected symbol near '%v', it is not a GraphQL name expression", lineNum, token)
            panic(err)
        }
    }
    return &Name{lineNum, token}
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
    enumValueDefinition.Value      = parseName(lexer)
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
    operationName       = parseOperationName(lexer)
    fmt.Printf("\033[33mOperationName: %v \033[0m\n", operationName.Name)
    variableDefinitions = parseVariableDefinitions(lexer)
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

func parseOperationName(lexer *Lexer) *OperationName {
    fmt.Printf("\033[31m[INTO] func parseOperationName  \033[0m\n")

    name := parseName(lexer)
    return &OperationName{name.LineNum, name}
}

/* 
    Parse VariableDefinitions 
    VariableDefinitions ::= <"("> VariableDefinition+ <")">
    VariableDefinition ::= <Ignored> VariableName <":"> <Ignored> Type <Ignored> DefaultValue? <Ignored>
*/

func parseVariableDefinitions(lexer *Lexer) []*VariableDefinition {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinitions  \033[0m\n")

    var VariableDefinitions []*VariableDefinition
    if lexer.LookAhead() != TOKEN_LEFT_PAREN { // variable definitation should start with "("
        return VariableDefinitions
    }
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    // parse variable def until token is ")"
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        VariableDefinitions = append(VariableDefinitions, parseVariableDefinition(lexer))
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return VariableDefinitions
}

func parseVariableDefinition(lexer *Lexer) *VariableDefinition {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinition  \033[0m\n")

    lineNum := lexer.GetLineNum()
    lexer.NextTokenIs(TOKEN_VAR_PREFIX)
    variableName := parseVariableName(lexer)
    lexer.NextTokenIs(TOKEN_COLON)
    variableType := parseType(lexer)
    var variableDefaultValue *DefaultValue
    if lexer.LookAhead() == TOKEN_EQUAL {
        variableDefaultValue = parseDefaultValue(lexer)
    }
    return &VariableDefinition{lineNum, variableName, variableType, variableDefaultValue}
}

/*
    Parse VariableName
    VariableName ::= <"$"> Name
 */
func parseVariableName(lexer *Lexer) *VariableName {
    fmt.Printf("\033[31m[INTO] func parseVariableName  \033[0m\n")

    name := parseName(lexer)
    return &VariableName{lexer.GetLineNum(), name}
}

/**
 * Parse Type
 * Type ::= TypeName | ListType | NonNullType
 * TypeName ::= Name
 * ListType ::= <"["> Type <"]">
 * NonNullType ::= TypeName <"!"> | ListType <"!">
 */

func parseType(lexer *Lexer) Type {
    fmt.Printf("\033[31m[INTO] func parseType  \033[0m\n")

    var typeRet Type

    // parse type
    switch lexer.LookAhead() {
    case TOKEN_IDENTIFIER:   // named type
        typeRet = parseTypeName(lexer)
    case TOKEN_LEFT_BRACKET: // list type, start with "["
        typeRet = parseListType(lexer)
    }

    // for NonNullType
    if lexer.LookAhead() == TOKEN_NOT_NULL {
        typeRet = parseNonNullType(lexer, typeRet)   
    }
    return typeRet
}

func parseTypeName(lexer *Lexer) *NamedType {
    fmt.Printf("\033[31m[INTO] func parseTypeName  \033[0m\n")

    name := parseName(lexer)
    return &NamedType{name.LineNum, name}
}

func parseListType(lexer *Lexer) *ListType {
    fmt.Printf("\033[31m[INTO] func parseListType  \033[0m\n")

    namedType := parseType(lexer) 
    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET) // and end with "]"
    return &ListType{lexer.GetLineNum(), namedType}
}

func parseNonNullType(lexer *Lexer, previousType Type) *NonNullType {
    fmt.Printf("\033[31m[INTO] func parseNonNullType  \033[0m\n")

    return &NonNullType{lexer.GetLineNum(), previousType}
}

/**
 * Parse DefaultValue
 * DefaultValue ::= <"="> <Ignored> Value
 */

func parseDefaultValue(lexer *Lexer) *DefaultValue {
    fmt.Printf("\033[31m[INTO] func parseDefaultValue  \033[0m\n")

    value := parseValue(lexer)
    return &DefaultValue{lexer.GetLineNum(), value}
}


/*
    Parse Value
    Value ::= VariableName | IntValue | FloatValue | ListValue | StringValue | BooleanValue | EnumValue | ObjectValue
    IntValue ::= #"[\+\-0-9]+"
    FloatValue ::= #"[\+\-0-9]+\.[0-9]"
    ListValue ::= <"["> <"]"> | <"["> OneOrMoreValue <"]">
    OneOrMoreValue ::= [Value <Ignored>]+
    StringValue ::= <"\""><"\""> | <"\""> StringCharacter+ <"\"">
    StringCharacter ::= #"[\x{9}\x{20}\x{21}\x{23}-\x{5B}\x{5D}-\uFFFF]" | "\\" "u" EscapedUnicode | "\\" EscapedCharacter
    BooleanValue ::= "true" | "false"
    EnumValue ::= #"(?!(true|false|null))[_A-Za-z][_0-9A-Za-z]*"
    ObjectValue ::= <"{"> ObjectField <"}">
 */

func parseValue(lexer *Lexer) Value {
    fmt.Printf("\033[31m[INTO] func parseValue  \033[0m\n")

    var value Value
    switch lexer.LookAhead() {
    case TOKEN_VAR_PREFIX: // VariableName, start with "$"
        value = parseVariableValue(lexer)
    case TOKEN_NUMBER:     // number, include IntValue, FloatValue
        value = parseNumberValue(lexer)
    case TOKEN_QUOTE:      // string
        value, _ = parseStringValue(lexer)
    case TOKEN_TRUE:
        value = parseBooleanValue(lexer)
    case TOKEN_FALSE:
        value = parseBooleanValue(lexer)
    case TOKEN_IDENTIFIER:
        return nil
    default:
        return nil
    }
    return value
}

func isFloat(token string) bool {
    i := strings.Index(token, ".")
    if i < 0 {
        return false
    }
    return true
}

func parseVariableValue(lexer *Lexer) Value {
    fmt.Printf("\033[31m[INTO] func parseVariableValue  \033[0m\n")

    lexer.NextTokenIs(TOKEN_VAR_PREFIX) // start with $
    _, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    return VariableValue{lexer.GetLineNum(), token}

}

func parseNumberValue(lexer *Lexer) Value {
    fmt.Printf("\033[31m[INTO] func parseNumberValue  \033[0m\n")

    _, token := lexer.NextTokenIs(TOKEN_NUMBER)
    if isFloat(token) {
        num, _ := strconv.ParseFloat(token, 64)
        return FloatValue{lexer.GetLineNum(), num}
    } else {
        num, _ := strconv.Atoi(token)
        return IntValue{lexer.GetLineNum(), num}
    }
    return nil
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
func parseBooleanValue(lexer *Lexer) BooleanValue {
    fmt.Printf("\033[31m[INTO] func parseBooleanValue  \033[0m\n")

    tokenType := lexer.LookAhead()
    if tokenType == TOKEN_TRUE {
        lexer.NextTokenIs(TOKEN_TRUE)
        return BooleanValue{lexer.GetLineNum(), true}
    }
    lexer.NextTokenIs(TOKEN_FALSE)
    return BooleanValue{lexer.GetLineNum(), false}
}
// 
// func EnumValue(lexer *Lexer) *EnumValue {
//     return nil
// }
// 
// func ObjectValue(lexer *Lexer) *ObjectValue {
//     return nil
// }

/*
    Parse ObjectField
    ObjectField ::= <Ignored> Name <":"> <Ignored> Value <Ignored>
 */

func parseObjectField(lexer *Lexer) *ObjectField {
    return nil
}



/* 
    Parse Directives
    Directives ::= Directive+
    Directive ::= <"@"> Name Arguments? <Ignored>
 */

func parseDirectives(lexer *Lexer) ([]*Directive, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectives  \033[0m\n")

    var directives []*Directive
    if lexer.LookAhead() != TOKEN_AT {
        return directives, nil
    }
    for lexer.LookAhead() == TOKEN_AT { // Directive start with "@"
        directive := parseDirective(lexer)
        directives = append(directives, directive)
    }
    return directives, nil
}

func parseDirective(lexer *Lexer) *Directive {
    fmt.Printf("\033[31m[INTO] func parseDirective  \033[0m\n")

    name      := parseName(lexer)
    arguments := parseArguments(lexer)
    return &Directive{lexer.GetLineNum(), name, arguments}
}

/*
    Parse Arguments
    Arguments ::= <"("> <Ignored> Argument+ <")">
    Argument ::= ArgumentName <":"> <Ignored> ArgumentValue <Ignored>*
    ArgumentName ::= Name 
    ArgumentValue ::= Value | VariableName
 */

func parseArguments(lexer *Lexer) []*Argument {
    fmt.Printf("\033[31m[INTO] func parseArguments  \033[0m\n")

    var arguments []*Argument 
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        arguments = append(arguments, parseArgument(lexer))
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return arguments
}

func parseArgument(lexer *Lexer) *Argument {
    fmt.Printf("\033[31m[INTO] func parseArgument  \033[0m\n")

    argumentName := parseArgumentName(lexer)
    lexer.NextTokenIs(TOKEN_COLON)
    argumentValue := parseArgumentValue(lexer)
    return &Argument{lexer.GetLineNum(), argumentName, argumentValue}
}

func parseArgumentName(lexer *Lexer) *ArgumentName {
    fmt.Printf("\033[31m[INTO] func parseArgumentName  \033[0m\n")

    name := parseName(lexer)
    return &ArgumentName{name.LineNum, name}
}

func parseArgumentValue(lexer *Lexer) *ArgumentValue {
    fmt.Printf("\033[31m[INTO] func parseArgumentValue  \033[0m\n")

    value := parseValue(lexer)
    return &ArgumentValue{lexer.GetLineNum(), value}
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

    var alias *Alias
    var fieldName *FieldName
    var arguments []*Argument
    var directives []*Directive
    var selectionSet *SelectionSet
    var err error
    //  Alias & FieldName
    name := parseName(lexer)
    fmt.Printf("parseField.parseName() %v\n", name)
    lineNum := lexer.GetLineNum()
    nextToken := lexer.LookAhead()
    if nextToken == TOKEN_COLON { // suffix is ":"
        alias = &Alias{lineNum, name}
        fieldName = parseFieldName(lexer)
    } else {
        fieldName = &FieldName{lineNum, name}
    } 
    fmt.Printf("parseField.fieldName{lineNum: %v, name: %v} -> %v\n", fieldName.LineNum, fieldName.Name, fieldName)

    // Arguments
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        arguments = parseArguments(lexer)
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

func parseFieldName(lexer *Lexer) *FieldName {
    fmt.Printf("\033[31m[INTO] func parseFieldName  \033[0m\n")

    name := parseName(lexer)
    return &FieldName{lexer.GetLineNum(), name}
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