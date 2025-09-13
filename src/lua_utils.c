#include "lua_utils.h"
#include "raylib.h"
#include "lua/lualib.h"
#include "lua/lauxlib.h"
#include "display.h"

#define FPS 60

Color framebuffer[WIDTH][HEIGHT];

void CallLuaFunction(lua_State* L, const char* funcName) {
    lua_getglobal(L, funcName);
    if (lua_isfunction(L, -1)) {
        if (lua_pcall(L, 0, 0, 0) != LUA_OK) {
            printf("Error running %s: %s\n", funcName, lua_tostring(L, -1));
            lua_pop(L, 1);
        }
    } else {
        lua_pop(L, 1);
    }
}

void RunLuaScriptLoop(lua_State* L, const char* script) {
    if (luaL_loadstring(L, script) != LUA_OK) {
        printf("Failed to load Lua script: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        return;
    }

    if (lua_pcall(L, 0, 0, 0) != LUA_OK) {
        printf("Error running Lua script: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        return;
    }

    CallLuaFunction(L, "init");

    SetTargetFPS(FPS);
    CallLuaFunction(L, "update");

    BeginDrawing();
    ClearBackground(BLACK);

    CallLuaFunction(L, "draw");

    for (int x = 0; x < WIDTH; x++) {
        for (int y = 0; y < HEIGHT; y++) {
            DrawPixel(x, y, framebuffer[x][y]);
        }
    }

    EndDrawing();
}

void AnnihilateLuaState(lua_State* L) {
    lua_close(L);
}

int pset(lua_State* L) {
    int x = luaL_checkinteger(L, 1);
    int y = luaL_checkinteger(L, 2);
    int r = luaL_checkinteger(L, 3);
    int g = luaL_checkinteger(L, 4);
    int b = luaL_checkinteger(L, 5);

    if (x >= 0 && x < WIDTH && y >= 0 && y < HEIGHT &&
        r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255) {
        framebuffer[x][y] = (Color){ r, g, b, 255 };
        return 0;
    }

    return 0;
}

lua_State* NewLuaState() {
    lua_State* L = luaL_newstate();
    luaL_openlibs(L);

    lua_newtable(L);
    lua_pushcfunction(L, pset);
    lua_setfield(L, -2, "pset");
    lua_setglobal(L, "gfx");

    return L;
}
