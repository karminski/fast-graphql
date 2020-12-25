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

    // "github.com/davecgh/go-spew/spew"

)


/**
 * Lexical Tokens Expression
 * Token                ::= Punctuator | Name | IntValue | FloatValue | StringValue
 * Punctuator           ::= "!" | "$" | "(" | ")" | "..." | ":" | "=" | "@" | "[" | "]" | "{" | "|" | "}"
 * Name                 ::= [_A-Za-z][_0-9A-Za-z]*
 * IntValue             ::= IntegerPart 
 * IntegerPart          ::= NegativeSign? Ignored "0" Ignored | NegativeSign? Ignored NonZeroDigit Ignored Digit+? Ignored 
 * NegativeSign         ::= "-"
 * Digit                ::= [0-9]
 * NonZeroDigit         ::= Digit - "0"
 * FloatValue           ::= IntegerPart FractionalPart | IntegerPart ExponentPart | IntegerPart FractionalPart ExponentPart 
 * FractionalPart       ::= "." Digit+
 * ExponentPart         ::= ExponentIndicator Sign? Digit+
 * ExponentIndicator    ::= "e" | "E"
 * Sign                 ::= "+" | "-"
 * StringValue          ::= '"""' '"""' | '"' '"' | '"' StringCharacter '"' | '"""' BlockStringCharacter '"""'
 * StringCharacter      ::= SourceCharacter - '"' | SourceCharacter - "\" | SourceCharacter - LineTerminator | "\u" EscapedUnicode | "\" EscapedCharacter 
 * EscapedUnicode       ::= [#x0000-#xFFFF] 
 * EscapedCharacter     ::= '"' | '\' | '/' | 'b' | 'f' | 'n' | 'r' | 't'
 * BlockStringCharacter ::= SourceCharacter - '"""' | SourceCharacter - '\"""' | '\"""' 
 *
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

func parseNumberValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseNumberValue  \033[0m\n")

    var isFloat = func(token string) bool {
        i := strings.Index(token, ".")
        if i < 0 {
            return false
        }
        return true
    }

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


/**
 * Document Expression
 * Document ::= Ignored Definition+ Ignored
 * 
 */
func parseDocument(lexer *Lexer) (*Document, error) {
    fmt.Println("\n\n\033[33m////////////////////////////////////////// Parser Start ///////////////////////////////////////\033[0m\n")
    fmt.Printf("\033[31m[INTO] func parseDocument  \033[0m\n")

    var document Document
    var err      error

    // LastLineNum
    document.LastLineNum = lexer.GetLineNum()
    // Definition+
    if document.Definitions, err = parseDefinitions(lexer); err != nil {
        return nil, err
    }
    return &document, nil
}

func isDocumentEnd(tokenType int) bool {
    if tokenType == TOKEN_EOF {
        return true
    }
    return false
}


/**
 * Definition Expression 
 * Definition           ::= ExecutableDefinition | TypeSystemDefinition | TypeSystemExtension
 * ExecutableDefinition ::= OperationDefinition | FragmentDefinition
 *
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


/**
 * OperationDefinition Expression
 * OperationDefinition ::= SelectionSet | Ignored OperationType Ignored Name? Ignored VariableDefinitions? Ignored Directives? SelectionSet 
 * OperationType       ::= "query" | "mutation" | "subscription" 
 *
 */
func parseOperationDefinition(lexer *Lexer) (*OperationDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseOperationDefinition  \033[0m\n")

    var operationDefinition OperationDefinition
    var err                 error

    // LineNum
    operationDefinition.LineNum = lexer.GetLineNum()
    // short query operation
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        goto SHORT_QUERY_OPERATION
    }
    // OperationType
    operationDefinition.OperationType, operationDefinition.OperationTypeName = parseOperationType(lexer)
    // Name?
    if lexer.LookAhead() == TOKEN_IDENTIFIER {
        if operationDefinition.Name, err = parseName(lexer); err != nil {
            return nil, err
        }
    }
    // VariableDefinitions?
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if operationDefinition.VariableDefinitions, err = parseVariableDefinitions(lexer); err != nil {
            return nil, err
        }
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if operationDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    // SelectionSet
    SHORT_QUERY_OPERATION:
        if operationDefinition.SelectionSet, err = parseSelectionSet(lexer); err != nil {
            return nil, err
        }

    return &operationDefinition, nil
}

func parseOperationType(lexer *Lexer) (int, string) {
    fmt.Printf("\033[31m[INTO] func parseOperationType  \033[0m\n")

    var operation int
    switch lexer.LookAhead() {
    case TOKEN_QUERY:        // operation "query"
        lexer.NextTokenIs(TOKEN_QUERY)
        operation = TOKEN_QUERY
    case TOKEN_MUTATION:     // operation "mutation"
        lexer.NextTokenIs(TOKEN_MUTATION)
        operation = TOKEN_MUTATION
    case TOKEN_SUBSCRIPTION: // operation "subscription"
        lexer.NextTokenIs(TOKEN_SUBSCRIPTION)
        operation = TOKEN_SUBSCRIPTION
    default:                 // anonymous operation must be query operation
        lexer.NextTokenIs(TOKEN_QUERY)
        operation = TOKEN_QUERY 
    }
    return operation, tokenNameMap[operation]
}


/**
 * SelectionSet Expression
 * SelectionSet ::= "{" Ignored Selection+ Ignored "}" Ignored
 * Selection    ::= Field Ignored | FragmentSpread Ignored | InlineFragment Ignored
 * 
 */
func parseSelectionSet(lexer *Lexer) (*SelectionSet, error) {
    fmt.Printf("\033[31m[INTO] func parseSelectionSet  \033[0m\n")

    var selectionSet SelectionSet

    // LineNum
    selectionSet.LineNum = lexer.GetLineNum() 
    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // Selection+
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        var selectionInterface interface{}
        var err                error 
        if selectionInterface, err = parseSelection(lexer); err != nil {
            return nil, err
        }
        selectionSet.Selections = append(selectionSet.Selections, selectionInterface.(Selection))
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return &selectionSet, nil
}

