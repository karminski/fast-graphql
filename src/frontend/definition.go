// definition.go

package frontend


type Definition interface{
    GetDefinitionType() string
}

// Definition should be 
var _ Definition = (*OperationDefinition)(nil)
// var _ Definition = (*FragmentDefinition)(nil)


/**
 * TypeSystemDefinition
 * TypeSystemDefinition ::= TypeDefinition | InterfaceDefinition | UnionDefinition | SchemaDefinition | EnumDefinition | InputDefinition | DirectiveDefinition | TypeExtensionDefinition | ScalarDefinition
 */

const TypeSystemDefinitionType = "TypeSystemDefinition"

type TypeSystemDefinition struct {
    LineNum       int 
}

func (typeSystemDefinition *TypeSystemDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

/**
 * OperationDefinition
 * OperationDefinition ::= <Ignored> OperationType? <Ignored> OperationName? <Ignored> VariableDefinitions? <Ignored> Directives? SelectionSet
 */

const OperationDefinitionType = "OperationDefinition"

type OperationDefinition struct {
    LineNum                int
    OperationType         *OperationType
    OperationName         *OperationName
    VariableDefinitions []*VariableDefinition
    Directives          []*Directive
    SelectionSet          *SelectionSet
}

const (
    OperationTypeQuery        = "query"
    OperationTypeMutation     = "mutation"
    OperationTypeSubscription = "subscription"
)

type OperationType struct {
    LineNum       int 
    OperationName string // anonymous operation if Token is empty.
    Operation     int
}

func (operationDefinition *OperationDefinition) GetDefinitionType() string {
    return OperationDefinitionType
}

func (operationDefinition *OperationDefinition) IsQuery() bool {
    if operationDefinition.OperationType.Operation == TOKEN_QUERY {
        return true
    }
    return false
}

func (operationDefinition *OperationDefinition) IsMutation() bool {
    if operationDefinition.OperationType.Operation == TOKEN_MUTATION {
        return true
    }
    return false
}

func (operationDefinition *OperationDefinition) IsSubscription() bool {
    if operationDefinition.OperationType.Operation == TOKEN_SUBSCRIPTION {
        return true
    }
    return false
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
    Value      Value
}


/**
 * Value Definition
 * Value ::= VariableName | IntValue | FloatValue | ListValue | StringValue | BooleanValue | EnumValue | ObjectValue
 * VariableName ::=  
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

var _ Value = (*VariableValue)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ObjectValue)(nil)

type VariableValue struct {
    LineNum int 
    Value   string   
}

type IntValue struct {
    LineNum int
    Value   int
}

type FloatValue struct {
    LineNum int 
    Value   float64
}

type ListValue struct {
    LineNum int 
    Value   []Value
}

type StringValue struct {
    LineNum int
    Value   string
}

type BooleanValue struct {
    LineNum int 
    Value   bool
}

type EnumValue struct {
    LineNum int 
    Value   string
}

type ObjectValue struct {
    LineNum int 
    Fields []*ObjectField
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
    Value      Value
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

func (selectionSet *SelectionSet) GetSelections() []Selection {
    return selectionSet.Selections
}

type Selection interface {
    GetSelectionSet() *SelectionSet
}

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

func (field *Field) GetSelectionSet() *SelectionSet {
    return field.SelectionSet
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

const FragmentDefinitionType = "FragmentDefinition"

type FragmentDefinition struct {
    LineNum int
}

func (fragmentDefinition *FragmentDefinition) GetDefinitionType() string {
    return FragmentDefinitionType
}

type FragmentSpread struct {
    LineNum    int
    FragmentName    *FragmentName
    Directives      []*Directive
}

// FragmentSpread does not have SelectionSet section.
func (fragmentSpread *FragmentSpread) GetSelectionSet() *SelectionSet {
    return nil 
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

func (inlineFragment *InlineFragment) GetSelectionSet() *SelectionSet {
    return inlineFragment.SelectionSet
}


/**
 * ## TypeCondition
 * TypeCondition ::= <"on"> <Ignored> TypeName <Ignored>
 */

type TypeCondition struct {
    LineNum     int
    TypeName    *TypeName
}



