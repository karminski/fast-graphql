package backend

import (
    "sync"

    "fast-graphql/src/luajitter"
    
    "strings"
    
    _ "net/http"
    _ "io"
    _ "net/http/pprof"

)

func closeVM(vm *luajitter.LuaState) {
    err := vm.Close()
    if err != nil {
        panic(err)
    }
}

func Emit(code string) {
    mutex.Lock()
    vm := luajitter.NewState()
    vm.SetGlobal("addFunc", AddValues)
    vm.DoString(code)
    closeVM(vm)
    defer mutex.Unlock()
}




func callSechemaResolveFunctionByJIT(fieldName string) {
    // build resolve params for resolve function
    var resolveParams ResolveParams
    var err           error
    resolveParams.Source = resolvedData
    if resolveParams.Arguments, err = getFieldArgumentsMap(g, field.Arguments); err != nil {
        return nil, err
    }

    // get resolve function
    resolveFunction := objectFields[fieldName].ResolveFunction

    // execute
    if resolvedData, err = resolveFunction(resolveParams); err != nil {
        return nil, err
    }

    // pass resolved data to [jit]
    

    
}


func readResult() string {
    r, err := vm.GetGlobal("_G.result")
    if err != nil {
        panic(err)
    }
    fmt.Println(r)
}

func buildSchemaResolveFunctionMap(objectFields ObjectFields) map[string]interface{} {
	var resolveFunctionMap map[string]interface{}
	for fieldName, objectField := range objectFields {
		if objectField.ResolveFunction != nil {
			resolveFunctionMap[fieldName] = objectField.ResolveFunction
		}
	} 
}


func steppingSelectionSet(g *GlobalVariables, request Request, selectionSet *frontend.SelectionSet, objectFields ObjectFields, resolvedData interface{}) (interface{}, error) {
    if selectionSet == nil {
        return nil, errors.New("jit.steppingSelectionSet(): empty selectionSet input.")
    }

    // resolve SelectionSet.Selections
    selections  := selectionSet.GetSelections()
    var err            error
    for i, selection := range selections {
        field     := selection.(*frontend.Field)
        fieldName := field.GetFieldNameString()

        // stringify
        g.Stringifier.buildFieldPrefix(fieldName)

        // resolve
        if resolvedResult, err = resolveField(g, request, fieldName, field, objectFields, resolvedData); err != nil {
            return nil, err
        }
        finalResult[fieldName] = resolvedResult  

    }

    // stringify
    g.Stringifier.buildObjectEnd()

    return finalResult, nil
}