#include "lua_utils.h"
#include "raylib.h"
#include "lua/lualib.h"
#include "lua/lauxlib.h"
#include "display.h"
#include "editors/play_editor.h"
#include "editors/sprite_editor.h"
#include <string.h>

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

char* CheckLuaString(lua_State* L, const char* code) {
    if (!code) return strdup("No code provided");
    if (luaL_loadstring(L, code) != LUA_OK) {
        const char* err = lua_tostring(L, -1);
        char* result = strdup(err ? err : "unknown lua error :P womp womp");
        lua_pop(L, 1);
        return result;
    }

    lua_pop(L, 1);
    return NULL; // success! :D
}

void RunLuaScriptLoop(lua_State* L, const char* script) {
    const char* error = CheckLuaString(L, script);
    if (error) {
        isPlaying = false;
        return;
    }
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

static int l_sprite(lua_State* L) {
    int id = luaL_checkinteger(L, 1);
    int x = luaL_checkinteger(L, 2);
    int y = luaL_checkinteger(L, 3);

    struct Sprite* s = GetSpriteById(id);
    if (!s) return 0;

    for (size_t sy = 0; sy < s->height; sy++) {
        for (size_t sx = 0; sx < s->width; sx++) {
            int dx = x + (int)sx;
            int dy = y + (int)sy;
            if (dx < 0 || dy < 0 || dx >= WIDTH || dy >= HEIGHT) continue;

            Color c = s->data[sy][sx];
            if (c.a == 0) continue;

            framebuffer[dy][dx] = c;
        }
    }
    return 0;
}


lua_State* NewLuaState() {
    lua_State* L = luaL_newstate();
    luaL_openlibs(L);

    lua_newtable(L);
    lua_pushcfunction(L, pset);
    lua_setfield(L, -2, "pset");
    lua_pushcfunction(L, l_sprite); // gfx.sprite(id,x,y)
    lua_setfield(L, -2, "sprite");
    lua_setglobal(L, "gfx");

    return L;
}
