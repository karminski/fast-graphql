// executor.go
package backend

import (
    "fast-graphql/src/frontend"
    "fmt"
    "errors"
    "log"
    "reflect"
    // "strconv"

    "github.com/davecgh/go-spew/spew"

)

const DUMP=false

type Request struct {
    // GraphQL Schema config for server side
    Schema Schema 

    // GraphQL Query string from client side
    Query string 
}

type Result struct {
    Data  interface{} `json:"data"`
    Error string      `json:"error"`
}

// get field name string from Field in AST
func getFieldName(field *frontend.Field) string {
        return field.FieldName.Name.Value
}

func Execute(request Request) (*Result) {
    // debugging
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    // process input
    document := frontend.Compile(request.Query)
    // @todo: THE DOCUMENT NEED VALIDATE!
    
    if DUMP {
        fmt.Printf("\n")
        fmt.Printf("\033[33m    [DUMP] Document:  \033[0m\n")
        spewo.Dump(document)
        fmt.Printf("\033[33m    [DUMP] Request:  \033[0m\n")
        spewo.Dump(request)
    }

    // get top layer SelectionSet.Fields and request.Schema.ObjectFields
    operationDefinition, _ := document.GetOperationDefinition()
    selectionSet := operationDefinition.SelectionSet
    // selectionSetFields := getSelectionSetFields(selectionSet)
    objectFields       := request.Schema.GetQueryObjectFields()
    // execute
    resolveSelectionSet(selectionSet, objectFields)

    return nil
}




// get name mapped Fields from SelectionSet
func getSelectionSetFields(selectionSet *frontend.SelectionSet) map[string]*frontend.Field {
    fields := make(map[string]*frontend.Field)
    selections := selectionSet.GetSelections()
    for _, selection := range selections {
        field := selection.(*frontend.Field)
        fieldName := field.FieldName.Name.Value
        fields[fieldName] = field
    }
    return fields
}


func resolveSelectionSet(selectionSet *frontend.SelectionSet, objectFields ObjectFields) {
    selections        := selectionSet.GetSelections()
    for _, selection := range selections {
        // prepare data
        field := selection.(*frontend.Field)
        fieldName := getFieldName(field)
        // resolve Field
        resolveField(fieldName, field, objectFields)
    }
}


func getResolveFunction(fieldName string, objectFields ObjectFields) ResolveFunction {
    resolveFunction := objectFields[fieldName].ResolveFunction
    // build in type, provide default resolve function
    if resolveFunction == nil {
        return nil
    }
    return resolveFunction
}

func resolveField(fieldName string, field *frontend.Field, objectFields ObjectFields) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveField  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}
    fmt.Printf("\033[33m    [DUMP] fieldName:  \033[0m\n")
    spewo.Dump(fieldName)
    fmt.Printf("\033[33m    [DUMP] field:  \033[0m\n")
    spewo.Dump(field)
    fmt.Printf("\033[33m    [DUMP] objectFields:  \033[0m\n")
    spewo.Dump(objectFields)

    if _, ok := objectFields[fieldName]; !ok {
        err := "resolveField(): input document field name "+fieldName+" does not defined in schema."
        return nil, errors.New(err)
    }
    
    resolveFunction := getResolveFunction(fieldName, objectFields)
    if resolveFunction == nil {
        fmt.Printf("\033[33m    [HIT!] resolveFunction == nil  \033[0m\n")
        resolveFunction = func (i interface{}) (interface{}, error) {
            return i, nil
        }
    }
    resolvedData, _ := resolveFunction(nil) // @todo: resolveFunction need input query param
    fmt.Printf("\033[33m    [DUMP] objectFields[documentFieldName]:  \033[0m\n")
    spewo.Dump(objectFields[fieldName])
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    // resolve sub-Field
    targetSelectionSet := field.SelectionSet
    targetObjectField := objectFields[fieldName]
    targetObjectFieldType := objectFields[fieldName].Type
    // go
    resolveSubField(targetSelectionSet, targetObjectField, targetObjectFieldType, resolvedData)
    return resolvedData, nil
}


func resolveSubField(selectionSet *frontend.SelectionSet, objectField *ObjectField, targetType FieldType, resolvedData interface{}) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveSubField  \033[0m\n")
    // get resolve target type

    if _, ok := targetType.(*List); ok {
        resolveListData(selectionSet, objectField, resolvedData)
    } 

    if _, ok := targetType.(*Scalar); ok {
        resolveScalarData(selectionSet, objectField, resolvedData)
    }

    if _, ok := targetType.(*Object); ok {
        resolveObjectData()
    }
    
}

