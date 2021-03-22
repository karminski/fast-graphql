package luajitter

/*
#include "go_luajit.h"
*/
import "C"

type LocalData interface {
	LuaValue() *C.struct_lua_value
	HomeVM() *LuaState
	Close() error
}

type LocalLuaData struct {
	value  *C.struct_lua_value
	homeVM *LuaState
}

func (d *LocalLuaData) LuaValue() *C.struct_lua_value {
	return d.value
}

func (d *LocalLuaData) HomeVM() *LuaState {
	return d.homeVM
}

func (d *LocalLuaData) Close() error {
	if d.value != nil {
		C.free_lua_value(d.homeVM._l, d.value)
		d.value = nil
	}

	return nil
}
