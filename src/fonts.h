#ifndef DOFI_FONTS_H
#define DOFI_FONTS_H

#include "raylib.h"

extern Font GeneralFont;

void LoadFonts(void);
void UnloadFonts(void);

void RenderString(const char *str, int x, int y);

#endif