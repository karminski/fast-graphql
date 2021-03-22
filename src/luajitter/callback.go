package luajitter

/*
#include "go_luajit.h"
*/
import "C"
import (
	"github.com/baohavan/go-pointer"
	"unsafe"
)

//export releaseCGOHandle
func releaseCGOHandle(handle unsafe.Pointer) *C.lua_err {
	pointer.Unref(handle)
	return nil
}

//export callbackGoFunction
func callbackGoFunction(_L *C.lua_State, handle unsafe.Pointer, args C.lua_args, ret *C.lua_return) {
	handlePtr := pointer.Restore(handle)
	goFunction, ok := handlePtr.(func([]interface{}) ([]interface{}, error))
	if !ok {
		ret.err.message = C.CString("attempted to call go function with non-callback object")
		C.increment_allocs()
		return
	}

	state := vmMap[_L]
	argCount := int(args.valueCount)
	argsList := (*[1 << 30]*C.struct_lua_value)(unsafe.Pointer(args.values))
	goArgs := buildGoValues(state, argCount, argsList)

	retVals, err := goFunction(goArgs)
	if err != nil {
		ret.err.message = C.CString(err.Error())
		C.increment_allocs()
		return
	}

	//We don't need to alloc values for local vals or for nil, so figure out
	//how many we need C to generate
	createValues := 0
	for _, val := range retVals {
		if val == nil {
			continue
		}

		switch val.(type) {
		case *LocalLuaData, *LocalLuaFunction:
			continue
		}

		createValues++
	}

	ret.err = nil
	ret.valueCount = C.int(len(retVals))
	ret.values = C.build_values(_L, ret.valueCount, C.int(createValues))
	allValues := (*[1 << 30]*C.struct_lua_value)(unsafe.Pointer(ret.values))

	//We have N allocated values for M slots, but we need to make sure the right slots are populated
	nextFilledValue := createValues - 1
	for i := len(retVals) - 1; i >= 0; i-- {
		if retVals[i] == nil {
			continue
		}

		switch retVals[i].(type) {
		case *LocalLuaData, *LocalLuaFunction:
			continue
		}

		val := allValues[nextFilledValue]
		allValues[nextFilledValue] = nil
		allValues[i] = val
		nextFilledValue--
	}

	for idx, singleVal := range retVals {
		var value *C.lua_value
		value, err = fromGoValue(state, singleVal, allValues[idx])
		allValues[idx] = value
		if err != nil {
			break
		}
	}
}
