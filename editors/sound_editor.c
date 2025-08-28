#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../scale.h"

void RenderSoundEditor(void)
{
    int scale = GetDisplayScale();
    int navY = GetScaledNavHeight();
    DrawText("Sound Editor", 10 * scale, navY + 40 * scale, 20 * scale, BLUE);
}
