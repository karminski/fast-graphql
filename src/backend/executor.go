// executor.go
package backend

import (
    "fast-graphql/src/frontend"
    "fmt"
    "errors"
    "log"
    "reflect"
    // "strconv"
    "os"
    "github.com/davecgh/go-spew/spew"

)

const DUMP = false

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
    finalResult := Result{} 
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
        os.Exit(1)
    }

    // get top layer SelectionSet.Fields and request.Schema.ObjectFields
    operationDefinition, _ := document.GetOperationDefinition()
    selectionSet := operationDefinition.SelectionSet
    // selectionSetFields := getSelectionSetFields(selectionSet)
    objectFields       := request.Schema.GetQueryObjectFields()
    // execute
    resolvedResult, _ := resolveSelectionSet(selectionSet, objectFields, nil)
    fmt.Printf("\033[33m    [DUMP] resolvedResult:  \033[0m\n")
    spewo.Dump(resolvedResult)
    finalResult.Data = resolvedResult
    return &finalResult
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


func resolveSelectionSet(selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
    selections  := selectionSet.GetSelections()
    finalResult := make(map[string]interface{}, len(selections))
    for _, selection := range selections {
        // prepare data
        field := selection.(*frontend.Field)
        fieldName := getFieldName(field)
        // resolve Field
        resolvedResult, _ := resolveField(fieldName, field, objectFields, resolvedData)
        finalResult[fieldName] = resolvedResult   
    }
    return finalResult, nil
}


func getResolveFunction(fieldName string, objectFields ObjectFields) ResolveFunction {
    resolveFunction := objectFields[fieldName].ResolveFunction
    // build in type, provide default resolve function
    if resolveFunction == nil {
        return nil
    }
    return resolveFunction
}



func getArgumentsMap(arguments []*frontend.Argument) map[string]interface{} {
    argumentsMap := make(map[string]interface{}, len(arguments))
    for _, argument := range arguments {
        // detect value type & fill
        interfaceValue := argument.ArgumentValue.Value
        if val, ok := interfaceValue.(frontend.IntValue); ok {
            argumentsMap[argument.ArgumentName.Name.Value] = val.Value
        } else if val, ok := interfaceValue.(frontend.StringValue); ok {
            argumentsMap[argument.ArgumentName.Name.Value] = val.Value
        } else {
            argumentsMap[argument.ArgumentName.Name.Value] = nil
        }
    }
    return argumentsMap
}

func resolveField(fieldName string, field *frontend.Field, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
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
    
    // check resolve function or extend last resolved data 
    resolveFunction := getResolveFunction(fieldName, objectFields)

    if resolveFunction == nil {
        fmt.Printf("\033[33m    [HIT!] resolveFunction == nil  \033[0m\n")
        resolveFunction = func (p ResolveParams) (interface{}, error) {
            return p, nil
        }
    }
    // no context resovedData input, check GraphQL Request Arguments
    if resolvedData == nil {
        p := ResolveParams{}
        // GraphQL Request Arguments are avaliable
        if field.Arguments != nil {
            p.Arguments = getArgumentsMap(field.Arguments)
        }
        resolvedData, _ = resolveFunction(p) 
        fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
        spewo.Dump(resolvedData)
        os.Exit(1)
    }
    fmt.Printf("\033[33m    [DUMP] objectFields[documentFieldName]:  \033[0m\n")
    spewo.Dump(objectFields[fieldName])
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    // resolve sub-Field
    targetSelectionSet := field.SelectionSet
    targetObjectField := objectFields[fieldName]
    targetObjectFieldType := objectFields[fieldName].Type
    // go
    resolvedSubData, _ := resolveSubField(targetSelectionSet, targetObjectField, targetObjectFieldType, resolvedData)
    return resolvedSubData, nil
}


