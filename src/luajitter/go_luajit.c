#include "go_luajit.h"

lua_err *internal_dostring(lua_State *_L, char *script) {
	int retVal = luaL_dostring(_L, script);
	return get_lua_error(_L, retVal);
}

lua_State *new_luajit_state() {
	lua_State *_L = luaL_newstate();
	luaL_openlibs(_L);

	luaL_newmetatable(_L, MT_GOCALLBACK);
	lua_pushliteral(_L,"__call");
	lua_pushcfunction(_L,&execute_go_callback);
	lua_settable(_L,-3);

	lua_pushliteral(_L,"__gc");
	lua_pushcfunction(_L,&release_cgo_handle);
	lua_settable(_L,-3);
	lua_pop(_L,1);

	init_pools(_L);

	return _L;
}

void close_lua(lua_State *_L) {
    free_pools(_L);
    lua_close(_L);
}