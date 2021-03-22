package luajitter

/*
#include "go_luajit.h"
*/
import "C"

import "unsafe"

type LocalLuaFunction struct {
	LocalLuaData
}

func (f *LocalLuaFunction) Call(args ...interface{}) ([]interface{}, error) {
	luaArgs := C.lua_args{
		valueCount: C.int(len(args)),
		values:     nil,
	}

	argsIn := make([]*C.struct_lua_value, len(args))
	for i := 0; i < len(args); i++ {
		argsIn[i] = nil
	}

	defer C.free_temporary_lua_value_array(f.HomeVM()._l, *(***C.struct_lua_value)(unsafe.Pointer(&argsIn)),C.int( len(args)))

	var err error
	for ind, arg := range args {
		val, err := fromGoValue(f.HomeVM(), arg, nil)
		if err != nil {
			return nil, err
		}

		argsIn[ind] = val
	}

	if len(argsIn) > 0 {
		luaArgs.values = &argsIn[0]
	}

	var allRetVals []interface{}
	retVal := C.call_function(f.HomeVM()._l, f.LuaValue(), luaArgs)
	if retVal.err != nil {
		defer C.free_lua_error(retVal.err)
		err = LuaErrorToGo(retVal.err)
	} else if retVal.valueCount > 0 {
		defer C.free_temporary_lua_return(f.HomeVM()._l, retVal, C._Bool(true))
		valueList := (*[1 << 30]*C.struct_lua_value)(unsafe.Pointer(retVal.values))
		allRetVals = buildGoValues(f.HomeVM(), int(retVal.valueCount), valueList)
	}

	return allRetVals, err
}
