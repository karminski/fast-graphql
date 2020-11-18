// parser.go


package frontend

import (
    // "regexp"
    // "strings"
    // "strconv"
    "fmt"
)


/**
 * Parse Document
 * Document ::= <Ignored> Definition+ <Ignored>
 */

func parseDocument(lexer *Lexer) *Document {
    fmt.Println("parseDocument\n")
    return &Document{
        LastLineNum:       lexer.GetLineNum(),
        Definitions:       parseDefinitions(lexer),
        // ReturnExpressions: parseReturnExpressions(lexer),
    }
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
func parseDefinitions(lexer *Lexer) []Definition {
    var definitions []Definition
    for !isDocumentEnd(lexer.LookAhead()) {
        definition := parseDefinition(lexer)
        fmt.Printf("- definitions ------\n%v\n", definition)
        definitions = append(definitions, definition)
    }   
    return definitions
}

func parseDefinition(lexer *Lexer) Definition {
    // parse OperationType and OperationName
    switch lexer.LookAhead() {
    // type system definitation
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
        return nil
    }
}

func parseTypeSystemDefinition(lexer *Lexer) *TypeSystemDefinition {
    return nil
}

// parseTypeSystemDefinition
// parseObjectTypeDefinition
// parseInterfaceTypeDefinition
// parseUnionTypeDefinition
// parseSchemaDefinition
// parseEnumTypeDefinition
// parseInputObjectTypeDefinition
// parseDirectiveDefinition
// parseTypeExtensionDefinition
// parseScalarTypeDefinition


/**
 * Parse OperationDefinition & OperationType
 * OperationDefinition ::= <Ignored> OperationType? <Ignored> OperationName? <Ignored> VariableDefinitions? <Ignored> Directives? SelectionSet
 * OperationType ::= Query | Mutation | Subscription  
 */

func parseOperationDefinition(lexer *Lexer) *OperationDefinition {
    fmt.Println("parseOperationDefinition:")
    var lineNum               int
    var OperationType         *OperationType
    var OperationName         *OperationName
    var VariableDefinitions []*VariableDefinition
    var Directives          []*Directive
    var SelectionSet          *SelectionSet

    lineNum = lexer.GetLineNum()

    // anonymous operation
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        goto SHORT_QUERY_OPERATION
    }

    // named operation
    OperationType       = parseOperationType(lexer)
    fmt.Printf("\033[34mOperationType: %v \033[0m\n", OperationType)
    OperationName       = parseOperationName(lexer)
    fmt.Printf("\033[34mOperationName: %v \033[0m\n", OperationName.Name)
    VariableDefinitions = parseVariableDefinitions(lexer)
    fmt.Printf("\033[34mVariableDefinitions: %v \033[0m\n", VariableDefinitions)
    Directives          = parseDirectives(lexer)
    fmt.Printf("\033[34mDirectives: %v \033[0m\n", Directives)

    SHORT_QUERY_OPERATION:
        fmt.Printf("\033[34mParse SHORT_QUERY_OPERATION:  \033[0m\n")
        lexer.NextTokenIs(TOKEN_LEFT_BRACE)
        SelectionSet    = parseSelectionSet(lexer)
        lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
        fmt.Printf("\033[34mSelectionSet: %v \033[0m\n", SelectionSet)

    // build OperationDefinition
    return &OperationDefinition{
        lineNum,
        OperationType,
        OperationName,
        VariableDefinitions,
        Directives,
        SelectionSet,
    }
}


func parseOperationType(lexer *Lexer) *OperationType {
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
    name := parseName(lexer)
    return &OperationName{name.LineNum, name}
}

/* 
    Parse VariableDefinitions 
    VariableDefinitions ::= <"("> VariableDefinition+ <")">
    VariableDefinition ::= <Ignored> VariableName <":"> <Ignored> Type <Ignored> DefaultValue? <Ignored>
*/

func parseVariableDefinitions(lexer *Lexer) []*VariableDefinition {
    var VariableDefinitions []*VariableDefinition
    if lexer.LookAhead() != TOKEN_LEFT_PAREN { // variable definitation should start with "("
        return VariableDefinitions
    }
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    // parse variable def until token is ")"
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        VariableDefinitions = append(VariableDefinitions, parseVariableDefinition(lexer))
    }
    return VariableDefinitions
}

func parseVariableDefinition(lexer *Lexer) *VariableDefinition {
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
    name := parseName(lexer)
    return &NamedType{name.LineNum, name}
}

