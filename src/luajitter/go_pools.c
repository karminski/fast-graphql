#include "go_luajit.h"

static const char *EntryPoolKey = "internal_entry_pool";
static const char *ValuePoolKey = "internal_value_pool";
static const char *TablePoolKey = "internal_table_pool";

void init_single_pool(ObjectPool *pool, int newSize) {
    void **newEntries = chmalloc(sizeof(void*)*newSize);
    for (int i = 0; i < pool->count; i++) {
        newEntries[i] = pool->entries[i];
    }
    if (pool->entries != NULL) {
        chfree(pool->entries);
    }
    pool->entries = newEntries;
    pool->maxSize = newSize;
}

void init_pools(lua_State *L) {
    ObjectPool *entryPool = chmalloc(sizeof(ObjectPool));
    entryPool->entries = NULL;
    entryPool->count = 0;
    init_single_pool(entryPool, 10000);

    ObjectPool *valuePool = chmalloc(sizeof(ObjectPool));
    valuePool->entries = NULL;
    valuePool->count = 0;
    init_single_pool(valuePool, 10000);

    ObjectPool *tablePool = chmalloc(sizeof(ObjectPool));
    tablePool->entries = NULL;
    tablePool->count = 0;
    init_single_pool(tablePool, 1000);

    lua_pushlightuserdata(L, (void *)&EntryPoolKey);
    lua_pushlightuserdata(L, (void*)entryPool);
    lua_settable(L, LUA_REGISTRYINDEX);

    lua_pushlightuserdata(L, (void *)&ValuePoolKey);
    lua_pushlightuserdata(L, (void*)valuePool);
    lua_settable(L, LUA_REGISTRYINDEX);

    lua_pushlightuserdata(L, (void *)&TablePoolKey);
    lua_pushlightuserdata(L, (void*)tablePool);
    lua_settable(L, LUA_REGISTRYINDEX);
}

ObjectPool *get_obj_pool(lua_State *L, const char **key) {
  lua_pushlightuserdata(L, (void *)key);
  lua_gettable(L, LUA_REGISTRYINDEX);
  ObjectPool *pool = (ObjectPool *)lua_touserdata(L, -1);
  lua_pop(L, 1);
  return pool;
}

void free_pools(lua_State *L) {
    ObjectPool *entries = get_obj_pool(L, &EntryPoolKey);
    for (int i = 0; i < entries->count; i++) {
        chfree(entries->entries[i]);
    }
    chfree(entries->entries);
    chfree(entries);

    ObjectPool *values = get_obj_pool(L, &ValuePoolKey);
    for (int i = 0; i < values->count; i++) {
        chfree(values->entries[i]);
    }
    chfree(values->entries);
    chfree(values);

    ObjectPool *tables = get_obj_pool(L, &TablePoolKey);
    for (int i = 0; i < tables->count; i++) {
        chfree(tables->entries[i]);
    }
    chfree(tables->entries);
    chfree(tables);

    lua_pushlightuserdata(L, (void *)&EntryPoolKey);
    lua_pushnil(L);
    lua_settable(L, LUA_REGISTRYINDEX);

    lua_pushlightuserdata(L, (void *)&ValuePoolKey);
    lua_pushnil(L);
    lua_settable(L, LUA_REGISTRYINDEX);

    lua_pushlightuserdata(L, (void *)&TablePoolKey);
    lua_pushnil(L);
    lua_settable(L, LUA_REGISTRYINDEX);
}

void *get_from_pool(ObjectPool *pool) {
    if(pool->count == 0) {
        return NULL;
    }
    pool->count--;
    void *item = pool->entries[pool->count];
    pool->entries[pool->count] = NULL;
    return item;
}

void add_to_pool(ObjectPool *pool, void *item) {
    while (pool->count+1 >= pool->maxSize) {
        init_single_pool(pool, pool->maxSize*2);
    }
    pool->entries[pool->count] = item;
    pool->count++;
}

lua_unrolled_table *make_lua_unrolled_table(lua_State *L) {
    ObjectPool *pool = get_obj_pool(L, &TablePoolKey);
    lua_unrolled_table *table = (lua_unrolled_table*)get_from_pool(pool);
    if (table != NULL)
        return table;
    return chmalloc(sizeof(lua_unrolled_table));
}

void return_lua_unrolled_table(lua_State *L, lua_unrolled_table *table) {
    ObjectPool *pool = get_obj_pool(L, &TablePoolKey);
    add_to_pool(pool, (void*)table);
}

lua_value *make_lua_value(lua_State *L) {
    ObjectPool *pool = get_obj_pool(L, &ValuePoolKey);
    lua_value *value = (lua_value*)get_from_pool(pool);
    if (value != NULL)
        return value;
    return chmalloc(sizeof(lua_value));
}

void return_lua_value(lua_State *L, lua_value *value) {
    ObjectPool *pool = get_obj_pool(L, &ValuePoolKey);
    add_to_pool(pool, (void*)value);
}

lua_table_entry *make_lua_table_entry(lua_State *L) {
    ObjectPool *pool = get_obj_pool(L, &EntryPoolKey);
    lua_table_entry *entry = (lua_table_entry*)get_from_pool(pool);
    if (entry != NULL)
        return entry;
    return chmalloc(sizeof(lua_table_entry));
}

void return_lua_table_entry(lua_State *L, lua_table_entry *entry) {
    ObjectPool *pool = get_obj_pool(L, &EntryPoolKey);
    add_to_pool(pool, (void*)entry);
}