func resolveListData(selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveListData  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    resolvedDataValue := reflect.ValueOf(resolvedData)
    targetObjectFields := objectField.Type.(*List).Payload.(*Object).Fields
    // traverse list
    for i:=0; i<resolvedDataValue.Len(); i++ {
        resolvedDataElement := resolvedDataValue.Index(i).Interface()
        fmt.Printf("\033[33m    [DUMP] resolvedDataElement:  \033[0m\n")
        spewo.Dump(resolvedDataElement)
        fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
        spewo.Dump(objectField)
        fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
        spewo.Dump(selectionSet)
        // execute
        resolveSelectionSet(selectionSet, targetObjectFields)
        
    }
    return nil, nil

}

func resolveScalarData(selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveScalarData  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    // call resolve function
    fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
    spewo.Dump(selectionSet)
    fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    spewo.Dump(objectField)
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)
    return nil, nil
}

func resolveObjectData() (interface{}, error) {
fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveObjectData  \033[0m\n")
    return nil, nil
}


// types

type Type interface {
    GetName() string
}

type FieldType interface {
    GetName() string
}

// List types

type List struct {
    Payload Type `json:payload`
}

func (list *List) GetName() string {
    return fmt.Sprintf("%v", list.Payload)
}

func NewList(i Type) *List {
    list := &List{}

    if i == nil {
        log.Fatal("NewList() input is nil")
        return list
    }

    list.Payload = i
    return list
}

// scalar definition

type ScalarTemplate struct {
    Name        string `json:name`
    Description string `json:description`
    ResolveFunction ResolveFunction `json:"-"`
}

type Scalar struct {
    Name            string          `json:name`
    Description     string          `json:description`
    ResolveFunction ResolveFunction `json:"-"`
}

func (scalar *Scalar) GetName() string {
    return scalar.Name
}


func NewScalar(scalarTemplate ScalarTemplate) *Scalar {
    scalar := &Scalar{}

    // check scalar template
    if scalarTemplate.Name == "" {
        err := "scalarTemplate.Name is not defined"
        log.Fatal(err)
    }

    scalar.Name        = scalarTemplate.Name
    scalar.Description = scalarTemplate.Description

    // @todo: Scalar should provide serialize, parse value, parse literal function.
    
    return scalar
}

// scalar types
var Int = NewScalar(ScalarTemplate{
    Name: "Int",
    Description: "GraphQL Int type",
    // ResolveFunction: func (i interface{}) (interface{}, error) {
    //     if intValue, err := strconv.Atoi(i.(string)); err == nil {
    //         return intValue, nil
    //     }
    //     return nil, nil
    // },
    ResolveFunction: func (i interface{}) (interface{}, error) {
        return i, nil
    },
})

var String = NewScalar(ScalarTemplate{
    Name: "String",
    Description: "GraphQL String type",
    ResolveFunction: func (i interface{}) (interface{}, error) {
        return i, nil
    },
})

// Object Syntax

type ObjectFields map[string]*ObjectField

type ObjectTemplate struct {
    Name   string 
    Fields ObjectFields
}

type Object struct {
    Name   string
    Fields ObjectFields
}

func (object *Object) GetName() string {
    return object.Name
}

func (object *Object) GetFields() ObjectFields {
    return object.Fields
} 

type ObjectField struct {
    Name            string               `json:name`
    Type            FieldType            `json:type`
    Description     string               `json:description`
    Arguments       *Arguments           `json:arguments`    
    ResolveFunction ResolveFunction      `json:"-"`
}

type Arguments map[string]*Argument

type Argument struct {
    Name string    `json:name` 
    Type FieldType `json:type`
}

type ResolveFunction func(i interface{}) (interface{}, error)


func NewObject(objectTemplate ObjectTemplate) (*Object, error) {
    object := &Object{}

    // check object input
    if objectTemplate.Name == "" {
        err := errors.New("ObjectTemplate.Name is not defined")
        return nil, err
    }
    
    object.Name = objectTemplate.Name
    object.Fields = objectTemplate.Fields
    return object, nil
}

// Schema Syntax

type SchemaTemplate struct {
    Query        *Object
    Mutation     *Object 
    Subscription *Object 
}

type Schema struct {
    Query        *Object
    Mutation     *Object 
    Subscription *Object 
}

func (schema *Schema) GetQueryObject() *Object {
    return schema.Query
}

func (schema *Schema) GetQueryObjectFields() ObjectFields {
    return schema.Query.Fields
}

func (schema *Schema) GetMutationObject() *Object {
    return schema.Mutation
}

func (schema *Schema) GetSubscriptionObject() *Object {
    return schema.Subscription
}

func NewSchema(schemaTemplate SchemaTemplate) (Schema, error) {
    schema := Schema{}

    // check query
    if schemaTemplate.Query == nil {
        err := errors.New("SchemaTemplate.Query is not defined")
        return schema, err
    }

    // fill schema
    schema.Query = schemaTemplate.Query


    return schema, nil
}