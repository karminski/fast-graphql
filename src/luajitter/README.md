# LuaJitter

Blazing fast LuaJIT bindings with great ergonomics.  Uses Go 1.14.

## Installing

### MacOS

1. Install luajit with Homebrew: `brew install luajit`
1. Locate the pkg-config file (it is usually in or around `/usr/local/Cellar/luajit/2.0.5/lib/pkgconfig/luajit.pc`)
1. Modify the LDFLAGS at the bottom of the luajit.pc file to remove -pagezero_size and -image_base arguments
1. `go test ./...` should now work

### Linux

1. Download the LuaJIT source from http://luajit.org/download.html
1. Follow the POSIX installation instructions at http://luajit.org/install.html
1. Locate the pkg-config file
1. Modify the LDFLAGS at the bottom of the luajit.pc file to remove -pagezero_size and -image_base arguments
1. `go test ./...` should now work

### Windows

1. Install MinGW from http://www.mingw.org/wiki/Install_MinGW (you can skip this step if you have Git Bash installed, it can be your MinGW terminal)
1. Install the Twilight Dragon GCC compiler from http://tdm-gcc.tdragon.net/download
1. Download the LuaJIT source from http://luajit.org/download.html
1. Run your MinGW terminal as administrator and navigate to the LuaJIT source on your hard drive
1. Run `mingw32-make`, then navigate to the `src` subdirectory and locate `lua51.dll`
1. Copy `lua51.dll` to whatever folder you intend to use LuaJitter from
1. `go test ./...` should now work

## Usage

```go
package main 

import (
    "fmt"
    "github.com/cannibalvox/luajitter"
)

func AddValues(args []interface{}) ([]interface{}, error) {
    lValue := args[0].(float64)
    rValue := args[1].(float64)
    return []interface{}{lValue + rValue}, nil
}

func closeVM(vm *luajitter.LuaState) {
    err := vm.Close()
    if err != nil {
        panic(err)
    }
}

func main() {
    //Create/destroy states
    vm := luajitter.NewState()
    defer closeVM(vm)

    //Execute arbitrary code
    err := vm.DoString(`
    print("Hello, World")
    someGlobal = {
        itsANumber = 5,
        itsAString = "WOW",
        itsAFunction = function(l,r) return l+r end,
    }
 `)
    if err != nil {
        panic(err)
    }
 
    //Access globals with dot-separated paths
    num, err := vm.GetGlobal("someGlobal.itsANumber")
    if err != nil {
        panic(err)
    }
    fmt.Println(num)
 
    str, err := vm.GetGlobal("someGlobal.itsAString")
    if err != nil {
        panic(err)
    }
    fmt.Println(str)
 
    //Call Lua functions from go
    fObj, err := vm.GetGlobal("someGlobal.itsAFunction")
    if err != nil {
        panic(err)
    }
 
    f := fObj.(*luajitter.LocalLuaFunction)
    sumRet, err := f.Call(6, 3)
    if err != nil {
        panic(err)
    }
    fmt.Println(sumRet[0])
 
    //Manipulate lua globals from go
    err = vm.SetGlobal("someGlobal.itsAString", "NEW STRING")
    if err != nil {
        panic(err)
    }
 
    str, err = vm.GetGlobal("someGlobal.itsAString")
    if err != nil {
        panic(err)
    }
    fmt.Println(str)
 
    //Set global + initialize intervening tables
    err = vm.InitGlobal("someGlobal.subGlobal.subGlobal.value", true)
    if err != nil {
        panic(err)
    }
 
    b, err := vm.GetGlobal("someGlobal.subGlobal.subGlobal.value")
    if err != nil {
        panic(err)
    }
    fmt.Println(b)
 
    //Catch lua errors from go
    err = vm.DoString(`
    someGlobal.errorFunc = function() error("lua failed!") end
 `)
    if err != nil {
        panic(err)
    }
 
    errFObj, err := vm.GetGlobal("someGlobal.errorFunc")
    if err != nil {
        panic(err)
    }
 
    errF := errFObj.(*luajitter.LocalLuaFunction)
    _, err = errF.Call()
    fmt.Println(err.Error())

    //Call go functions from lua
    err = vm.SetGlobal("addFunc", AddValues)
    if err != nil {
        panic(err)
    }
 
    err = vm.DoString(`
    print(addFunc(1,7))
 `)
    if err != nil {
        panic(err)
    }
}
```