func resolveSubField(selectionSet *frontend.SelectionSet, objectField *ObjectField, targetType FieldType, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveSubField  \033[0m\n")
    // get resolve target type

    if _, ok := targetType.(*List); ok {
        return resolveListData(selectionSet, objectField, resolvedData)
    } 

    if _, ok := targetType.(*Scalar); ok {
        return resolveScalarData(selectionSet, objectField, resolvedData)
    }

    if _, ok := targetType.(*Object); ok {
        return resolveObjectData(selectionSet, objectField, resolvedData)
    }
    return nil, nil
    
}

func resolveListData(selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveListData  \033[0m\n")
    // spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    resolvedDataValue := reflect.ValueOf(resolvedData)
    targetObjectFields := objectField.Type.(*List).Payload.(*Object).Fields
    // allocate space for list data returns
    finalResult := make([]interface{}, 0, resolvedDataValue.Len())
    // traverse list
    for i:=0; i<resolvedDataValue.Len(); i++ {
        resolvedDataElement := resolvedDataValue.Index(i).Interface()
        // fmt.Printf("\033[33m    [DUMP] resolvedDataElement:  \033[0m\n")
        // spewo.Dump(resolvedDataElement)
        // fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
        // spewo.Dump(objectField)
        // fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
        // spewo.Dump(selectionSet)
        // execute
        selectionSetResult, _ := resolveSelectionSet(selectionSet, targetObjectFields, resolvedDataElement)
        finalResult = append(finalResult, selectionSetResult)
    }
    return finalResult, nil
}

func resolveScalarData(selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveScalarData  \033[0m\n")
    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    // fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
    // spewo.Dump(selectionSet)
    // fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    // spewo.Dump(objectField)
    // fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    // spewo.Dump(resolvedData)

    // call resolve function
    resolveFunction := objectField.Type.(*Scalar).ResolveFunction
    targetFieldName := objectField.Name
    r0 := getResolvedDataTargetFieldValue(resolvedData, targetFieldName)
    fmt.Printf("\033[33m    [DUMP] getResolvedDataTargetFieldValue result:  \033[0m\n")
    spewo.Dump(r0)
    // convert 
    p := ResolveParams{}
    p.Context = r0
    r1, _ := resolveFunction(p)
    fmt.Printf("\033[43;37m    [DUMP] resolveFunction result:  \033[0m\n")
    spewo.Dump(r1)
    return r1, nil
}

func resolveObjectData(selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}) (interface{}, error) {
    fmt.Printf("\n")
    fmt.Printf("\033[31m[INTO] func resolveObjectData  \033[0m\n")

    spewo := spew.ConfigState{ Indent: "    ", DisablePointerAddresses: true}

    fmt.Printf("\033[33m    [DUMP] selectionSet:  \033[0m\n")
    spewo.Dump(selectionSet)
    fmt.Printf("\033[33m    [DUMP] objectField:  \033[0m\n")
    spewo.Dump(objectField)
    fmt.Printf("\033[33m    [DUMP] resolvedData:  \033[0m\n")
    spewo.Dump(resolvedData)

    // go
    targetObjectFields := objectField.Type.(*Object).Fields
    selectionSetResult, _ := resolveSelectionSet(selectionSet, targetObjectFields, resolvedData)
    return selectionSetResult, nil
}

func getResolvedDataTargetFieldValue(resolvedData interface{}, targetFieldName string) (interface{}) {
    val := reflect.ValueOf(resolvedData)
    for i := 0; i < val.Type().NumField(); i++ {
        if val.Type().Field(i).Tag.Get("json") == targetFieldName {
            return reflect.Indirect(val).FieldByName(val.Type().Field(i).Name)
        }
    }
    return nil
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
    Name            string          `json:name`
    Description     string          `json:description`
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

    scalar.Name            = scalarTemplate.Name
    scalar.Description     = scalarTemplate.Description
    scalar.ResolveFunction = scalarTemplate.ResolveFunction

    // @todo: Scalar should provide serialize, parse value, parse literal function.
    
    return scalar
}


// scalar types
var Int = NewScalar(ScalarTemplate{
    Name: "Int",
    Description: "GraphQL Int type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).Int(), nil
    },
})

var String = NewScalar(ScalarTemplate{
    Name: "String",
    Description: "GraphQL String type",
    ResolveFunction: func (p ResolveParams) (interface{}, error) {
        return p.Context.(reflect.Value).String(), nil
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

type ResolveFunction func(p ResolveParams) (interface{}, error)

// resolve params for ResolveFunction()
type ResolveParams struct {

    // context from executor
    Context interface{}

    // arguments map from request
    Arguments map[string]interface{}
}


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