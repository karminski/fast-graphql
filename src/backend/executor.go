// executor.go
package backend

import (
    "reflect"
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
   
   
    
}

// types

type Type interface {
    Name() string
}

type FieldType interface {
    Name() string
}

// List types

type List struct {
    Payload Type `json:payload`
}

func (list *List) Name() string {
    return fmt.Sprintf("%v", list.Payload)
}

func NewList(i Type) (*List, error) {
    list := &List{}

    if i == nil {
        var err error
        err = "NewList() input is nil"
        return nil, err
    }

    list.Payload = i
    return list, nil
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

func NewScalar(scalarTemplate ScalarTemplate) *Scalar {
    scalar := &Scalar{}

    // check scalar template
    if scalarTemplate.Name == nil {
        var err error
        err = "scalarTemplate.Name is not defined"
        return nil, err
    }

    Scalar.Name        = scalarTemplate.Name
    Scalar.Description = scalarTemplate.Description

    // @todo: Scalar should provide serialize, parse value, parse literal function.
    
    return Scalar
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

type ObjectField struct {
    Name            string               `json:name`
    Type            FieldType            `json:type`
    Description     string               `json:description`
    Arguments       Arguments            `json:arguments`    
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
        var err error
        err = "ObjectTemplate.Name is not defined"
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
        var err error
        err = "SchemaTemplate.Query is not defined"
        return nil, err
    }

    // fill schema
    schema.Query = schemaTemplate.Query


    return schema, nil
}