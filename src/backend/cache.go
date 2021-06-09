package backend

import (
    "fast-graphql/src/frontend"
    "fast-graphql/src/graphql"
	"errors"
	"sync"
	"crypto/md5"
	"encoding/hex"
	// "fmt"


    "github.com/karminski/fastreflect"
    // "github.com/davecgh/go-spew/spew"


)

var selectionSetCache sync.Map // map[selectionSetHash]cachedSelectionSet

type StringifyFunc func(s *Stringifier, value interface{})

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





func buildStringField(s *Stringifier, value interface{}) {
    s.buildString(value.(string))
}

func buildIntField(s *Stringifier, value interface{}) {
    s.buildInt(value.(int))
}

func buildFloat64Field(s *Stringifier, value interface{}) {
    s.buildFloat64(value.(float64))
}

func buildBoolField(s *Stringifier, value interface{}) {
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





func resolveCachedSelectionSet(g *GlobalVariables, request graphql.Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}, css cachedSelectionSet) (string, error) {
	// stringify
    g.Stringifier.buildObjectStart()


	var err 		   error

	// resolve field
	for i, cf := range css.Fields {
		// stringify
        g.Stringifier.buildFieldPrefix(cf.Name)

		if _, err = resolveCachedField(g, request, selectionSet, cf.Name, cf, objectFields, resolvedData, css); err != nil {
			return "", err
		} 
		// stringify
        if i < len(css.Fields) {
            g.Stringifier.buildComma() 
        }
	}

	// stringify
    g.Stringifier.buildObjectEnd()

    return "", nil
}

func resolveCachedField(g *GlobalVariables, request graphql.Request, selectionSet *frontend.SelectionSet, fieldName string, cf cachedField, objectFields ObjectFields, resolvedData interface{}, css cachedSelectionSet) (interface{}, error) {
    var err error
    // resolve by user defined function
    if cf.ResolveFunction != nil { // user defined resolve function avaliable
        if resolvedData, err = cachedSchemaResolveFunction(g, request, resolvedData, cf); err != nil {
        	return nil, err
        }
    } 
    // resolve by default function
    if resolvedData, err = cachedDefaultResolveFunction(g, request, selectionSet, objectFields[fieldName], resolvedData, cf); err != nil {
    	return nil, err
    }

    return resolvedData, nil
}




func cachedDefaultResolveFunction(g *GlobalVariables, request graphql.Request, selectionSet *frontend.SelectionSet, objectField *ObjectField, resolvedData interface{}, cf cachedField) (interface{}, error) {
	switch cf.Type {
	case FIELD_TYPE_SCALAR:
		r0 := fastreflect.StructFieldByName(resolvedData, cf.Name)
		cf.StringifyFunc(g.Stringifier, r0)
		return nil, nil
	case FIELD_TYPE_LIST:
    	allListElements := fastreflect.SliceAllElements(resolvedData)
	    targetObjectFields := objectField.Type.(*List).Payload.(*Object).Fields

    	// stringify
    	g.Stringifier.buildArrayStart()

    	// traverse list
    	stopPos := len(allListElements) - 1
    	for i, elements := range allListElements {
			cssHash := GetSelectionSetHash(g.queryHash, objectField.Name)
			if css, ok :=loadSelectionSet(cssHash); ok {
    			if _, err := resolveCachedSelectionSet(g, request, selectionSet, targetObjectFields, elements, css); err != nil {
    				return nil, err
    			}
    		    // stringify
    			if i < stopPos {
    			    g.Stringifier.buildComma() 
    			}
    		}
    	}
    	// stringify
    	g.Stringifier.buildArrayEnd()
    	return nil, nil
	case FIELD_TYPE_OBJECT:
		targetObjectFields := objectField.Type.(*Object).Fields

        r0 := fastreflect.StructFieldByName(resolvedData, objectField.Name)
    	if r0 != nil {
    	    resolvedData = r0
    	}
		cssHash := GetSelectionSetHash(g.queryHash, objectField.Name)
		if css, ok :=loadSelectionSet(cssHash); ok {
			return resolveCachedSelectionSet(g, request, selectionSet, targetObjectFields, resolvedData, css)
		}
	}
    return nil, errors.New("defaultResolveFunction(): can not resolve target field.")
}


func cachedSchemaResolveFunction(g *GlobalVariables, request graphql.Request, resolvedData interface{}, cf cachedField) (interface{}, error) {
	// build resolve params for resolve function
    var resolveParams ResolveParams
    var err           error
    resolveParams.Source = resolvedData
    resolveParams.Arguments = make(map[string]interface{}, len(cf.Arguments))
    for arg, _ := range cf.Arguments {
    	resolveParams.Arguments[arg] = g.QueryVariablesMap[arg]
    }

    // resolve
    resolvedData, err = cf.ResolveFunction(resolveParams)
    return resolvedData, err
}