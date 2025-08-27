#include "raylib.h"
#include "fonts.h"
#include "colors.h"

Font GeneralFont;

void LoadFonts(void)
{
    if (FileExists("assets/fonts/pico-8.ttf"))
    {
        GeneralFont = LoadFont("assets/fonts/pico-8.ttf");
    }
    else
    {
        GeneralFont = GetFontDefault();
    }
}

void UnloadFonts(void)
{
    UnloadFont(GeneralFont);
}

void RenderString(const char *str, int x, int y)
{
    int offsetX = 0;
    char buf[2] = {0, 0};
    while (*str)
    {
        buf[0] = *str;
        Vector2 pos = (Vector2){x + offsetX, (float)y};
        DrawTextPro(GeneralFont, buf, pos, (Vector2){0, 0}, 0.0f, 20, 1, systemPalette[4]);
        offsetX += (int)(MeasureTextEx(GeneralFont, buf, 20, 1).x) + 2;
        str++;
    }
}