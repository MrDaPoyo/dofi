#ifndef DOFI_PLAY_EDITOR_H
#define DOFI_PLAY_EDITOR_H

#include <stdbool.h>
#include "../lua/lua.h"

extern bool isPlaying;
extern bool isErrorValidated;
extern bool isScriptLoaded;
extern char gLastError[256];
extern lua_State LuaVM;

void ResetPlayEditor(void);

#endif //DOFI_PLAY_EDITOR_H