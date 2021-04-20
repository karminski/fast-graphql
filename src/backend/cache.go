package backend

import (
	"errors"
	"sync"
    "fast-graphql/src/frontend"
	"crypto/md5"
	"encoding/hex"

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





func resolveCachedSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (string, error) {
	fmap := NewResolveFunctionMap()
    buildSchemaResolveFunctionMap(objectFields, fmap)
    return "", nil
}


func resolveFieldByCache() {

}

func resolveFieldBySchema() {

}


