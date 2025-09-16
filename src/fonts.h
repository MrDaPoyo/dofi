#ifndef DOFI_FONTS_H
#define DOFI_FONTS_H

#define FONT_WIDTH 3
#define FONT_HEIGHT 5
#define FONT_GAP 2
#define FONT_SPACING 1

#include "raylib.h"
#include <stddef.h>
#include <stdlib.h>

extern Font GeneralFont;

void LoadFonts(void);
void UnloadFonts(void);

void RenderString(char *str, int x, int y);
void RenderTextEditorString(char* str, size_t length, int x, int y);
void RenderStringWrap(char* str, int x, int y);
int IsComment(const char* line);

#endif