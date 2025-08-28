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
    const float fontSize = 5.0f;
    const float spacing = 1.0f;
    while (*str)
    {
        buf[0] = *str;
        Vector2 pos = (Vector2){(float)(x + offsetX), (float)y};
        DrawTextPro(GeneralFont, buf, pos, (Vector2){0, 0}, 0.0f, fontSize, spacing, systemPalette[4]);
        offsetX += (int)(MeasureTextEx(GeneralFont, buf, fontSize, spacing).x) + 1;
        str++;
    }
}