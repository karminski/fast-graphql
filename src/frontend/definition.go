// definition.go

package frontend


type Definition interface{
    GetDefinitionType() string
}

// Definition should be 
var _ Definition = (*TypeSystemDefinition)(nil)
var _ Definition = (*EnumTypeDefinition)(nil)
var _ Definition = (*OperationDefinition)(nil)
// var _ Definition = (*FragmentDefinition)(nil)


/**
 * TypeSystemDefinition
 * TypeSystemDefinition ::= TypeDefinition | InterfaceDefinition | UnionDefinition | SchemaDefinition | EnumTypeDefinition | InputDefinition | DirectiveDefinition | TypeExtensionDefinition | ScalarDefinition
 */

const TypeSystemDefinitionType = "TypeSystemDefinition"

type TypeSystemDefinition struct{
    LineNum int
    TypeName *Name 

}

func (typeSystemDefinition *TypeSystemDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}


/**
 * EnumTypeDefinition
 * EnumDefinition ::= Description? <"enum"> <Ignored> TypeName <Ignored> Directives? <Ignored> <"{"> EnumValuesDefinition <"}"> <Ignored>
 * EnumValuesDefinition ::= EnumValueDefinition+
 * EnumValueDefinition ::= Description? <Ignored> EnumValue <Ignored> Directives? 
 * EnumValue ::= Name
 */
type EnumTypeDefinition struct {
    LineNum        int 
    Name          *Name
    Description    StringValue
    Directives  []*Directive
    Values      []*EnumValueDefinition
}

