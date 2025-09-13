#include "editors.h"
#include "text_editor.h"
#include "../display.h"
#include "../lua_utils.h"

#include "../lua/lua.h"
#include "../lua/lauxlib.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static lua_State* gLuaVM = NULL;
static bool gScriptLoaded = false;

extern Color framebuffer[WIDTH][HEIGHT];

static void CallLuaFunctionSafe(lua_State* L, const char* funcName) {
    lua_getglobal(L, funcName);
    if (lua_isfunction(L, -1)) {
        if (lua_pcall(L, 0, 0, 0) != LUA_OK) {
            printf("[lua] error running %s: %s\n", funcName, lua_tostring(L, -1));
            lua_pop(L, 1);
        }
    } else {
        lua_pop(L, 1);
    }
}

static bool LoadEntireFile(const char* path, char** outBuf) {
    *outBuf = NULL;
    FILE* f = fopen(path, "rb");
    if (!f) return false;
    fseek(f, 0, SEEK_END);
    long sz = ftell(f);
    if (sz < 0) { fclose(f); return false; }
    rewind(f);
    char* buf = (char*)malloc((size_t)sz + 1);
    if (!buf) { fclose(f); return false; }
    size_t rd = fread(buf, 1, (size_t)sz, f);
    fclose(f);
    buf[rd] = '\0';
    *outBuf = buf;
    return true;
}

static void EnsureScriptLoaded(void) {
    if (gScriptLoaded) return;

    if (!gLuaVM) {
        gLuaVM = NewLuaState();
    }

    char* script = NULL;
    if (!LoadEntireFile("assets/examples/gradient.lua", &script)) {
        fprintf(stderr, "[play] Failed to load script file.\n");
        return;
    }

    if (luaL_loadstring(gLuaVM, script) != LUA_OK) {
        fprintf(stderr, "[lua] load error: %s\n", lua_tostring(gLuaVM, -1));
        lua_pop(gLuaVM, 1);
        free(script);
        return;
    }

    if (lua_pcall(gLuaVM, 0, 0, 0) != LUA_OK) {
        fprintf(stderr, "[lua] runtime error: %s\n", lua_tostring(gLuaVM, -1));
        lua_pop(gLuaVM, 1);
        free(script);
        return;
    }

    free(script);

    CallLuaFunctionSafe(gLuaVM, "init");
    gScriptLoaded = true;
}

static void UpdateLua(void) {
    if (!gScriptLoaded) return;
    CallLuaFunctionSafe(gLuaVM, "update");
}

static void DrawLua(void) {
    if (!gScriptLoaded) return;
    CallLuaFunctionSafe(gLuaVM, "draw");

    for (int x = 0; x < WIDTH; x++) {
        for (int y = 0; y < HEIGHT; y++) {
            Color c = framebuffer[x][y];
            if (c.a != 0) {
                DrawPixel(x, y, c);
            }
        }
    }
}

void RenderPlayEditor(void) {
    if (displayNavbar) {
        HideNavbar();
    }

    EnsureScriptLoaded();

    ClearBackground(BLACK);

    UpdateLua();
    DrawLua();
}