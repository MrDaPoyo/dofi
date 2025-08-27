#include "raylib.h"
#include "../colors.h"
#include "editors.h"

#define NAVBAR_HEIGHT 16

void RenderPlayEditor(void)
{
    DrawText("Play Editor", 10, NAVBAR_HEIGHT + 40, 20, systemPalette[4]);
}
