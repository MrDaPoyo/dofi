#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderPlayEditor(void)
{
    int navY = GetNavHeight();
    DrawText("Play Editor", 10, navY + 40, 20, systemPalette[4]);
}
