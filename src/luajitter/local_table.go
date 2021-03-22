package luajitter

/*
#include "go_luajit.h"
*/
import "C"

import "unsafe"

type LocalLuaTable struct {
	LocalLuaData
}

func (table *LocalLuaTable) convertSingleUnrolledValue(value *C.struct_lua_value) (interface{}, error) {
	if value.valueType == C.LUA_TUNROLLEDTABLE {
		value.temporary = C._Bool(true)
		tablePtr := (**C.struct_lua_unrolled_table)(unsafe.Pointer(&value.data))
		return table.convertUnrolledTable(*tablePtr)
	}

	return buildGoValue(table.HomeVM(), value), nil
}

func (table *LocalLuaTable) convertUnrolledTable(unrolled *C.struct_lua_unrolled_table) (map[interface{}]interface{}, error) {
	retVal := make(map[interface{}]interface{}, unrolled.arraySize+unrolled.hashSize)
	var entry *C.struct_lua_table_entry
	entry = unrolled.first

	for entry != nil {
		key, err := table.convertSingleUnrolledValue(entry.key)
		if err != nil {
			C.free_lua_value(table.HomeVM()._l, entry.key)
			return nil, err
		}

		value, err := table.convertSingleUnrolledValue(entry.value)
		if err != nil {
			C.free_lua_value(table.HomeVM()._l, entry.key)
			C.free_lua_value(table.HomeVM()._l, entry.value)
			return nil, err
		}

		retVal[key] = value

		entry = entry.next
	}

	return retVal, nil
}

func (table *LocalLuaTable) Unroll() (map[interface{}]interface{}, error) {
	result := C.unroll_table(table.HomeVM()._l, table.LuaValue())
	if result.err != nil {
		defer C.free_lua_error(result.err)
		return nil, LuaErrorToGo(result.err)
	} else if result.value != nil {
		result.value.temporary = C._Bool(true)
		unrollTablePtr := (**C.struct_lua_unrolled_table)(unsafe.Pointer(&result.value.data))
		unrollTable := *unrollTablePtr
		retVal, err := table.convertUnrolledTable(unrollTable)
		if err != nil {
			C.free_lua_value(table.HomeVM()._l, result.value)
			return nil, err
		} else {
			C.free_temporary_lua_value(table.HomeVM()._l, result.value)
			return retVal, nil
		}
	}
	return nil, nil
}