func parseSelection(lexer *Lexer) (interface{}, error) {
    fmt.Printf("\033[31m[INTO] func parseSelection  \033[0m\n")

    switch lexer.LookAhead() {
    case TOKEN_DOTS:
        lexer.NextTokenIs(TOKEN_DOTS)
        if lexer.LookAhead() == TOKEN_IDENTIFIER { 
            return parseFragmentSpread(lexer)
        } else{
            return parseInlineFragment(lexer)
        }
    default:
        return parseField(lexer)
    }
}

/**
 * Field Expression 
 * Field ::= Alias? Ignored Name Ignored Arguments? Ignored Directives? Ignored SelectionSet? Ignored
 *
 * Alias Expression 
 * Alias ::= Name Ignored ":" Ignored
 * 
 */
func parseField(lexer *Lexer) (*Field, error) {
    fmt.Printf("\033[31m[INTO] func parseField  \033[0m\n")

    var field           Field
    var err             error

    // lineNum
    field.LineNum = lexer.GetLineNum()

    //  Alias & Name
    var name *Name
    if name, err = parseName(lexer); err != nil {
        return nil ,err
    }
    if lexer.LookAhead() == TOKEN_COLON { // suffix is ":", it's Alias
        field.Alias = &Alias{lexer.GetLineNum(), name}
        if field.Name, err = parseName(lexer); err != nil {
            return nil, err
        }
    } else {
        field.Name = name
    } 

    // Arguments
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if field.Arguments, err = parseArguments(lexer); err != nil {
            return nil, err
        }
    }

    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        if field.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }

    // SelectionSet
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        fmt.Printf("\033[33m into more SelectionSet: \033[0m\n")
        if field.SelectionSet, err = parseSelectionSet(lexer); err != nil {
            return nil, err
        }
        fmt.Printf("\033[33m out more SelectionSet: \033[0m\n")
    }
    return &field, nil
}

/**
 * Arguments Expression
 * Arguments ::= "(" Ignored Argument+ Ignored ")" Ignored
 * Argument  ::= Name Ignored ":" Ignored Value Ignored
 *
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

    // LineNum
    argument.LineNum = lexer.GetLineNum()
    // Name
    if argument.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    lexer.NextTokenIs(TOKEN_COLON)
    // Value
    if argument.Value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return &argument, nil
}

/**
 * FragmentSpread Expression 
 * FragmentSpread     ::= "..." Ignored FragmentName Ignored Directives? Ignored
 * InlineFragment     ::= "..." Ignored TypeCondition? Ignored Directives? Ignored SelectionSet? Ignored
 * FragmentDefinition ::= "fragment" Ignored FragmentName Ignored TypeCondition Ignored Directives? Ignored SelectionSet Ignored
 * FragmentName       ::= Name - "on"
 * TypeCondition      ::= "on" Ignored NamedType Ignored
 *
 */
func parseFragmentSpread(lexer *Lexer) (*FragmentSpread, error) {
    fmt.Printf("\033[31m[INTO] func parseFragmentSpread  \033[0m\n")

    var fragmentSpread FragmentSpread
    var err            error

    // "..."
    lexer.NextTokenIs(TOKEN_DOTS)
    // FragmentName
    if fragmentSpread.Name, err = parseFragmentName(lexer); err != nil {
        return nil, err
    }
    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        if fragmentSpread.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    return &fragmentSpread, nil
}

func parseInlineFragment(lexer *Lexer) (*InlineFragment, error) {
    fmt.Printf("\033[31m[INTO] func parseInlineFragment  \033[0m\n")

    var inlineFragment InlineFragment
    var err            error

    // "..."
    lexer.NextTokenIs(TOKEN_DOTS)
    // TypeCondition
    if lexer.LookAhead() == TOKEN_ON {
        if inlineFragment.TypeCondition, err = parseTypeCondition(lexer); err != nil {
            return nil, err
        }
    }
    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        if inlineFragment.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    // SelectionSet
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if inlineFragment.SelectionSet, err = parseSelectionSet(lexer); err != nil {
            return nil, err
        }
    }
    return &inlineFragment, nil
}


