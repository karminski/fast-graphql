#include <stdlib.h>
#include <stdio.h>
#include <luajit.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <string.h>
#include <errno.h>

#define MT_GOCALLBACK "GO_CALLBACK"

#include "go_diag_memory.h"
#include "go_luaerrors.h"
#include "go_luatypes.h"
#include "go_pools.h"
#include "go_luainterface.h"

#include "go_callbacks.h"

extern lua_err *internal_dostring(lua_State *_L, char *script);
extern lua_State *new_luajit_state();
extern void close_lua(lua_State *_L);