#ifndef DOFI_FONTS_H
#define DOFI_FONTS_H

#define FONT_WIDTH 3
#define FONT_HEIGHT 5
#define FONT_GAP 2
#define FONT_SPACING 1

#include "raylib.h"

extern Font GeneralFont;

void LoadFonts(void);
void UnloadFonts(void);

void RenderString(const char *str, int x, int y);

#endif