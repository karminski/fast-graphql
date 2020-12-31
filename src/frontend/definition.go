// definition.go

package frontend


/**
 * Name Definitation
 * Name ::= #"[_A-Za-z][_0-9A-Za-z]*"
 */

type Name struct {
    LineNum int
    Value   string
}


/**
 * Definition Expression 
 * Definition           ::= ExecutableDefinition | TypeSystemDefinition | TypeSystemExtension
 * ExecutableDefinition ::= OperationDefinition | FragmentDefinition
 *
 */
type Definition interface{
    GetDefinitionType() string
}

// ExecutableDefinition
var _ Definition = (*OperationDefinition)(nil)
var _ Definition = (*FragmentDefinition)(nil)
// TypeSystemDefinition
var _ Definition = (*SchemaDefinition)(nil)
var _ Definition = (*ScalarTypeDefinition)(nil)
var _ Definition = (*ObjectTypeDefinition)(nil)
var _ Definition = (*InterfaceTypeDefinition)(nil)
var _ Definition = (*UnionTypeDefinition)(nil)
var _ Definition = (*EnumTypeDefinition)(nil)
var _ Definition = (*InputObjectTypeDefinition)(nil)
var _ Definition = (*DirectiveDefinition)(nil)
// TypeSystemExtension
var _ Definition = (*SchemaExtension)(nil)
var _ Definition = (*ScalarTypeExtension)(nil)
var _ Definition = (*ObjectTypeExtension)(nil)
var _ Definition = (*InterfaceTypeExtension)(nil)
var _ Definition = (*UnionTypeExtension)(nil)
var _ Definition = (*EnumTypeExtension)(nil)
var _ Definition = (*InputObjectTypeExtension)(nil)


/**
 * OperationDefinition
 * OperationDefinition ::= <Ignored> OperationType? <Ignored> OperationName? <Ignored> VariableDefinitions? <Ignored> Directives? SelectionSet
 */

const OperationDefinitionType = "OperationDefinition"

type OperationDefinition struct {
    LineNum                int
    OperationType          int
    OperationTypeName      string
    Name                  *Name
    VariableDefinitions []*VariableDefinition
    Directives          []*Directive
    SelectionSet          *SelectionSet
}

const (
    OperationTypeQuery            = TOKEN_QUERY
    OperationTypeMutation         = TOKEN_MUTATION
    OperationTypeSubscription     = TOKEN_SUBSCRIPTION
    OperationTypeQueryName        = "query"
    OperationTypeMutationName     = "mutation"
    OperationTypeSubscriptionName = "subscription"
)

func (operationDefinition *OperationDefinition) GetDefinitionType() string {
    return OperationDefinitionType
}

func (operationDefinition *OperationDefinition) IsQuery() bool {
    if operationDefinition.OperationType == TOKEN_QUERY {
        return true
    }
    return false
}

func (operationDefinition *OperationDefinition) IsMutation() bool {
    if operationDefinition.OperationType == TOKEN_MUTATION {
        return true
    }
    return false
}

func (operationDefinition *OperationDefinition) IsSubscription() bool {
    if operationDefinition.OperationType == TOKEN_SUBSCRIPTION {
        return true
    }
    return false
}


/**
 * SelectionSet Expression 
 * SelectionSet ::= "{" Ignored Selection+ Ignored "}" Ignored
 * Selection    ::= Field Ignored | FragmentSpread Ignored | InlineFragment Ignored
 *
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
 * Field Expression 
 * Field ::= Alias? Ignored Name Ignored Arguments? Ignored Directives? Ignored SelectionSet? Ignored
 *
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


/**
 * Alias Expression 
 * Alias ::= Name Ignored ":" Ignored
 *
 */
type Alias struct {
    LineNum     int 
    Name       *Name
}


/**
 * Arguments Expression 
 * Arguments ::= "(" Ignored Argument+ Ignored ")" Ignored
 * Argument  ::= Name Ignored ":" Ignored Value Ignored
 *
 */
