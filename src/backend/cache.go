package backend

import (
	"errors"
	"sync"
    "fast-graphql/src/frontend"
	"crypto/md5"
	"encoding/hex"
	"fmt"

)

var selectionSetCache sync.Map // map[selectionSetHash]cachedSelectionSet

type StringifyFunc func(s Stringifier, name string, value interface{})

type cachedSelectionSet struct {
	Name 		string
	Fields 		[]cachedField
}

type cachedField struct {
	Name 		    string
	Type 		    int
	Arguments 		map[string]interface{}
	StringifyFunc   StringifyFunc
	ResolveFunction ResolveFunction
}

func GetSelectionSetHash(queryHash [16]byte, name string) [16]byte {
	key := hex.EncodeToString(queryHash[:]) + name
	return md5.Sum([]byte(key))
}

func saveSelectionSet(k [16]byte, c cachedSelectionSet) {
	selectionSetCache.LoadOrStore(k, c)
}


func loadSelectionSet(k [16]byte) (cachedSelectionSet, bool)  {
	if c, ok := selectionSetCache.Load(k); ok {
		return c.(cachedSelectionSet), true
	}
	var c cachedSelectionSet 
	return c, false
}





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





func resolveCachedSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}, css cachedSelectionSet) (string, error) {
	for _, cf := range css.Fields {
		resolveCachedField(g, request, fieldName, cf, objectFields, resolvedData, css)
	}

    return "", nil
}

func resolveCachedField(g *GlobalVariables, request Request, fieldName string, cf cachedField, objectFields ObjectFields, resolvedData interface{}, css cachedSelectionSet) {
    var err error
	fmt.Printf("resolveCachedField cf: %s\n", cf.Name)
    
    // resolve
    if css.ResolveFunction != nil { // user defined resolve function avaliable
        resolvedData, err = cachedSchemaResolveFunction(g, request, fieldName, field, objectFields, resolvedData, cf)
        fmt.Printf("css.ResolveFunction != nil\n")
    } 

    cachedDefaultResolveFunction(g, request, objectFields[fieldName], resolvedData, cf)
}




func cachedDefaultResolveFunction(g *GlobalVariables, request Request, objectField *ObjectField, resolvedData interface{}, cf *cachedField) {
	switch cf.Type {
	case FIELD_TYPE_SCALAR:
		r0 := fastreflect.StructFieldByName(resolvedData, cf.Name)
		cf.StringifyFunc(g.Stringifier, cf.Name, r0)
	case FIELD_TYPE_LIST:
		
	case FIELD_TYPE_OBJECT:
		resolveCachedSelectionSet(g, request)
	}
}


func cachedSchemaResolveFunction(g *GlobalVariables, request Request, objectField *ObjectField, resolvedData interface{}, cf *cachedField) (interface{}, err) {
	// build resolve params for resolve function
    var resolveParams ResolveParams
    var err           error
    resolveParams.Source = resolvedData
    for arg, _ := cf.Arguments {
    	cf.Arguments[arg] = g.QueryVariablesMap[arg]
    }
    resolveParams.Arguments = cf.Arguments[arg]

    // resolve
    resolvedData, err = cf.ResolveFunction(resolveParams)
    return resolvedData, err
}