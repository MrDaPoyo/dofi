#include "play_editor.h"

#include "editors.h"
#include "text_editor.h"
#include "../colors.h"
#include "../display.h"
#include "../lua_utils.h"

#include "../lua/lua.h"
#include "../lua/lauxlib.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static lua_State* gLuaVM = NULL;
bool gScriptLoaded = false;
RenderTexture2D playRenderTex;
bool gRenderTexInit = false;
char gLastError[256] = {0};
bool isPlaying = false;
bool isErrorValidated = false;

void ResetPlayEditor(void) {
    gScriptLoaded = false;
    isErrorValidated = false;
    hasCodeChanged = false;

    if (gLuaVM) {
        lua_close(gLuaVM);
        gLuaVM = NULL;
    }
}

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
    if (hasCodeChanged) {
        ResetPlayEditor();
    }

    if (gScriptLoaded) return;
    if (isErrorValidated) return;

    if (!gLuaVM) {
        gLuaVM = NewLuaState();
    }

    char* scriptRaw = RetrieveAllCodeFromTextEditors(editors);
    if (!scriptRaw) {
        snprintf(gLastError, sizeof(gLastError), "allocation failure while retrieving code");
        isErrorValidated = true;
        return;
    }
    char* script = scriptRaw;

    if (luaL_loadstring(gLuaVM, script) != LUA_OK) {
        const char* err = lua_tostring(gLuaVM, -1);
        fprintf(stderr, "[lua] load error: %s\n", err);
        isErrorValidated = true;
        if (err) {
            snprintf(gLastError, sizeof(gLastError), "%s", err);
        }
        lua_pop(gLuaVM, 1);
        free(script);
        return;
    }

    if (lua_pcall(gLuaVM, 0, 0, 0) != LUA_OK) {
        const char* err = lua_tostring(gLuaVM, -1);
        fprintf(stderr, "[lua] runtime error: %s\n", err);
        isErrorValidated = true;
        if (err) {
            snprintf(gLastError, sizeof(gLastError), "%s", err);
        }
        lua_pop(gLuaVM, 1);
        free(script);
        return;
    }

    CallLuaFunctionSafe(gLuaVM, "init");
    gScriptLoaded = true;
    free(script);
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
        playRenderTex = LoadRenderTexture(WIDTH, HEIGHT);
        SetTextureFilter(playRenderTex.texture, TEXTURE_FILTER_POINT);
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
 
    UpdateTexture(playRenderTex.texture, img.data);
 
    ClearBackground(BLACK);
    DrawTextureEx(playRenderTex.texture, (Vector2){0,0}, 0, SCALE, WHITE);

    if (!gScriptLoaded && gLastError[0] != '\0') {
        DrawRectangle(0, 0, WIDTH, NAVBAR_HEIGHT, systemPalette[2]);
        RenderString("error detected!!", FONT_SPACING, NAVBAR_HEIGHT / 2 - FONT_HEIGHT / 2);
        char shortErr[256];
        snprintf(shortErr, sizeof(shortErr), "%s", gLastError);
        RenderStringWrap(shortErr, FONT_SPACING, NAVBAR_HEIGHT + FONT_GAP);
    }
}