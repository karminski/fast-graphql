package backend

import (
	"errors"
    "fast-graphql/src/frontend"

)



type StringifyFunc func(name string, value interface{})



func buildStringField(s Stringifier, field string, value interface{}) {
    s.buildString(field)
    s.buildColon()
    s.buildString(value.(string))
}

func buildIntField(s Stringifier, field string, value interface{}) {
    s.buildString(field)
    s.buildColon()
    s.buildInt(value.(int))
}

func buildFloat64Field(s Stringifier, field string, value interface{}) {
    s.buildString(field)
    s.buildColon()
    s.buildFloat64(value.(float64))
}

func buildBoolField(s Stringifier, field string, value interface{}) {
    s.buildString(field)
    s.buildColon()
    s.buildBool(value.(bool))
}


func NewResolveFunctionMap() map[string]interface{} {
	var fmap map[string]interface{}
	fmap = make(map[string]interface{})
	return fmap
}


func buildSchemaResolveFunctionMap(objectFields ObjectFields, fmap map[string]interface{})  {	
	for fieldName, objectField := range objectFields {
		if objectField.ResolveFunction != nil {
			fmap[fieldName] = objectField.ResolveFunction
		}
		if subFields := getSubObjectFields(objectField); subFields != nil {
			buildSchemaResolveFunctionMap(subFields, fmap)
		}
	} 
}

// main JIT resolve method
func steppingSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (string, error) {
	fmap := NewResolveFunctionMap()
    buildSchemaResolveFunctionMap(objectFields, fmap)

    return "", nil

}


func callResolveFuncByArguments(args []interface{}, fmap map[string]interface{}) (interface{}, error) {
	var field string
	var p ResolveParams


	// call resolve function
	if _, ok := fmap[field]; !ok {
		return nil, errors.New("callResolveFuncByArguments(): unknown field input.")
	}
	return fmap[field].(ResolveFunction)(p)
	
} 

func callResolveFuncBySource(args []interface{}, fmap map[string]interface{}) (interface{}, error) {
	var field string
	var p ResolveParams


	// call resolve function
	if _, ok := fmap[field]; !ok {
		return nil, errors.New("callResolveFuncBySource(): unknown field input.")
	}
	return fmap[field].(ResolveFunction)(p)
} 