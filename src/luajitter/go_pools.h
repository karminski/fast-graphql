struct ObjectPool {
    void **entries;
    int count;
    int maxSize;
};
typedef struct ObjectPool ObjectPool;


void init_pools(lua_State *L);
void free_pools(lua_State *L);
extern lua_value *make_lua_value(lua_State *L);
void return_lua_value(lua_State *L, lua_value *value);
extern lua_table_entry *make_lua_table_entry(lua_State *L);
void return_lua_table_entry(lua_State *L, lua_table_entry *entry);
extern lua_unrolled_table *make_lua_unrolled_table(lua_State *L);
void return_lua_unrolled_table(lua_State *L, lua_unrolled_table *table);
