#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderSoundEditor(void)
{
    if (!displayNavbar) {
        ShowNavbar();
    }
    int navY = GetNavHeight();
    DrawText("Sound Editor", 10, navY + 40, 20, BLUE);
}
