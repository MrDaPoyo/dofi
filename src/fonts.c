#include "raylib.h"
#include "fonts.h"
#include "colors.h"

Font GeneralFont;

void LoadFonts(void) {
    if (FileExists("assets/fonts/dofi.ttf")) {
        GeneralFont = LoadFont("assets/fonts/dofi.ttf");
    } else {
        GeneralFont = GetFontDefault();
    }
}

void UnloadFonts(void) {
    UnloadFont(GeneralFont);
}

void RenderString(const char* str, int x, int y) {
    int offsetX = 0;
    char buf[2] = { 0, 0 };
    const float fontSize = 5.0f;
    const float spacing = 1.0f;

    const int cellWidth = FONT_WIDTH + FONT_SPACING;

    while (*str) {
        buf[0] = *str;

        Vector2 pos = (Vector2){ (float)(x + offsetX), (float)y };

        DrawTextPro(
            GeneralFont,
            buf,
            pos,
            (Vector2){ 0, 0 },
            0.0f,
            fontSize,
            spacing,
            systemPalette[4]
            );

        offsetX += cellWidth;
        str++;
    }
}
