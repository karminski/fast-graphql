package jit

import (
    "sync"

    "luajitter-demo/src/luajitter"
    
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