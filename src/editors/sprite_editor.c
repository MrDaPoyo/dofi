#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderSpriteEditor(void)
{
    if (!displayNavbar) {
        ShowNavbar();
    }
    int navY = GetNavHeight();
    DrawText("Sprite Editor", 10, navY + 40, 20, RED);
}