type Argument struct {
    LineNum  int
    Name    *Name
    Value    Value
}



/**
 * FragmentSpread Expression 
 * FragmentSpread     ::= "..." Ignored FragmentName Ignored Directives? Ignored
 * InlineFragment     ::= "..." Ignored TypeCondition? Ignored Directives? Ignored SelectionSet Ignored
 * FragmentDefinition ::= "fragment" Ignored FragmentName Ignored TypeCondition Ignored Directives? Ignored SelectionSet Ignored
 * FragmentName       ::= Name - "on"
 * TypeCondition      ::= "on" Ignored NamedType Ignored
 *
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
    Name            *Name
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
    NamedType  *NamedType
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
type Value interface {
}

var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*ObjectValue)(nil)

type Variable *Name

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
    Value   interface{}
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

type ObjectField struct {
    LineNum    int
    Name       *Name 
    Value      Value
}


/**
 * VariableDefinitions Expression 
 * VariableDefinitions ::= "(" VariableDefinition+ ")"
 * VariableDefinition  ::= Variable Ignored ":" Ignored Type Ignored DefaultValue? Ignored
 * Variable            ::= "$" Name
 * DefaultValue        ::= "=" Ignored Value
 *
 */
type VariableDefinition struct {
    LineNum        int 
    Variable      *Name
    Type           Type
    DefaultValue   Value
}


/**
 * Type Expression 
 * Type        ::= NamedType | ListType | NonNullType
 * NamedType   ::= Name
 * ListType    ::= "[" Type "]"
 * NonNullType ::= NamedType "!" | ListType "!"
 *
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
 * TypeSystemDefinition
 * TypeSystemDefinition ::= TypeDefinition | InterfaceDefinition | UnionDefinition | SchemaDefinition | EnumTypeDefinition | InputDefinition | DirectiveDefinition | TypeExtensionDefinition | ScalarDefinition
 */

const TypeSystemDefinitionType = "TypeSystemDefinition"


/**
 * SchemaDefinition Expression 
 * SchemaDefinition ::= "schema" Ignored Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored
 * SchemaExtension  ::= "extend" Ignored "schema" Directives? Ignored "{" Ignored OperationTypeDefinition+ Ignored "}" Ignored | "extend" Ignored "schema" Directives Ignored
 *
 */
const SchemaDefinitionType = "SchemaDefinition"

type SchemaDefinition struct {
    LineNum                     int 
    Directives               []*Directive
    OperationTypeDefinitions []*OperationTypeDefinition
}

func (schemaDefinition *SchemaDefinition) GetDefinitionType() string {
    return SchemaDefinitionType
}

const SchemaExtensionType = "SchemaExtension"

type SchemaExtension struct {
    LineNum                    int 
    Directives               []*Directive 
    OperationTypeDefinitions []*OperationTypeDefinition
}

func (schemaExtension *SchemaExtension) GetDefinitionType() string {
    return SchemaExtensionType
}

/**
 * OperationTypeDefinition Expression 
 * OperationTypeDefinition ::= OperationType Ignored ":" Ignored NamedType Ignored
 *
 */
