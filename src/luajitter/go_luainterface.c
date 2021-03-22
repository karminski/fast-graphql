#include "go_luajit.h"

lua_return call_function(lua_State *_L, lua_value *func, lua_args args) {
    lua_return retVal = {};
    int startTop = lua_gettop(_L);
    lua_err *err = push_lua_value(_L, func);
    if (err != NULL) {
        retVal.err = err;
        return retVal;
    }
    err = push_lua_args(_L, args);
    if (err != NULL) {
        retVal.err = err;
        return retVal;
    }
    int resultCode = lua_pcall(_L, args.valueCount, LUA_MULTRET, 0);
    retVal.err = get_lua_error(_L, resultCode);
    
    if (retVal.err == NULL) {
        int endTop = lua_gettop(_L);
        int popValues = endTop - startTop;
        if (popValues == 0)
            return retVal;
        
        return pop_lua_values(_L, popValues);
    }

    return retVal;
}

typedef lua_result (*last_segment_handler)(lua_State *_L, int depth, const char *fullPath, const char *segment);

lua_err *create_walk_error(const char *fullPath, const char *path, const char *error) {
    int segLen;
    for (segLen = 0; path[segLen] != '\0' && path[segLen] != '.'; segLen++) {}

    const char *format = "Failed path walk in '%s' on '%.*s': %s";
    int length = strlen(format) - 8 + strlen(fullPath) + segLen + strlen(error) ;
    char *fullError = chmalloc(sizeof(char)*(length+1));
    snprintf(fullError,length,format,fullPath,segLen,path,error);
    return create_lua_error(fullError);
}

void push_walk_key(lua_State *_L, const char *path, int segLen) {
    char *outPath;
    char *expectedFinishPath = (char*)(path+segLen);

    int result = strtol(path, &outPath, 10);
    if ((result != 0 || errno == 0) && expectedFinishPath == outPath) {
        lua_pushinteger(_L, result);
        return;
    }

    lua_pushlstring(_L, path, (size_t)segLen);
}

lua_result walk_next_segment(lua_State *_L, int depth, const char *fullPath, const char *path, last_segment_handler handler, _Bool fillIntermediateTables) {
    lua_result retVal = {};
    if (lua_isnil(_L, -1)) {
        retVal.err = create_walk_error(fullPath, path, "nil segment");
        return retVal;
    }

    if (!lua_istable(_L, -1)) {
        retVal.err = create_walk_error(fullPath, path, "not table");
        return retVal;
    }

    int segLen;
    for (segLen = 0; path[segLen] != '\0' && path[segLen] != '.'; segLen++) {}
    
    if (segLen == 0) {
        retVal.err = create_walk_error(fullPath, path, "walked path segment zero length");
        return retVal;
    }

    if (path[segLen] == '.') {
        push_walk_key(_L, path, segLen);
        lua_gettable(_L, -2);
        if (fillIntermediateTables && lua_isnil(_L, -1)) {
            //Remove nil from the stack
            lua_pop(_L, 1);

            //Push key & new table
            push_walk_key(_L, path, segLen);
            lua_newtable(_L);
            //Put new table into old one
            lua_settable(_L, -3);

            //Retry get now that the new table is in there
            push_walk_key(_L, path, segLen);
            lua_gettable(_L, -2);
        }
        retVal = walk_next_segment(_L, depth+1, fullPath, path+segLen+1, handler, fillIntermediateTables);
        lua_pop(_L, 1);
        return retVal;
    }

    return handler(_L, depth, fullPath, path);
}

lua_result walk_table_path(lua_State *_L, int valueIndex, const char *path, last_segment_handler handler, _Bool fillIntermediateTables) {
    lua_pushvalue(_L, valueIndex);
    lua_result retVal = walk_next_segment(_L, 1, path, path, handler, fillIntermediateTables);
    lua_pop(_L, 1);
    return retVal;
}

lua_result get_global_handler(lua_State *_L, int depth, const char *fullPath, const char *path) {
    push_walk_key(_L, path, strlen(path));
    lua_gettable(_L, -2);
    return convert_stack_value(_L);
}

lua_result get_global(lua_State *_L, const char *path, _Bool fillIntermediateTables) {
    lua_result retVal = walk_table_path(_L, LUA_GLOBALSINDEX, path, get_global_handler, fillIntermediateTables);
    return retVal;
}

lua_result set_global_handler(lua_State *_L, int depth, const char *fullPath, const char *path) {
    push_walk_key(_L, path, strlen(path));
    lua_pushvalue(_L, -2-depth);
    lua_settable(_L, -3);
    lua_result res = {};
    return res;
}

lua_err *set_global(lua_State *_L, const char *path, lua_value *value, _Bool fillIntermediateTables) {
    lua_err *err = push_lua_value(_L, value);
    if (err != NULL)
        return err;
    lua_result result = walk_table_path(_L, LUA_GLOBALSINDEX, path, set_global_handler, fillIntermediateTables);
    lua_pop(_L, 1);
    if (result.value) {
        free_lua_value(_L, result.value);
    }
    return result.err;
}
