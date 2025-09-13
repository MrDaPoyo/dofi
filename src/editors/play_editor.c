#include "editors.h"
#include "text_editor.h"
#include "../display.h"
#include "../lua_utils.h"

#include "../lua/lua.h"
#include "../lua/lauxlib.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define WIDTH 128
#define HEIGHT 128
#define SCALE 4

static lua_State* gLuaVM = NULL;
static bool gScriptLoaded = false;
static RenderTexture2D gRenderTex;
static bool gRenderTexInit = false;

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

/*
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
*/

static void EnsureScriptLoaded(void) {
    if (gScriptLoaded) return;

    if (!gLuaVM) {
        gLuaVM = NewLuaState();
    }

    char* script = NULL;

    /*
    if (!LoadEntireFile("assets/examples/gradient.lua", &script)) {
        fprintf(stderr, "[play] Failed to load script file.\n");
        return;
    }
    */

    script = RetrieveAllCodeFromTextEditors(editors);

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
}

static void InitRenderTexture(void) {
    if (!gRenderTexInit) {
        gRenderTex = LoadRenderTexture(WIDTH, HEIGHT);
        SetTextureFilter(gRenderTex.texture, TEXTURE_FILTER_POINT);
        gRenderTexInit = true;
    }
}

void RenderPlayEditor(void) {
    EnsureScriptLoaded();
    InitRenderTexture();

    for (int x = 0; x < WIDTH; x++)
        for (int y = 0; y < HEIGHT; y++)
            framebuffer[x][y] = (Color){0,0,0,255};

    UpdateLua();
    DrawLua();

    Image img = {
        .data = framebuffer,
        .width = WIDTH,
        .height = HEIGHT,
        .mipmaps = 1,
        .format = PIXELFORMAT_UNCOMPRESSED_R8G8B8A8
    };

    UpdateTexture(gRenderTex.texture, img.data);

    ClearBackground(BLACK);
    DrawTextureEx(gRenderTex.texture, (Vector2){0,0}, 0, SCALE, WHITE);
}