type OperationTypeDefinition struct {
    LineNum            int 
    OperationType      int
    OperationTypeName  string 
    NamedType         *NamedType
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

func (scalarTypeDefinition *ScalarTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

const TypeExtensionType = "TypeExtension"

type ScalarTypeExtension struct {
    LineNum       int 
    Name         *Name 
    Directives []*Directive
}

func (scalarTypeExtension *ScalarTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
}

/**
 * ObjectTypeDefinition      ::= Description? Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition? gnored
 * ObjectTypeExtension       ::= "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives? Ignored FieldsDefinition Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces? Ignored Directives Ignored | "extend" Ignored "type" Ignored Name Ignored ImplementsInterfaces Ignored
 *
 */
type ObjectTypeDefinition struct {
    LineNum                 int 
    Description             StringValue
    Name                   *Name 
    ImplementsInterfaces   *ImplementsInterfaces
    Directives           []*Directive
    FieldsDefinition     []*FieldDefinition
}

func (objectTypeDefinition *ObjectTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

type ObjectTypeExtension struct {
    LineNum                 int 
    Name                   *Name 
    ImplementsInterfaces   *ImplementsInterfaces
    Directives           []*Directive
    FieldsDefinition     []*FieldDefinition
}

func (objectTypeExtension *ObjectTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
}


/**
 * ImplementsInterfaces      ::= "implements" Ignored "&"? Ignored NamedType Ignored | ImplementsInterfaces Ignored "&" Ignored NamedType Ignored
 *
 */
type ImplementsInterfaces struct {
    LineNum       int 
    NamedTypes []*NamedType
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
    FieldsDefinition []*FieldDefinition
}

func (interfaceTypeDefinition *InterfaceTypeDefinition)GetDefinitionType() string {
    return TypeSystemDefinitionType
} 

type InterfaceTypeExtension struct {
    LineNum             int 
    Name               *Name 
    Directives       []*Directive
    FieldsDefinition []*FieldDefinition
}

func (interfaceTypeExtension *InterfaceTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
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
    UnionMemberTypes   *UnionMemberTypes
}

func (unionTypeDefinition *UnionTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

type UnionMemberTypes struct {
    LineNum       int
    NamedTypes []*NamedType
}

type UnionTypeExtension struct {
    LineNum             int
    Name               *Name
    Directives       []*Directive
    UnionMemberTypes   *UnionMemberTypes
}

func (unionTypeExtension *UnionTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
}


/**
 * EnumTypeDefinition        ::= Description? Ignored "enum" Ignored Name Ignored Directives? Ignored EnumValuesDefinition? Ignored
 * EnumValuesDefinition      ::= "{" Ignored EnumValueDefinition+ Ignored "}" Ignored
 * EnumValueDefinition       ::= Description? Ignored EnumValue Ignored Directives? Ignored
 * EnumTypeExtension         ::= "extend" Ignored "enum" Ignored Name Ignored Directives? Ignored EnumValuesDefinition Ignored | "extend" Ignored "enum" Ignored Name Ignored Directives Ignored
 *
 */
type EnumTypeDefinition struct {
    LineNum                 int 
    Name                   *Name
    Description             StringValue
    Directives           []*Directive
    EnumValuesDefinition []*EnumValueDefinition
}

func (enumTypeDefinition EnumTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

type EnumValueDefinition struct {
    LineNum        int
    Description    StringValue 
    EnumValue      EnumValue
    Directives  []*Directive
}

type EnumTypeExtension struct {
    LineNum                 int
    Name                   *Name
    Directives           []*Directive
    EnumValuesDefinition []*EnumValueDefinition
}

func (enumTypeExtension *EnumTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
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
    InputFieldsDefinition []*InputValueDefinition
}

func (inputObjectTypeDefinition InputObjectTypeDefinition) GetDefinitionType() string {
    return TypeSystemDefinitionType
}

type InputObjectTypeExtension struct {
    LineNum                   int
    Name                     *Name
    Directives             []*Directive
    InputFieldsDefinition  []*InputValueDefinition
}

func (inputObjectTypeExtension *InputObjectTypeExtension) GetDefinitionType() string {
    return TypeExtensionType
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
const DirectiveDefinitionType = "DirectiveDefinition"

type DirectiveDefinition struct {
    LineNum                int
    Description            StringValue
    Name                  *Name 
    ArgumentsDefinition []*InputValueDefinition
    DirectiveLocations  []string
}

func (directiveDefinition DirectiveDefinition) GetDefinitionType() string {
    return DirectiveDefinitionType
}

type DirectiveLocations []string

type DirectiveLocation string


/**
 * FieldsDefinition Expression
 * FieldsDefinition ::= "{" Ignored FieldDefinition+ Ignored "}"
 * FieldDefinition  ::= Description? Ignored Name Ignored ArgumentsDefinition? Ignored ":" Ignored Type Ignored Directives? Ignored
 *
 */
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

