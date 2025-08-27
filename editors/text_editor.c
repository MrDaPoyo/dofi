#include "raylib.h"
#include "../colors.h"

#define NAVBAR_HEIGHT 16

void RenderTextEditor(void)
{
    DrawText("Text Editor", 10, NAVBAR_HEIGHT + 40, 20, systemPalette[4]);
}
