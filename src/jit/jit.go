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




