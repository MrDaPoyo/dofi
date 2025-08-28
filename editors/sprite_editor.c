#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../scale.h"

void RenderSpriteEditor(void)
{
    int scale = GetDisplayScale();
    int navY = GetScaledNavHeight();
    DrawText("Sprite Editor", 10 * scale, navY + 40 * scale, 20 * scale, RED);
}
