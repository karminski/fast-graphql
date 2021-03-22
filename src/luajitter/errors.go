package luajitter

/*
#include "go_luajit.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func LuaErrorToGo(err *C.lua_err) error {
	if err == nil {
		return nil
	}
	if err == C.INVALID_ERROR {
		panic("INVALID ERROR RAISED FROM LUA")
	}
	outErr := errors.New(C.GoString(err.message))
	return outErr
}

var luaErrSize C.size_t = C.size_t(unsafe.Sizeof(C.lua_err{}))

func GoErrorToLua(err error) *C.lua_err {
	if err == nil {
		return nil
	}

	outErr := (*C.struct_lua_err)(C.chmalloc(luaErrSize))
	outErr.message = C.CString(err.Error())
	C.increment_allocs()
	return outErr
}