func parseFragmentDefinition(lexer *Lexer) (*FragmentDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseInlineFragment  \033[0m\n")

    var fragmentDefinition FragmentDefinition
    var err                error

    // "fragment"
    lexer.NextTokenIs(TOKEN_FRAGMENT)
    // FragmentName
    if fragmentDefinition.Name, err = parseFragmentName(lexer); err != nil {
        return nil, err
    }
    // TypeCondition
    if fragmentDefinition.TypeCondition, err = parseTypeCondition(lexer); err != nil {
        return nil, err
    }
    // Directives
    if lexer.LookAhead() == TOKEN_AT {
        if fragmentDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    // SelectionSet
    if fragmentDefinition.SelectionSet, err = parseSelectionSet(lexer); err != nil {
        return nil, err
    }
    return &fragmentDefinition, nil
}

func parseFragmentName(lexer *Lexer) (*Name, error) {
    fmt.Printf("\033[31m[INTO] func parseInlineFragment  \033[0m\n")

    var name *Name
    var err   error
    shouldNotBe := map[string]bool{tokenNameMap[TOKEN_ON]: true}

    // Name
    if name, err = parseName(lexer); err != nil {
        return nil, err
    }
    if _, ok := shouldNotBe[name.Value]; ok {
        err = errors.New("parseFragmentName(): FragmentName should not be reserved word 'on'.")
        return nil, err
    }
    return name, nil
}

func parseTypeCondition(lexer *Lexer) (*Name, error) {
    fmt.Printf("\033[31m[INTO] func parseInlineFragment  \033[0m\n")

    var name *Name 
    var err   error

    // "on"
    lexer.NextTokenIs(TOKEN_ON)
    // NamedType
    if name, err = parseName(lexer); err != nil {
        return nil, err
    }
    return name, nil
}


/**
 * Value Expression 
 * Value            ::= Variable | IntValue | FloatValue | StringValue | BooleanValue | NullValue | EnumValue | ListValue | ObjectValue
 * BooleanValue     ::= "true" | "false"
 * NullValue        ::= "null"
 * EnumValue        ::=  Name - "true" | Name - "false" | Name - "null"
 * ListValue        ::= "[" "]" | "[" Value+ "]"
 * ObjectValue      ::= "{" "}" | "{" ObjectField+ "}"
 * ObjectField      ::= Ignored Name Ignored ":" Ignored Value Ignored
 *
 */
func parseValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseValue  \033[0m\n")

    var value Value
    var err error

    switch token := lexer.LookAhead(); {
    case token == TOKEN_VAR_PREFIX:  // Variable, start with "$"
        if value, err = parseVariable(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_NUMBER:      // number, include IntValue, FloatValue
        if value, err = parseNumberValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_QUOTE:       // StringValue, "string"
        if value, err = parseStringValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_TRIQUOTE:    // StringValue, """string""" 
        if value, err = parseStringValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_HEXQUOTE:    // StringValue, """""" 
        if value, err = parseStringValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_DUOQUOTE:    // StringValue, ""
        if value, err = parseStringValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_TRUE:        // BooleanValue
        if value, err = parseBooleanValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_FALSE:       // BooleanValue
        if value, err = parseBooleanValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_NULL:        // NullValue
        if value, err = parseNullValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_IDENTIFIER:  // EnumValue
        if value, err = parseEnumValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_LEFT_BRACKET: // ListValue
        if value, err = parseListValue(lexer); err != nil {
            return nil, err
        }
    case token == TOKEN_LEFT_BRACE:   // ObjectValue
        if value, err = parseObjectValue(lexer); err != nil {
            return nil, err
        }
    default:
        err := errors.New("parseValue(): unexpected value type '" + tokenNameMap[token] + "'." )
        return nil, err
    }
    return value, nil
}

func parseBooleanValue(lexer *Lexer) (BooleanValue, error) {
    fmt.Printf("\033[31m[INTO] func parseBooleanValue  \033[0m\n")

    if lexer.LookAhead() == TOKEN_TRUE {
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
    fmt.Printf("\033[31m[INTO] func parseBooleanValue  \033[0m\n")

    var enumValue EnumValue
    var err       error
    shouldNotBe := map[string]bool{tokenNameMap[TOKEN_TRUE]: true, tokenNameMap[TOKEN_FALSE]: true, tokenNameMap[TOKEN_NULL]: true}

    // LineNum
    enumValue.LineNum = lexer.GetLineNum()
    // Name
    if enumValue.Value, err = parseName(lexer); err != nil {
        return enumValue, err
    }
    if _, ok := shouldNotBe[enumValue.Value.Value]; ok {
        err = errors.New("parseEnumValue(): enum value can not be 'true' or 'false' or 'null'.")
        return enumValue, err
    }
    return enumValue, nil
}

func parseListValue(lexer *Lexer) (ListValue, error) {
    fmt.Printf("\033[31m[INTO] func parseListValue  \033[0m\n")

    var listValue ListValue

    // "["
    lexer.NextTokenIs(TOKEN_LEFT_BRACKET)
    // Value+
    for lexer.LookAhead() != TOKEN_RIGHT_BRACKET {
        var value Value 
        var err   error
        if value, err = parseValue(lexer); err != nil {
            return listValue, err
        }
        listValue.Value = append(listValue.Value, value)
    }
    // "]"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET)
    return listValue, nil
}

func parseObjectValue(lexer *Lexer) (ObjectValue, error) {
    fmt.Printf("\033[31m[INTO] func parseObjectValue  \033[0m\n")

    var objectValue ObjectValue

    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // ObjectField+
    for lexer.LookAhead() != TOKEN_RIGHT_BRACE {
        var objectField *ObjectField
        var err          error
        if objectField, err = parseObjectField(lexer); err != nil {
            return objectValue, err
        }
        objectValue.Value = append(objectValue.Value, objectField)
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE) 
    return objectValue, nil
}

func parseObjectField(lexer *Lexer) (*ObjectField, error) {
    fmt.Printf("\033[31m[INTO] func parseObjectField  \033[0m\n")

    var objectField ObjectField
    var err         error

    // Name
    if objectField.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    // Value
    if objectField.Value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return &objectField, nil
}


/**
 * VariableDefinitions Expression 
 * VariableDefinitions ::= "(" VariableDefinition+ ")"
 * VariableDefinition  ::= Variable ":" Ignored Type Ignored DefaultValue? Ignored
 * Variable            ::= "$" Name
 * DefaultValue        ::= "=" Ignored Value
 *
 */
func parseVariableDefinitions(lexer *Lexer) ([]*VariableDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinitions  \033[0m\n")

    var VariableDefinitions []*VariableDefinition

    // "("
    lexer.NextTokenIs(TOKEN_LEFT_PAREN) 
    // VariableDefinition+
    for lexer.LookAhead() != TOKEN_RIGHT_PAREN {
        var variableDefinition *VariableDefinition
        var err                 error
        if variableDefinition, err = parseVariableDefinition(lexer); err != nil {
            return nil, err
        }
        VariableDefinitions = append(VariableDefinitions, variableDefinition)
    }
    // ")"
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return VariableDefinitions, nil
}

func parseVariableDefinition(lexer *Lexer) (*VariableDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseVariableDefinition  \033[0m\n")

    var variableDefinition VariableDefinition
    var err                error

    // LineNum
    variableDefinition.LineNum = lexer.GetLineNum()
    // Variable
    if variableDefinition.Variable, err = parseVariable(lexer); err != nil {
        return nil, err
    }
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    // Type
    if variableDefinition.Type, err = parseType(lexer); err != nil {
        return nil, err
    }
    // DefaultValue?
    if lexer.LookAhead() == TOKEN_EQUAL {
        if variableDefinition.DefaultValue, err = parseDefaultValue(lexer); err != nil {
            return nil, err
        }
    }
    return &variableDefinition, nil
}

func parseVariable(lexer *Lexer) (*Name, error) {
    fmt.Printf("\033[31m[INTO] func parseVariable  \033[0m\n")

    // "$"
    lexer.NextTokenIs(TOKEN_VAR_PREFIX)
    // Name
    return parseName(lexer)    
}

func parseDefaultValue(lexer *Lexer) (Value, error) {
    fmt.Printf("\033[31m[INTO] func parseDefaultValue  \033[0m\n")

    var value Value 
    var err   error

    // "="
    lexer.NextTokenIs(TOKEN_EQUAL)
    // Value
    if value, err = parseValue(lexer); err != nil {
        return nil, err
    }
    return value, nil
}


/**
 * Type Expression 
 * Type        ::= NamedType | ListType | NonNullType
 * NamedType   ::= Name
 * ListType    ::= "[" Type+ "]"
 * NonNullType ::= NamedType "!" | ListType "!"
 *
 */
func parseType(lexer *Lexer) (Type, error) {
    fmt.Printf("\033[31m[INTO] func parseType  \033[0m\n")

    var typeRet Type
    var err     error

    // NamedType & ListType
    switch lexer.LookAhead() {
    case TOKEN_IDENTIFIER:   // NamedType
        if typeRet, err = parseNamedType(lexer); err != nil {
            return nil, err
        }
    case TOKEN_LEFT_BRACKET: // ListType, start with "["
        if typeRet, err = parseListType(lexer); err != nil {
            return nil, err
        }
    }
    // NonNullType
    if lexer.LookAhead() == TOKEN_NOT_NULL {
        if typeRet, err = parseNonNullType(lexer, typeRet); err != nil {
            return nil, err
        }   
    }
    return typeRet, nil
}

var parseNamedType = parseName 

func parseListType(lexer *Lexer) (ListType, error) {
    fmt.Printf("\033[31m[INTO] func parseListType  \033[0m\n")

    var listType ListType

    // "["
    lexer.NextTokenIs(TOKEN_LEFT_BRACKET) 
    // Type+
    for lexer.LookAhead() != TOKEN_RIGHT_BRACKET {
        var typeRet Type
        var err     error
        if typeRet, err = parseType(lexer); err != nil {
            return listType, err
        }
        listType.Type = append(listType.Type, typeRet)
    }
    // "]"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACKET)
    return listType, nil
}

func parseNonNullType(lexer *Lexer, previousType Type) (NonNullType, error) {
    fmt.Printf("\033[31m[INTO] func parseNonNullType  \033[0m\n")

    return NonNullType{lexer.GetLineNum(), previousType}, nil
}


/**
 * Directives Expression
 * Directives ::= Directive+
 * Directive  ::= "@" Ignored Name Ignored Arguments? Ignored
 *
 */
func parseDirectives(lexer *Lexer) ([]*Directive, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectives  \033[0m\n")

    var directives []*Directive
    
    // Directive+
    for lexer.LookAhead() == TOKEN_AT { 
        var directive    *Directive
        var err           error
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
    var err       error

    // LineNum
    directive.LineNum = lexer.GetLineNum()
    // "@"
    lexer.NextTokenIs(TOKEN_AT)
    // Name
    if directive.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // Arguments?
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if directive.Arguments, err = parseArguments(lexer); err != nil {
            return nil, err
        }
    }
    return &directive, nil
}

/**
 * TypeSystemDefinition Expression 
 * TypeSystemDefinition ::= SchemaDefinition | TypeDefinition | DirectiveDefinition
 * TypeSystemExtension  ::= SchemaExtension | TypeExtension
 *
 * Description Expression 
 * Description ::= StringValue
 * 
 */
// @todo: unfinished
func parseTypeSystemDefinition(lexer *Lexer) (Definition, error) {
    fmt.Printf("\033[31m[INTO] func parseTypeSystemDefinition  \033[0m\n")

    var description StringValue 
    var err         error

    // Description
    nextToken := lexer.LookAhead()
    if nextToken == TOKEN_QUOTE || nextToken == TOKEN_DUOQUOTE ||  nextToken == TOKEN_TRIQUOTE || nextToken == TOKEN_HEXQUOTE {
        if description, err = parseDescription(lexer); err != nil {
            return nil, err
        }
    }
    // SchemaDefinition | TypeDefinition | DirectiveDefinition
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

func parseTypeSystemExtension(lexer *Lexer) {

}

var parseDescription = parseStringValue


/**
 * SchemaDefinition Expression 
 * SchemaDefinition ::= "schema" Ignored Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored
 * SchemaExtension  ::= "extend" Ignored "schema" Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored | "extend" Ignored "schema" Directives Ignored
 *
 */
func parseSchemaDefinition(lexer *Lexer) (SchemaDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseSchemaDefinition  \033[0m\n")

    var schemaDefinition SchemaDefinition
    var err              error

    // LineNum
    schemaDefinition.LineNum = lexer.GetLineNum()
    // "schema"
    lexer.NextTokenIs(TOKEN_SCHEMA)
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if schemaDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return schemaDefinition, err
        }
    }
    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // OperationTypeDefinition+
    if schemaDefinition.OperationTypeDefinition, err = parseOperationTypeDefinition(lexer); err != nil {
        return schemaDefinition, err
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return schemaDefinition, nil
}

func parseSchemaExtension(lexer *Lexer) (SchemaExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseSchemaExtension  \033[0m\n")

    var schemaExtension SchemaExtension
    var err              error

    // LineNum
    schemaExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "schema"
    lexer.NextTokenIs(TOKEN_SCHEMA)
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if schemaExtension.Directives, err = parseDirectives(lexer); err != nil {
            return schemaExtension, err
        }
    }
    // "{"?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        // "{"
        lexer.NextTokenIs(TOKEN_LEFT_BRACE)
        // OperationTypeDefinition+
        if schemaExtension.OperationTypeDefinition, err = parseOperationTypeDefinition(lexer); err != nil {
            return schemaExtension, err
        }
        // "}"
        lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    }
    return schemaExtension, nil
}

/**
 * OperationTypeDefinition Expression 
 * OperationTypeDefinition ::= OperationType Ignored ":" Ignored NamedType Ignored
 *
 */
func parseOperationTypeDefinition(lexer *Lexer) (OperationTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseSchemaExtension  \033[0m\n")

    var operationTypeDefinition OperationTypeDefinition
    var err                     error

    // LineNum
    operationTypeDefinition.LineNum = lexer.GetLineNum()
    // OperationType
    if operationTypeDefinition.OperationType, operationTypeDefinition.OperationTypeName = parseOperationType(lexer); err != nil {
        return operationTypeDefinition, err
    } 
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    // NamedType
    if operationTypeDefinition.NamedType, err = parseNamedType(lexer); err != nil {
        return operationTypeDefinition ,err
    }
    return operationTypeDefinition, nil
}


/**
 * TypeDefinition Expression 
 * TypeDefinition            ::= ScalarTypeDefinition | ObjectTypeDefinition | InterfaceTypeDefinition | UnionTypeDefinition | EnumTypeDefinition | InputObjectTypeDefinition
 * TypeExtension             ::= ScalarTypeExtension | ObjectTypeExtension | InterfaceTypeExtension | UnionTypeExtension | EnumTypeExtension | InputObjectTypeExtension
 *
 */
// @todo: unfinished
// func parseTypeDefinition(lexer *Lexer) (TypeDefinition, error) {
// 
// }
// 
// 
// func parseTypeExtension(lexer *Lexer) (TypeExtension, error) {
// 
// }

/**
 * ScalarTypeDefinition      ::= Description? Ignored "scalar" Ignored Name Ignored Directives? Ignored
 * ScalarTypeExtension       ::= "extend" Ignored "scalar" Ignored Name Ignored Directives Ignored
 *
 */
func parseScalarTypeDefinition(lexer *Lexer) (ScalarTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseScalarTypeDefinition  \033[0m\n")

    var scalarTypeDefinition ScalarTypeDefinition
    var err                  error

    // LineNum
    scalarTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "scalar"
    lexer.NextTokenIs(TOKEN_SCALAR)
    // Name
    if scalarTypeDefinition.Name, err = parseName(lexer); err != nil {
        return scalarTypeDefinition, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if scalarTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return scalarTypeDefinition, err
        }
    }
    return scalarTypeDefinition, nil
}

func parseScalarTypeExtension(lexer *Lexer) (ScalarTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseScalarTypeExtension  \033[0m\n")

    var scalarTypeExtension ScalarTypeExtension
    var err                 error

    // LineNum
    scalarTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "scalar"
    lexer.NextTokenIs(TOKEN_SCALAR)
    // Name 
    if scalarTypeExtension.Name, err = parseName(lexer); err != nil {
        return scalarTypeExtension, err
    }
    // Directives
    if scalarTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
        return scalarTypeExtension, err
    }
    return scalarTypeExtension, nil
}   


/**
 * ObjectTypeDefinition      ::= Description? Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition? gnored
 * ObjectTypeExtension       ::= "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces Ignored
 *
 */
func parseObjectTypeDefinition(lexer *Lexer) (ObjectTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseObjectTypeDefinition  \033[0m\n")

    var objectTypeDefinition ObjectTypeDefinition
    var err                  error

    // LineNum
    objectTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "type"
    lexer.NextTokenIs(TOKEN_TYPE)
    // Name
    if objectTypeDefinition.Name, err = parseName(lexer); err != nil {
        return objectTypeDefinition, err
    }
    // ImplementsInterfaces?
    if lexer.LookAhead() == TOKEN_IMPLEMENTS {
        if objectTypeDefinition.ImplementsInterfaces, err = parseImplementsInterfaces(lexer); err != nil {
            return objectTypeDefinition, err
        }
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if objectTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return objectTypeDefinition, err
        }
    }
    // FieldsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if objectTypeDefinition.FieldsDefinition, err = parseFieldsDefinition(lexer); err != nil {
            return objectTypeDefinition, err
        }
    }
    return objectTypeDefinition, nil
}

func parseObjectTypeExtension(lexer *Lexer) (ObjectTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseObjectTypeDefinition  \033[0m\n")

    var objectTypeExtension ObjectTypeExtension
    var err                 error

    // LineNum
    objectTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "type"
    lexer.NextTokenIs(TOKEN_TYPE)
    // Name
    if objectTypeExtension.Name, err = parseName(lexer); err != nil {
        return objectTypeExtension, err
    }
    // ImplementsInterfaces?
    if lexer.LookAhead() == TOKEN_IMPLEMENTS {
        if objectTypeExtension.ImplementsInterfaces, err = parseImplementsInterfaces(lexer); err != nil {
            return objectTypeExtension, err
        }
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if objectTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
            return objectTypeExtension, err
        }
    }
    // FieldsDefinition
    if objectTypeExtension.FieldsDefinition, err = parseFieldsDefinition(lexer); err != nil {
        return objectTypeExtension, err
    }
    return objectTypeExtension, nil
}


/**
 * ImplementsInterfaces      ::= "implements" Ignored "&"? Ignored NamedType Ignored | ImplementsInterfaces Ignored "&" Ignored NamedType Ignored
 *
 */
func parseImplementsInterfaces(lexer *Lexer) (ImplementsInterfaces, error) {
    fmt.Printf("\033[31m[INTO] func parseImplementsInterfaces  \033[0m\n")

    var implementsInterfaces ImplementsInterfaces

    // LineNum
    implementsInterfaces.LineNum = lexer.GetLineNum()
    namedTypeCounter := 0
    for {
        var namedType *Name
        var err        error
        // "&"?
        if lexer.LookAhead() == TOKEN_AND {
            lexer.NextTokenIs(TOKEN_AND)
            namedTypeCounter --
        }
        // NamedType
        if lexer.LookAhead() == TOKEN_IDENTIFIER {
            if namedTypeCounter > 0 { 
                break
            }
        }
        if namedType, err = parseNamedType(lexer); err != nil {
            return implementsInterfaces, err
        }
        implementsInterfaces.NamedTypes = append(implementsInterfaces.NamedTypes, namedType)
        namedTypeCounter++
    }
    return implementsInterfaces, nil
}


/**
 * InterfaceTypeDefinition   ::= Description? Ignored "interface" Ignored Name Ignored Directives? Ignored FieldsDefinition? Ignored
 * InterfaceTypeExtension    ::= "extend" Ignored "interface" Ignored Name Ignored Directives? Ignored FieldsDefinition Ignored | "extend" Ignored "interface" Ignored Name Ignored Directives Ignored 
 *
 */
func parseInterfaceTypeDefinition(lexer *Lexer) (InterfaceTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseInterfaceTypeDefinition  \033[0m\n")

    var interfaceTypeDefinition InterfaceTypeDefinition
    var err                     error

    // LineNum
    interfaceTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "interface"
    lexer.NextTokenIs(TOKEN_INTERFACE)
    // Name
    if interfaceTypeDefinition.Name, err = parseName(lexer); err != nil {
        return interfaceTypeDefinition, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if interfaceTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return interfaceTypeDefinition, err
        }
    }
    // FieldsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if interfaceTypeDefinition.FieldsDefinition, err = parseFieldsDefinition(lexer); err != nil {
            return interfaceTypeDefinition, err
        }
    }
    return interfaceTypeDefinition, nil
}

func parseInterfaceTypeExtension(lexer *Lexer) (InterfaceTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseInterfaceTypeExtension  \033[0m\n")

    var interfaceTypeExtension InterfaceTypeExtension
    var err                    error

    // LineNum
    interfaceTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "interface"
    lexer.NextTokenIs(TOKEN_INTERFACE)
    // Name
    if interfaceTypeExtension.Name, err = parseName(lexer); err != nil {
        return interfaceTypeExtension, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if interfaceTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
            return interfaceTypeExtension, err
        }
    }
    // FieldsDefinition
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if interfaceTypeExtension.FieldsDefinition, err = parseFieldsDefinition(lexer); err != nil {
            return interfaceTypeExtension, err
        }
    }
    return interfaceTypeExtension, nil
}


/**
 * UnionTypeDefinition       ::= Description? Ignored "union" Ignored Name Ignored Directives? Ignored UnionMemberTypes? Ignored
 * UnionMemberTypes          ::= "=" Ignored "|"? Ignored NamedType Ignored | UnionMemberTypes Ignored "|" Ignored NamedType Ignored
 * UnionTypeExtension        ::= "extend" Ignored "union" Ignored Name Ignored Directives? Ignored UnionMemberTypes? Ignored | "extend" Ignored "union" Ignored Name Ignored Directives Ignored
 *
 */

func parseUnionTypeDefinition(lexer *Lexer) (UnionTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseUnionTypeDefinition  \033[0m\n")

    var unionTypeDefinition UnionTypeDefinition
    var err                 error

    // LineNum
    unionTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "union"
    lexer.NextTokenIs(TOKEN_UNION)
    // Name
    if unionTypeDefinition.Name, err = parseName(lexer); err != nil {
        return unionTypeDefinition, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if unionTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return unionTypeDefinition, err
        }
    }
    // UnionMemberTypes?
    if lexer.LookAhead() == TOKEN_EQUAL {
        if unionTypeDefinition.UnionMemberTypes, err = parseUnionMemberTypes(lexer); err != nil {
            return unionTypeDefinition, err
        }
    }
    return unionTypeDefinition, nil
}

func parseUnionMemberTypes(lexer *Lexer) (UnionMemberTypes, error) {
    fmt.Printf("\033[31m[INTO] func parseUnionMemberTypes  \033[0m\n")

    var unionMemberTypes UnionMemberTypes

    // LineNum
    unionMemberTypes.LineNum = lexer.GetLineNum()
    // "=" 
    lexer.NextTokenIs(TOKEN_EQUAL)
    // UnionMemberTypes
    namedTypeCounter := 0
    for {
        var namedType *Name
        var err        error
        // "|"?
        if lexer.LookAhead() == TOKEN_VERTICAL_BAR {
            lexer.NextTokenIs(TOKEN_VERTICAL_BAR)
            namedTypeCounter --
        }
        // NamedType
        if lexer.LookAhead() == TOKEN_IDENTIFIER {
            if namedTypeCounter > 0 { 
                break
            }
        }
        if namedType, err = parseNamedType(lexer); err != nil {
            return unionMemberTypes, err
        }
        unionMemberTypes.NamedTypes = append(unionMemberTypes.NamedTypes, namedType)
        namedTypeCounter++
    }
    return unionMemberTypes, nil
}

func parseUnionTypeExtension(lexer *Lexer) (UnionTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseUnionTypeExtension  \033[0m\n")

    var unionTypeExtension UnionTypeExtension
    var err                    error

    // LineNum
    unionTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "union"
    lexer.NextTokenIs(TOKEN_UNION)
    // Name
    if unionTypeExtension.Name, err = parseName(lexer); err != nil {
        return unionTypeExtension, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if unionTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
            return unionTypeExtension, err
        }
    }
    // UnionMemberTypes
    if lexer.LookAhead() == TOKEN_EQUAL {
        if unionTypeExtension.UnionMemberTypes, err = parseUnionMemberTypes(lexer); err != nil {
            return unionTypeExtension, err
        }
    }
    return unionTypeExtension, nil
}

/**
 * EnumTypeDefinition        ::= Description? Ignored "enum" Ignored Name Ignored Directives? Ignored EnumValuesDefinition? Ignored
 * EnumValuesDefinition      ::= "{" Ignored EnumValueDefinition+ Ignored "}" Ignored
 * EnumValueDefinition       ::= Description? Ignored EnumValue Ignored Directives? Ignored
 * EnumTypeExtension         ::= "extend" Ignored "enum" Ignored Name Ignored Directives? Ignored EnumValuesDefinition Ignored | "extend" Ignored "enum" Ignored Name Ignored Directives Ignored
 *
 */
func parseEnumTypeDefinition(lexer *Lexer) (*EnumTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseEnumTypeDefinition  \033[0m\n")

    var enumTypeDefinition EnumTypeDefinition
    var err                error

    // LineNum
    enumTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "enum"
    lexer.NextTokenIs(TOKEN_ENUM)
    // Name
    if enumTypeDefinition.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if enumTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    // EnumValuesDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if enumTypeDefinition.EnumValuesDefinition, err = parseEnumValuesDefinition(lexer); err != nil {
            return nil, err
        }
    }
    return &enumTypeDefinition, nil
}

func parseEnumValuesDefinition(lexer *Lexer) ([]*EnumValueDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseEnumValuesDefinition  \033[0m\n")

    var enumValuesDefinition []*EnumValueDefinition

    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // enumValuesDefinition+
    for lexer.LookAhead() == TOKEN_RIGHT_BRACE { 
        var enumValueDefinition    *EnumValueDefinition
        var err           error
        if enumValueDefinition, err = parseEnumValueDefinition(lexer); err != nil {
            return nil, err
        }
        enumValuesDefinition = append(enumValuesDefinition, enumValueDefinition)
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return enumValuesDefinition, nil
}

func parseEnumValueDefinition(lexer *Lexer) (*EnumValueDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseEnumValueDefinition  \033[0m\n")

    var enumValueDefinition EnumValueDefinition
    var err                 error

    // LineNum
    enumValueDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // EnumValue
    if enumValueDefinition.EnumValue, err = parseEnumValue(lexer); err != nil {
        return nil, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if enumValueDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    return &enumValueDefinition, nil
}

func parseEnumTypeExtension(lexer *Lexer) (EnumTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseEnumTypeExtension  \033[0m\n")

    var enumTypeExtension EnumTypeExtension
    var err               error

    // LineNum
    enumTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "enum"
    lexer.NextTokenIs(TOKEN_ENUM)
    // Name
    if enumTypeExtension.Name, err = parseName(lexer); err != nil {
        return enumTypeExtension, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if enumTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
            return enumTypeExtension, err
        }
    }
    // EnumValuesDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if enumTypeExtension.EnumValuesDefinition, err = parseEnumValuesDefinition(lexer); err != nil {
            return enumTypeExtension, err
        }
    }
    return enumTypeExtension, nil
}

/**
 * InputObjectTypeDefinition ::= Description? Ignored "input" Ignored Name Ignored Directives? Ignored InputFieldsDefinition? Ignored
 * InputFieldsDefinition     ::= "{" Ignored InputValueDefinition+ Ignored "}" Ignored
 * InputObjectTypeExtension  ::= "extend" Ignored "input" Ignored Name Ignored Directives? Ignored InputFieldsDefinition Ignored | "extend" Ignored "input" Ignored Name Ignored Directives Ignored 
 *
 */
func parseInputObjectTypeDefinition(lexer *Lexer) (InputObjectTypeDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseInputObjectTypeDefinition  \033[0m\n")

    var inputObjectTypeDefinition InputObjectTypeDefinition
    var err               error

    // LineNum
    inputObjectTypeDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "input"
    lexer.NextTokenIs(TOKEN_INPUT)
    // Name
    if inputObjectTypeDefinition.Name, err = parseName(lexer); err != nil {
        return inputObjectTypeDefinition, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if inputObjectTypeDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return inputObjectTypeDefinition, err
        }
    }
    // InputFieldsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if inputObjectTypeDefinition.InputFieldsDefinition, err = parseInputFieldsDefinition(lexer); err != nil {
            return inputObjectTypeDefinition, err
        }
    }
    return inputObjectTypeDefinition, nil
}

func parseInputFieldsDefinition(lexer *Lexer) ([]*InputValueDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseInputFieldsDefinition  \033[0m\n")

    var inputFieldsDefinition []*InputValueDefinition

    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // InputValueDefinition+
    for lexer.LookAhead() == TOKEN_RIGHT_BRACE { 
        var inputValueDefinition *InputValueDefinition
        var err                   error
        if inputValueDefinition, err = parseInputValueDefinition(lexer); err != nil {
            return nil, err
        }
        inputFieldsDefinition = append(inputFieldsDefinition, inputValueDefinition)
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return inputFieldsDefinition, nil
}

func parseInputObjectTypeExtension(lexer *Lexer) (InputObjectTypeExtension, error) {
    fmt.Printf("\033[31m[INTO] func parseInputObjectTypeExtension  \033[0m\n")

    var inputObjectTypeExtension InputObjectTypeExtension
    var err               error

    // LineNum
    inputObjectTypeExtension.LineNum = lexer.GetLineNum()
    // "extend"
    lexer.NextTokenIs(TOKEN_EXTEND)
    // "input"
    lexer.NextTokenIs(TOKEN_INPUT)
    // Name
    if inputObjectTypeExtension.Name, err = parseName(lexer); err != nil {
        return inputObjectTypeExtension, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if inputObjectTypeExtension.Directives, err = parseDirectives(lexer); err != nil {
            return inputObjectTypeExtension, err
        }
    }
    // InputFieldsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_BRACE {
        if inputObjectTypeExtension.InputFieldsDefinition, err = parseInputFieldsDefinition(lexer); err != nil {
            return inputObjectTypeExtension, err
        }
    }
    return inputObjectTypeExtension, nil
}

/**
 * DirectiveDefinition Expression 
 * DirectiveDefinition         ::= Description? Ignored "directive" Ignored "@" Ignored Name Ignored ArgumentsDefinition? Ignored "on" Ignored DirectiveLocations Ignored
 * DirectiveLocations          ::= "|"? Ignored DirectiveLocation Ignored | DirectiveLocations Ignored "|" Ignored DirectiveLocation Ignored
 * DirectiveLocation           ::= ExecutableDirectiveLocation | TypeSystemDirectiveLocation
 * ExecutableDirectiveLocation ::= "QUERY" | "MUTATION" | "SUBSCRIPTION" | "FIELD" | "FRAGMENT_DEFINITION" | "FRAGMENT_SPREAD" | "INLINE_FRAGMENT" 
 * TypeSystemDirectiveLocation ::= "SCHEMA" | "SCALAR" | "OBJECT" | "FIELD_DEFINITION" | "ARGUMENT_DEFINITION" | "INTERFACE" | "UNION" | "ENUM" | "ENUM_VALUE" | "INPUT_OBJECT" | "INPUT_FIELD_DEFINITION" 
 *
 */
func parseDirectiveDefinition(lexer *Lexer) (DirectiveDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectiveDefinition  \033[0m\n")

    var directiveDefinition DirectiveDefinition
    var err               error

    // LineNum
    directiveDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // "directive"
    lexer.NextTokenIs(TOKEN_DIRECTIVE)
    // "@"
    lexer.NextTokenIs(TOKEN_AT)
    // Name
    if directiveDefinition.Name, err = parseName(lexer); err != nil {
        return directiveDefinition, err
    }
    // ArgumentsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if directiveDefinition.ArgumentsDefinition, err = parseArgumentsDefinition(lexer); err != nil {
            return directiveDefinition, err
        }
    }
    // "on"
    lexer.NextTokenIs(TOKEN_ON)
    // DirectiveLocations
    if directiveDefinition.DirectiveLocations, err = parseDirectiveLocations(lexer); err != nil {
        return directiveDefinition, err
    }
    return directiveDefinition, nil
}

func parseDirectiveLocations(lexer *Lexer) ([]string, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectiveLocations  \033[0m\n")

    var directiveLocations       []string

    // DirectiveLocations
    directiveLocationCounter  := 0
    for {
        var directiveLocation string
        var err               error
        // "|"?
        if lexer.LookAhead() == TOKEN_VERTICAL_BAR {
            lexer.NextTokenIs(TOKEN_VERTICAL_BAR)
            directiveLocationCounter --
        }
        // DirectiveLocation
        if lexer.LookAhead() == TOKEN_IDENTIFIER {
            if directiveLocationCounter > 0 { // this means the next token is not DirectiveLocation. there will a "|" between two DirectiveLocations 
                break
            }
        }
        if directiveLocation, err = parseDirectiveLocation(lexer); err != nil {
            return directiveLocations, err
        }
        directiveLocations = append(directiveLocations, directiveLocation)
        directiveLocationCounter++
    }
    return directiveLocations, nil
}

func parseDirectiveLocation(lexer *Lexer) (string, error) {
    fmt.Printf("\033[31m[INTO] func parseDirectiveLocation  \033[0m\n")

    executableDirectiveLocation := map[string]bool{
        "QUERY": true,
        "MUTATION": true,
        "SUBSCRIPTION": true,
        "FIELD": true,
        "FRAGMENT_DEFINITION": true,
        "FRAGMENT_SPREAD": true,
        "INLINE_FRAGMENT": true,}
    typeSystemDirectiveLocation :=  map[string]bool{
        "SCHEMA": true,
        "SCALAR": true,
        "OBJECT": true,
        "FIELD_DEFINITION": true,
        "ARGUMENT_DEFINITION": true,
        "INTERFACE": true,
        "UNION": true,
        "ENUM": true,
        "ENUM_VALUE": true,
        "INPUT_OBJECT": true,
        "INPUT_FIELD_DEFINITION": true,}

    // DirectiveLocation
    _, token := lexer.NextTokenIs(TOKEN_IDENTIFIER)
    _, ok1 := executableDirectiveLocation[token];
    _, ok2 := typeSystemDirectiveLocation[token];
    if !ok1 && !ok2 {
        err := errors.New("parseDirectiveLocation(): illegal DirectiveLocation '"+ token +"'.")
        return "", err
    }
    return token, nil
}

/**
 * FieldsDefinition Expression
 * FieldsDefinition ::= "{" Ignored FieldDefinition+ Ignored "}"
 * FieldDefinition  ::= Description? Ignored Name Ignored ArgumentsDefinition? Ignored ":" Ignored Type Ignored Directives? Ignored
 *
 */
func parseFieldsDefinition(lexer *Lexer) (FieldsDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseFieldsDefinition  \033[0m\n")

    var fieldsDefinition FieldsDefinition

    // "{"
    lexer.NextTokenIs(TOKEN_LEFT_BRACE)
    // FieldDefinition+
    for lexer.LookAhead() == TOKEN_RIGHT_BRACE { 
        var fieldDefinition *FieldDefinition
        var err              error
        if fieldDefinition, err = parseFieldDefinition(lexer); err != nil {
            return nil, err
        }
        fieldsDefinition = append(fieldsDefinition, fieldDefinition)
    }
    // "}"
    lexer.NextTokenIs(TOKEN_RIGHT_BRACE)
    return fieldsDefinition, nil
}

func parseFieldDefinition(lexer *Lexer) (*FieldDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseFieldDefinition  \033[0m\n")

    var fieldDefinition FieldDefinition
    var err             error

    // LineNum
    fieldDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // Name
    if fieldDefinition.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // ArgumentsDefinition?
    if lexer.LookAhead() == TOKEN_LEFT_PAREN {
        if fieldDefinition.ArgumentsDefinition, err = parseArgumentsDefinition(lexer); err != nil {
            return nil, err
        }
    }
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    // Type
    if fieldDefinition.Type, err = parseType(lexer); err != nil {
        return nil, err
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if fieldDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    return &fieldDefinition, nil
}


/**
 * ArgumentsDefinition Expression 
 * ArgumentsDefinition  ::= "(" Ignored InputValueDefinition+ Ignored ")" Ignored
 * InputValueDefinition ::= Description? Ignored Name Ignored ":" Ignored Type Ignored DefaultValue? Ignored Directives? Ignored
 *
 */
func parseArgumentsDefinition(lexer *Lexer) (ArgumentsDefinition, error) {
    fmt.Printf("\033[31m[INTO] func parseArgumentsDefinition  \033[0m\n")

    var argumentsDefinition ArgumentsDefinition

    // "("
    lexer.NextTokenIs(TOKEN_LEFT_PAREN)
    // FieldDefinition+
    for lexer.LookAhead() == TOKEN_RIGHT_PAREN { 
        var inputValueDefinition *InputValueDefinition
        var err              error
        if inputValueDefinition, err = parseInputValueDefinition(lexer); err != nil {
            return nil, err
        }
        argumentsDefinition = append(argumentsDefinition, inputValueDefinition)
    }
    // ")"
    lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
    return argumentsDefinition, nil    
}

func parseInputValueDefinition(lexer *Lexer) (*InputValueDefinition, error) {
     fmt.Printf("\033[31m[INTO] func parseInputValueDefinition  \033[0m\n")

    var inputValueDefinition InputValueDefinition
    var err                  error

    // LineNum
    inputValueDefinition.LineNum = lexer.GetLineNum()
    // Description? finished at parseTypeSystemDefinition()
    // Name
    if inputValueDefinition.Name, err = parseName(lexer); err != nil {
        return nil, err
    }
    // ":"
    lexer.NextTokenIs(TOKEN_COLON)
    // Type
    if inputValueDefinition.Type, err = parseType(lexer); err != nil {
        return nil, err
    }
    // DefaultValue?
    if lexer.LookAhead() == TOKEN_EQUAL {
        if inputValueDefinition.DefaultValue, err = parseDefaultValue(lexer); err != nil {
            return nil, err
        }
    }
    // Directives?
    if lexer.LookAhead() == TOKEN_AT {
        if inputValueDefinition.Directives, err = parseDirectives(lexer); err != nil {
            return nil, err
        }
    }
    return &inputValueDefinition, nil
}