func (enumTypeDefinition EnumTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

type EnumValueDefinition struct {
    LineNum        int
    Description    StringValue 
    Value         *Name
    Directives  []*Directive
}


/**
 * OperationDefinition
 * OperationDefinition ::= <Ignored> OperationType? <Ignored> OperationName? <Ignored> VariableDefinitions? <Ignored> Directives? SelectionSet
 */

const OperationDefinitionType = "OperationDefinition"

type OperationDefinition struct {
    LineNum                int
    OperationType         *OperationType
    Name                  *Name
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

type FieldDefinition struct {
    LineNum     int 
    TokenName   string
}


type VariableDefinition struct {
    LineNum        int 
    Variable      *Name
    Type           Type
    DefaultValue   Value
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

type NamedType Name

type ListType struct {
    LineNum     int
    Type        []Type
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


type Value interface {
}

var _ Value = (*VariableValue)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ListValue)(nil)
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

type StringValue struct {
    LineNum int
    Value   string
}

type BooleanValue struct {
    LineNum int 
    Value   bool
}

type NullValue struct {
    LineNum int
}

type EnumValue struct {
    LineNum int 
    Value   *Name
}

type ListValue struct {
    LineNum int 
    Value   []Value
}

type ObjectValue struct {
    LineNum int 
    Value []*ObjectField
}

/**
 * ObjectField Definition
 * ObjectField ::= <Ignored> Name <":"> <Ignored> Value <Ignored>
 */

type ObjectField struct {
    LineNum    int
    Name       *Name 
    Value      Value
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
 * Arguments Section
 * Arguments ::= <"("> <Ignored> Argument+ <Ignored> <")"> <Ignored>
 * Argument ::= Name <Ignored> <":"> <Ignored> Value <Ignored>
 */

type Argument struct {
    LineNum  int
    Name    *Name
    Value    Value
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
    Name            *Name
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


/**
 * FragmentSpread Section
 * FragmentSpread     ::= <"..."> <Ignored> FragmentName <Ignored> Directives? <Ignored>
 * FragmentDefinition ::= <"fragment"> <Ignored> FragmentName <Ignored> TypeCondition <Ignored> Directives? <Ignored> SelectionSet <Ignored>
 * FragmentName       ::= Name # but not <"on">
 * TypeCondition      ::= <"on"> <Ignored> NamedType <Ignored>
 */
const FragmentDefinitionType = "FragmentDefinition"

type FragmentSpread struct {
    LineNum       int
    Name         *Name
    Directives []*Directive
}

func (fragmentSpread *FragmentSpread) GetSelectionSet() *SelectionSet {
    return nil 
}

type FragmentDefinition struct {
	LineNum          int
	Name 	        *Name
	TypeCondition   *Name
	Directives    []*Directive
	SelectionSet    *SelectionSet
}

func (fragmentDefinition *FragmentDefinition) GetDefinitionType() string {
    return FragmentDefinitionType
}

type InlineFragment struct {
    LineNum          int 
    TypeCondition    *Name
    Directives     []*Directive
    SelectionSet     *SelectionSet
}

func (inlineFragment *InlineFragment) GetSelectionSet() *SelectionSet {
    return inlineFragment.SelectionSet
}

type TypeCondition struct {
    LineNum     int
    TypeName    *TypeName
}


/**
 * SchemaDefinition Expression 
 * SchemaDefinition ::= "schema" Ignored Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored
 * SchemaExtension  ::= "extend" Ignored "schema" Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored | "extend" Ignored "schema" Directives Ignored
 *
 */
type SchemaDefinition struct {
    LineNum                    int 
    Directives              []*Directive
    OperationTypeDefinition    OperationTypeDefinition
}

type SchemaExtension struct {
    LineNum                    int 
    Directives              []*Directive 
    OperationTypeDefinition    OperationTypeDefinition
}

/**
 * OperationTypeDefinition Expression 
 * OperationTypeDefinition ::= OperationType Ignored ":" Ignored NamedType Ignored
 *
 */
type OperationTypeDefinition struct {
    LineNum         int 
    OperationType  *OperationType
    NamedType      *NamedType
}


/**
 * TypeDefinition Expression 
 * TypeDefinition            ::= ScalarTypeDefinition | ObjectTypeDefinition | InterfaceTypeDefinition | UnionTypeDefinition | EnumTypeDefinition | InputObjectTypeDefinition
 * TypeExtension             ::= ScalarTypeExtension | ObjectTypeExtension | InterfaceTypeExtension | UnionTypeExtension | EnumTypeExtension | InputObjectTypeExtension
 *
 */


/**
 * ScalarTypeDefinition      ::= Description? Ignored "scalar" Ignored Name Ignored Directives? Ignored
 * ScalarTypeExtension       ::= "extend" Ignored "scalar" Ignored Name Ignored Directives Ignored
 *
 */
type ScalarTypeDefinition struct {
    LineNum       int 
    Description   StringValue
    Name         *Name 
    Directives []*Directive
}

type ScalarTypeExtension struct {
    LineNum       int 
    Name         *Name 
    Directives []*Directive
}

/**
 * ObjectTypeDefinition      ::= Description? Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition? gnored
 * ObjectTypeExtension       ::= "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces Ignored
 *
 */
type ObjectTypeDefinition struct {
    LineNum                 int 
    Name                   *Name 
    ImplementsInterfaces    ImplementsInterfaces
    Directives           []*Directive
    FieldsDefinition        FieldsDefinition
}

type ObjectTypeExtension struct {
    LineNum                 int 
    Name                   *Name 
    ImplementsInterfaces    ImplementsInterfaces
    Directives           []*Directive
    FieldsDefinition        FieldsDefinition
}


/**
 * ImplementsInterfaces      ::= "implements" Ignored "&"? Ignored NamedType Ignored | ImplementsInterfaces Ignored "&" Ignored NamedType Ignored
 *
 */
type ImplementsInterfaces struct {
    LineNum                 int 
    ImplementsInterfaces     ImplementsInterfaces
    NamedType               *NamedType
}


/**
 * InterfaceTypeDefinition   ::= Description? Ignored "interface" Ignored Name Ignored Directives? Ignored FieldsDefinition? Ignored
 * InterfaceTypeExtension    ::= "extend" Ignored "interface" Ignored Name Ignored Directives? Ignored FieldsDefinition Ignored | "extend" Ignored "interface" Ignored Name Ignored Directives Ignored 
 *
 */
type InterfaceTypeDefinition struct {
    LineNum             int 
    Description         StringValue
    Name               *Name 
    Directives       []*Directive
    FieldsDefinition    FieldsDefinition
}

type InterfaceTypeExtension struct {
    LineNum             int 
    Name               *Name 
    Directives       []*Directive
    FieldsDefinition    FieldsDefinition
}


/**
 * UnionTypeDefinition       ::= Description? Ignored "union" Ignored Name Ignored Directives? Ignored UnionMemberTypes? Ignored
 * UnionMemberTypes          ::= "=" Ignored "|"? Ignored NamedType Ignored | UnionMemberTypes Ignored "|" Ignored NamedType Ignored
 * UnionTypeExtension        ::= "extend" Ignored "union" Ignored Name Ignored Directives? Ignored UnionMemberTypes? Ignored | "extend" Ignored "union" Ignored Name Ignored Directives Ignored
 *
 */
type UnionTypeDefinition struct {
    LineNum             int
    Description         StringValue
    Name               *Name 
    Directives       []*Directive
    UnionMemberTypes    UnionMemberTypes
}

type UnionMemberTypes struct {
    LineNum           int
    UnionMemberTypes  UnionMemberTypes
    NamedType        *NamedType
}

type UnionTypeExtension struct {
    LineNum             int
    Name               *Name
    Directives       []*Directive
    UnionMemberTypes    UnionMemberTypes
}


/**
 * InputObjectTypeDefinition ::= Description? Ignored "input" Ignored Name Ignored Directives? Ignored InputFieldsDefinition? Ignored
 * InputFieldsDefinition     ::= "{" Ignored InputValueDefinition+ Ignored "}" Ignored
 * InputObjectTypeExtension  ::= "extend" Ignored "input" Ignored Name Ignored Directives? Ignored InputFieldsDefinition Ignored | "extend" Ignored "input" Ignored Name Ignored Directives Ignored 
 *
 */
type InputObjectTypeDefinition struct {
    LineNum                  int
    Description              StringValue
    Name                    *Name 
    Directives            []*Directive
    InputFieldsDefinition    InputFieldsDefinition
}

type InputFieldsDefinition []*InputValueDefinition

type InputObjectTypeExtension struct {
    LineNum                   int
    Name                     *Name
    Directives             []*Directive
    InputFieldsDefinition     InputFieldsDefinition
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
type DirectiveDefinition struct {
    LineNum                int
    Description            StringValue
    Name                  *Name 
    ArgumentsDefinition []*InputValueDefinition
    DirectiveLocations     DirectiveLocations
}

type DirectiveLocations []*DirectiveLocation

type DirectiveLocation     string


/**
 * FieldsDefinition Expression
 * FieldsDefinition ::= "{" Ignored FieldDefinition+ Ignored "}"
 * FieldDefinition  ::= Description? Ignored Name Ignored ArgumentsDefinition? Ignored ":" Ignored Type Ignored Directives? Ignored
 *
 */
type FieldsDefinition []*FieldDefinition

type FieldDefinition struct {
    LineNum                int
    Description            StringValue
    Name                  *Name 
    ArgumentsDefinition []*InputValueDefinition
    Type                   Type
    Directives          []*Directive
}


/**
 * ArgumentsDefinition Expression 
 * ArgumentsDefinition  ::= "(" Ignored InputValueDefinition+ Ignored ")" Ignored
 * InputValueDefinition ::= Description? Ignored Name Ignored ":" Ignored Type Ignored DefaultValue? Ignored Directives? Ignored
 *
 */

type ArgumentsDefinition []*InputValueDefinition

type InputValueDefinition struct {
    LineNum         int 
    Description     StringValue
    Name           *Name 
    Type            Type
    DefaultValue    Value
    Directives   []*Directive
}

