// executor.go
package backend

import (
    "fast-graphql/src/frontend"
    "fmt"
    "errors"
    "log"
)

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

func Execute(request Request) (*Result) {

    // process input
    AST := frontend.Compile(request.Query)

    // @todo: THE DOCUMENT NEED VALIDATE!
   
    // execute
    fmt.Printf("%v", request)
    if false {
        fmt.Printf("%v", AST)
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
    Name        string `json:name`
    Description string `json:description`
}

type Scalar struct {
    Name        string `json:name`
    Description string `json:description`
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
})

var String = NewScalar(ScalarTemplate{
    Name: "String",
    Description: "GraphQL String type",
})

// Object Syntax

type ObjectFields map[string]*ObjectField

type ObjectTemplate struct {
    Name string 
    Fields map[string]*ObjectField
}

type Object struct {
    Name string
    Fields map[string]*ObjectField
}

func (object *Object) GetName() string {
    return object.Name
}

type ObjectField struct {
    Name            string               `json:name`
    Type            FieldType            `json:type`
    Description     string               `json:description`
    Arguments       *Arguments            `json:arguments`    
    ResolveFunction FieldResolveFunction `json:"-"`
}

type Arguments map[string]*Argument

type Argument struct {
    Name string    `json:name` 
    Type FieldType `json:type`
}

type FieldResolveFunction func(i interface{}) (interface{}, error)


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