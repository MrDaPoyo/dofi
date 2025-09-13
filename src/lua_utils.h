#ifndef DOFI_LUA_H
#define DOFI_LUA_H

#include "lua/lua_utils.h"

lua_State *NewLuaState();
void AnnihilateLuaState(lua_State *L);
void RunLuaScript(lua_State *L, const char *script);

#endif // DOFI_LUA_H