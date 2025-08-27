#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"

#define NAVBAR_HEIGHT 16

void RenderTextEditor(void)
{
    DrawText("Text Editor", 10, NAVBAR_HEIGHT + 40, 20, systemPalette[4]);
    RenderString("hey there babe", 20, NAVBAR_HEIGHT + 80);
}
