#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../scale.h"

void RenderPlayEditor(void)
{
    int scale = GetDisplayScale();
    int navY = GetScaledNavHeight();
    DrawText("Play Editor", 10 * scale, navY + 40 * scale, 20 * scale, systemPalette[4]);
}
