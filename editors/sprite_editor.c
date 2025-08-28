#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderSpriteEditor(void)
{
    int navY = GetNavHeight();
    DrawText("Sprite Editor", 10, navY + 40, 20, RED);
}
