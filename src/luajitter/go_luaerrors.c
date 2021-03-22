#include "go_luajit.h"

const lua_err INVALID_ERROR_str = {"INVALID ERROR"};
lua_err *INVALID_ERROR = (lua_err *)&INVALID_ERROR_str;

lua_err *get_lua_error(lua_State *_L, int errCode) {
    if (errCode == 0)
        return NULL;
    if (errCode == LUA_ERRMEM)
        return create_lua_error("LUA OUT OF MEMORY");

	const char *message = lua_tolstring(_L, -1, NULL);
	if (message == NULL)
		return INVALID_ERROR;

	lua_pop(_L, 1);
	return create_lua_error_from_luastr(message);
}

lua_err *create_lua_error_from_luastr(const char *msg) {
	lua_err *err = chmalloc(sizeof(lua_err));
	char *newMessage = chmalloc(sizeof(char)*(strlen(msg)+1));
	memcpy(newMessage, msg, strlen(msg)+1);
	err->message = newMessage;

	return err;
}

lua_err *create_lua_error(char *msg) {
	lua_err *err = chmalloc(sizeof(lua_err));
	err->message = msg;

	return err;
}

void free_lua_error(lua_err *err) {
    if (err == NULL)
        return;
	chfree(err->message);
	err->message = NULL;
    chfree(err);
}

int raise_lua_error(lua_State *_L, lua_err *err) {
    if (err == NULL)
        return 0;
    lua_pushstring(_L, err->message);
    free_lua_error(err);
    return lua_error(_L);
}
