extern lua_return call_function(lua_State *_L, lua_value *func, lua_args args);
extern lua_result get_global(lua_State *_L, const char *path, _Bool fillIntermediateTables);
extern lua_err *set_global(lua_State *_L, const char *path, lua_value *value, _Bool fillIntermediateTables);