func parseListType(lexer *Lexer) *ListType {
    namedType := parseType(lexer) 
    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET) // and end with "]"
    return &ListType{lexer.GetLineNum(), namedType}
}

func parseNonNullType(lexer *Lexer, previousType Type) *NonNullType {
    return &NonNullType{lexer.GetLineNum(), previousType}
}

/**
 * Parse DefaultValue
 * DefaultValue ::= <"="> <Ignored> Value
 */

func parseDefaultValue(lexer *Lexer) *DefaultValue {
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

func parseValue(lexer *Lexer) IntValue {
    lexer.NextTokenIs(TOKEN_NUMBER)
    return IntValue{lexer.GetLineNum(), 1}
    // switch lexer.LookAhead() {
    // case TOKEN_VAR_PREFIX: // VariableName, start with "$"
    //     return nil
    //     // return parseVariableName(lexer)
    // case TOKEN_NUMBER: // number, include IntValue, FloatValue
    //     parseIntValue(lexer)
    // case TOKEN_IDENTIFIER:
    // 
    // default:
    //     return nil
    // }

}

// func parseIntValue(lexer *Lexer) *IntValue {
//     return nil
// }
// 
// func FloatValue(lexer *Lexer) *FloatValue {
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
// func StringValue(lexer *Lexer) *StringValue {
//     return nil
// }
// 
// func StringCharacter(lexer *Lexer) *StringCharacter {
//     return nil
// }
// 
// func BooleanValue(lexer *Lexer) *BooleanValue {
//     return nil
// }
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

func parseDirectives(lexer *Lexer) []*Directive {
    var directives []*Directive
    if lexer.LookAhead() != TOKEN_AT {
        return directives
    }
    for lexer.LookAhead() == TOKEN_AT { // Directive start with "@"
        directive := parseDirective(lexer)
        directives = append(directives, directive)
    }
    return directives
}

func parseDirective(lexer *Lexer) *Directive {
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
    var arguments []*Argument 
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        arguments = append(arguments, parseArgument(lexer))
    }
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return arguments
}

func parseArgument(lexer *Lexer) *Argument {
    argumentName := parseArgumentName(lexer)
    lexer.NextTokenIs(TOKEN_COLON)
    argumentValue := parseArgumentValue(lexer)
    return &Argument{lexer.GetLineNum(), argumentName, argumentValue}
}

func parseArgumentName(lexer *Lexer) *ArgumentName {
    name := parseName(lexer)
    return &ArgumentName{name.LineNum, name}
}

func parseArgumentValue(lexer *Lexer) *ArgumentValue {
    value := parseValue(lexer)
    return &ArgumentValue{value.LineNum, value}
}


/* 
    Parse SelectionSet & Selection 
    SelectionSet ::= <"{"> <Ignored> Selection+ <"}"> <Ignored>
    Selection ::= Field <Ignored> | FragmentSpread <Ignored> | InlineFragment <Ignored>
 */

func parseSelectionSet(lexer *Lexer) *SelectionSet {
    var selections []Selection
    lineNum := lexer.GetLineNum() 
    // parse variable def until token is "}"
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        selections = append(selections, parseSelection(lexer).(Selection))
    }

    return &SelectionSet{lineNum, selections}
}

func parseSelection(lexer *Lexer) interface{} {
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

func parseField(lexer *Lexer) *Field {
    var alias *Alias
    var fieldName *FieldName
    var arguments []*Argument
    var directives []*Directive
    var selectionSet *SelectionSet
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
        fmt.Printf("\033[34marguments: %v \033[0m\n", arguments)
    }

    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        directives = parseDirectives(lexer)
        fmt.Printf("\033[34mdirectives: %v \033[0m\n", directives)
    }

    // more SelectionSet
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        fmt.Printf("\033[34m into more SelectionSet: \033[0m\n")
        lexer.NextTokenIs(TOKEN_LEFT_BRACE)
        selectionSet = parseSelectionSet(lexer)
        lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
        fmt.Printf("\033[34m out more SelectionSet: \033[0m\n")
    }
    return &Field{lineNum, alias, fieldName, arguments, directives, selectionSet}
}

func parseFieldName(lexer *Lexer) *FieldName {
    name := parseName(lexer)
    return &FieldName{lexer.GetLineNum(), name}
}


/* 
    Parse FragmentSpread 
    FragmentSpread ::= <"..."> FragmentName <Ignored> Directives?
    FragmentName ::= Name
 */

func parseFragment(lexer *Lexer) interface{} {
    return nil
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
func parseFragmentDefinition(lexer *Lexer) *FragmentDefinition {
    return nil
}