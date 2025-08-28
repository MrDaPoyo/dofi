#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../scale.h"

#define GAP (2 * GetDisplayScale())

void RenderLine(const char *str, int index)
{
    int scale = GetDisplayScale();
    int navY = GetScaledNavHeight();
    int fontSize = 5 * scale;
    float fontHeight = MeasureTextEx(GeneralFont, "A", fontSize, 1.0f).y;

    int lineSpacing = (int)(fontHeight * 1.2f); 

    int y = navY + GAP + (index * lineSpacing);
    RenderString(str, GAP, y);
}


void RenderTextEditor(void)
{
    RenderLine("hey there babe 0", 0);
    RenderLine("hey there babe 1", 1);
    RenderLine("hey there babe 2", 2);
}
