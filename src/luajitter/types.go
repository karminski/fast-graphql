package luajitter

/*
#include "go_luajit.h"
*/
import "C"

import (
	"errors"
	"github.com/baohavan/go-pointer"
	"unsafe"
)

func outlyingAllocs() int {
	return int(C.outlying_allocs())
}

func clearAllocs() {
	C.clear_allocs()
}

var luaValueSize C.size_t = C.size_t(unsafe.Sizeof(C.lua_value{}))
var luaReturnSize C.size_t = C.size_t(unsafe.Sizeof(C.lua_return{}))

func fromGoValue(vm *LuaState, value interface{}, outValue *C.struct_lua_value) (cValue *C.struct_lua_value, err error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case uint64, uint32, int32, int64, int, uint, float64, float32:
		var castV float64
		switch innerV := v.(type) {
		case uint64:
			castV = float64(innerV)
		case uint32:
			castV = float64(innerV)
		case int32:
			castV = float64(innerV)
		case int64:
			castV = float64(innerV)
		case int:
			castV = float64(innerV)
		case uint:
			castV = float64(innerV)
		case float64:
			castV = innerV
		case float32:
			castV = float64(innerV)
		}
		if outValue == nil {
			outValue = C.make_lua_value(vm._l)
		}
		outValue.temporary = C._Bool(true)
		outValue.valueType = C.LUA_TNUMBER
		valData := (*C.double)(unsafe.Pointer(&outValue.data))
		*valData = C.double(castV)
	case bool:
		if outValue == nil {
			outValue = C.make_lua_value(vm._l)
		}

		outValue.temporary = C._Bool(true)
		outValue.valueType = C.LUA_TBOOLEAN
		valData := (*C._Bool)(unsafe.Pointer(&outValue.data))
		*valData = C._Bool(v)
	case string:
		if outValue == nil {
			outValue = C.make_lua_value(vm._l)
		}

		outValue.temporary = C._Bool(true)
		outValue.valueType = C.LUA_TSTRING
		valData := (**C.char)(unsafe.Pointer(&outValue.data))
		*valData = C.CString(v)
		C.increment_allocs()
		valDataArg := (*C.size_t)(unsafe.Pointer(&outValue.dataArg))
		*valDataArg = C.size_t(len(v))
	case *LocalLuaFunction, *LocalLuaData, *LocalLuaTable:
		castV := v.(LocalData)
		if outValue != nil {
			C.free_lua_value(vm._l, outValue)
		}
		if vm != castV.HomeVM() {
			return nil, errors.New("attempt to use local data in wrong VM")
		}
		outValue = castV.LuaValue()
	case func([]interface{}) ([]interface{}, error):
		if outValue == nil {
			outValue = C.make_lua_value(vm._l)
		}

		outValue.temporary = C._Bool(true)
		outValue.valueType = C.LUA_TUNLOADEDCALLBACK
		ptr := pointer.Save(v)

		valData := (*unsafe.Pointer)(unsafe.Pointer(&outValue.data))
		*valData = ptr
	case map[interface{}]interface{}:
		if outValue == nil {
			outValue = C.make_lua_value(vm._l)
		}

		outValue.temporary = C._Bool(true)
		outValue.valueType = C.LUA_TUNROLLEDTABLE
		valData := (*unsafe.Pointer)(unsafe.Pointer(&outValue.data))
		table := C.build_unrolled_table(vm._l, C.int(len(v)))
		*valData = unsafe.Pointer(table)

		entry := table.first

		for key, value := range v {
			entry.key, err = fromGoValue(vm, key, entry.key)
			if err != nil {
				C.free_temporary_lua_value(vm._l, outValue);
				return nil, err
			}

			entry.value, err = fromGoValue(vm, value, entry.value)
			if err != nil {
				C.free_temporary_lua_value(vm._l, outValue);
				return nil, err
			}

			entry = entry.next
		}
	default:
		return nil, errors.New("cannot marshal unknown type into lua")
	}

	return outValue, nil
}

func buildGoValue(vm *LuaState, value *C.struct_lua_value) interface{} {
	if value == nil {
		return nil
	}

	value.temporary = C._Bool(true)

	switch value.valueType {
	case C.LUA_TNUMBER:
		union := (*C.double)(unsafe.Pointer(&value.data))
		return float64(*union)
	case C.LUA_TBOOLEAN:
		union := (*C._Bool)(unsafe.Pointer(&value.data))
		return bool(*union == (C._Bool)(true))
	case C.LUA_TSTRING:
		union := (**C.char)(unsafe.Pointer(&(value.data)))
		return C.GoString(*union)
	case C.LUA_TTABLE:
		value.temporary = C._Bool(false)
		return &LocalLuaTable{
			LocalLuaData {
				value:  value,
				homeVM: vm,
			},
		}
	case C.LUA_TFUNCTION:
		isCFunction := (*C._Bool)(unsafe.Pointer(&value.dataArg))
		if *isCFunction == (C._Bool)(false) {
			value.temporary = C._Bool(false)
			return &LocalLuaFunction{
				LocalLuaData{
					value:  value,
					homeVM: vm,
				},
			}
		}

		fallthrough
	default:
		value.temporary = C._Bool(false)
		return &LocalLuaData{
			value:  value,
			homeVM: vm,
		}
	}
}

func buildGoValues(vm *LuaState, count int, values *[1 << 30]*C.struct_lua_value) []interface{} {
	goValues := make([]interface{}, count)
	for i := 0; i < count; i++ {
		goValues[i] = buildGoValue(vm, values[i])
	}

	return goValues
}
