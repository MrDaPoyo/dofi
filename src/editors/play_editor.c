#include "editors.h"
#include "text_editor.h"
#include "../display.h"
#include "../lua_utils.h"
#include "../lua/lua_utils.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void RenderPlayEditor(void) {
    if (displayNavbar) {
        HideNavbar();
    }

    FILE* filePtr = fopen("assets/examples/gradient.lua", "r");
    if (!filePtr) {
        fprintf(stderr, "Failed to open Lua file\n");
        return;
    }

    char tempBuf[1024];
    if (!fgets(tempBuf, sizeof(tempBuf), filePtr)) {
        fprintf(stderr, "Failed to read from Lua file\n");
        fclose(filePtr);
        return;
    }
    fclose(filePtr);

    lua_State* LuaVM = NewLuaState();

    RunLuaScriptLoop(LuaVM, tempBuf);

    AnnihilateLuaState(LuaVM);
}


