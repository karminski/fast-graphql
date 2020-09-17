// definition.go

package frontend

// 
const (

)


type Definition interface{}

type EmptyDefinition struct{}

type TypeSystemDefinition struct {
    LineNum       int 
}

type OperationDefinition struct {
    LineNum                int
    OperationType         *OperationType
    OperationName         *OperationName
    VariableDefinitions []*VariableDefinition
    Directives          []*Directive
    SelectionSet          *SelectionSet
}

type OperationType struct {
    LineNum       int 
    OperationName string // anonymous operation if Token is empty.
    Operation     int
}

/**
 * Name Definitation
 * Name ::= #"[_A-Za-z][_0-9A-Za-z]*"
 */

type Name struct {
    LineNum int
    Value   string
}

/**
 * OperationName Definitation
 * OperationName ::= Name
 */

type OperationName struct {
    LineNum int 
    Name    *Name
}

type FieldDefinition struct {
    LineNum     int 
    TokenName   string
}

/**
 * VariableDefinitions ::= <"("> VariableDefinition+ <")">
 * VariableDefinition ::= <Ignored> VariableName <":"> <Ignored> Type <Ignored> DefaultValue? <Ignored>
 */

type VariableDefinition struct {
    LineNum                 int 
    VariableName            *VariableName
    Type                    Type
    DefaultValue            *DefaultValue
}

type VariableName struct {
    LineNum int 
    Name    *Name
}



/**
 * Type Definition
 * Type ::= TypeName | ListType | NonNullType
 * TypeName ::= Name
 * ListType ::= <"["> Type <"]">
 * NonNullType ::= TypeName <"!"> | ListType <"!">
 */

type Type interface {

}

var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)

type TypeName struct {
    LineNum    int 
    Name       *Name
}

type NamedType struct {
    LineNum     int 
    Name        *Name
}

type ListType struct {
    LineNum     int
    Type        Type
}

type NonNullType struct {
    LineNum     int
    Type        Type
}


/**
 * DefaultValue ::= <"="> <Ignored> Value
 */

type DefaultValue struct {
    LineNum    int
    Value      *Value
}


/**
 * Value Definition
 * Value ::= VariableName | IntValue | FloatValue | ListValue | StringValue | BooleanValue | EnumValue | ObjectValue
 * IntValue ::= #"[\+\-0-9]+"
 * FloatValue ::= #"[\+\-0-9]+\.[0-9]"
 * ListValue ::= <"["> <"]"> | <"["> OneOrMoreValue <"]">
 * OneOrMoreValue ::= [Value <Ignored>]+
 * StringValue ::= <"\""><"\""> | <"\""> StringCharacter+ <"\"">
 * StringCharacter ::= #"[\x{9}\x{20}\x{21}\x{23}-\x{5B}\x{5D}-\uFFFF]" | "\\" "u" EscapedUnicode | "\\" EscapedCharacter
 * BooleanValue ::= "true" | "false"
 * EnumValue ::= #"(?!(true|false|null))[_A-Za-z][_0-9A-Za-z]*"
 * ObjectValue ::= <"{"> ObjectField <"}">
 */

type Value interface {

}

/**
 * ObjectField Definition
 * ObjectField ::= <Ignored> Name <":"> <Ignored> Value <Ignored>
 */

type ObjectField struct {
    LineNum    int
    Name       *Name 
    Value      *Value
}

/**
 * Directives Definition
 * Directives ::= Directive+
 * Directive ::= <"@"> Name Arguments? <Ignored>
 */

type Directive struct {
    LineNum     int
    Name        *Name
    Arguments   []*Argument
}


/**
 * Arguments ::= <"("> <Ignored> Argument+ <")">
 * Argument ::= ArgumentName <":"> <Ignored> ArgumentValue <Ignored>*
 * ArgumentName ::= Name 
 * ArgumentValue ::= Value | VariableName
 */

type Argument struct {
    LineNum         int
    ArgumentName    *ArgumentName
    ArgumentValue   *ArgumentValue
}

type ArgumentName struct {
    LineNum    int 
    Name       *Name
}

type ArgumentValue struct {
    LineNum    int 
    Value      *Value
}

/**
 * SelectionSet Definition
 * SelectionSet ::= <"{"> <Ignored> Selection+ <"}"> <Ignored>
 * Selection ::= Field <Ignored> | FragmentSpread <Ignored> | InlineFragment <Ignored>
 */

type SelectionSet struct {
    LineNum     int 
    Selections  []Selection
}

type Selection interface {}

var _ Selection = (*Field)(nil)
var _ Selection = (*FragmentSpread)(nil)
var _ Selection = (*InlineFragment)(nil)


/**
 * Parse Field 
 * Field ::= Alias? <Ignored> FieldName <Ignored> Arguments? <Ignored> Directives? SelectionSet?
 * Alias ::= Name <":">
 * FieldName ::= Name
 */
    
type Field struct {
    LineNum         int
    Alias           *Alias
    FieldName       *FieldName
    Arguments       []*Argument
    Directives      []*Directive
    SelectionSet    *SelectionSet
}

type Alias struct {
    LineNum    int 
    Name       *Name
}

type FieldName struct {
    LineNum    int 
    Name       *Name
}

/**
 * ## FragmentDefinition
 * FragmentDefinition ::= <"fragment"> <Ignored> FragmentName <Ignored> TypeCondition Directives? SelectionSet
 */

type FragmentDefinition interface {}

type FragmentSpread struct {
    LineNum    int
    FragmentName    *FragmentName
    Directives      []*Directive
}

type FragmentName struct {
    LineNum    int 
    Name       *Name
}

/**
 * InlineFragment Definition
 * InlineFragment ::= <"..."> <Ignored> TypeCondition? Directives? SelectionSet?
 */

type InlineFragment struct {
    LineNum          int 
    TypeCondition    *TypeCondition
    Directives       []*Directive
    SelectionSet     *SelectionSet
}


/**
 * ## TypeCondition
 * TypeCondition ::= <"on"> <Ignored> TypeName <Ignored>
 */

type TypeCondition struct {
    LineNum     int
    TypeName    *TypeName
}



