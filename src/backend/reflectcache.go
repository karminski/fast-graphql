package backend

// A field represents a single field found in a struct.
type field struct {
	name          string
	nameb         []byte     
	stringifyFunc StringifyFunc
	
}


type structFields struct {
	list      []field
	nameIndex map[string]int
}

var fieldCache sync.Map // map[reflect.Type]structFields


func getCachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, getTypeFields(t))
	return f.(structFields)
}

func getTypeFields(t reflect.Type) structFields {

}


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

// main JIT resolve method
func steppingSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (string, error) {
	fmap := NewResolveFunctionMap()
    buildSchemaResolveFunctionMap(objectFields, fmap)


}