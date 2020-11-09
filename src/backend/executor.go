// executor.go
package backend

import (
	"reflect"
)

type Result struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func Execute(ast *Document, fields *Fields) (*Result) {

	// @todo: THE DOCUMENT NEED VALIDATE!
	
	// execute resolve function
 	for fieldName, field := range ast {

 	}
	
}

// types

type ObjectTmpl struct {
	Name string 
	Fields map[string]*ObjectField
}

type ObjectField struct {
	Name String
	Type String

}

type Object struct {

}

func NewObject (objectTmpl ObjectTmpl) (*Object, error) {
	object := &Object{}

	// check object input
	if objectTmpl.Name == "" {
		var err error
		err = "ObjectTmpl.Name is not defined"
		return nil err
	}
	

}