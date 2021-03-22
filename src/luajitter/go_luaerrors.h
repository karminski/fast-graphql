struct lua_err {
	char *message;
};
typedef struct lua_err lua_err;

extern lua_err *INVALID_ERROR;
extern lua_err *get_lua_error(lua_State *_L, int errCode);
extern lua_err *create_lua_error_from_luastr(const char *msg);
extern lua_err *create_lua_error(char *msg);
extern void free_lua_error(lua_err *err);
extern int raise_lua_error(lua_State *_L, lua_err *err);