#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderMapEditor(void)
{
    int navY = GetNavHeight();
    DrawText("Map Editor", 10, navY + 40, 20, GREEN);
}
