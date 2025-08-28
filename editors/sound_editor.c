#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../display.h"

void RenderSoundEditor(void)
{
    int navY = GetNavHeight();
    DrawText("Sound Editor", 10, navY + 40, 20, BLUE);
}
