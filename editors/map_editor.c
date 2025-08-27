#include "raylib.h"
#include "../colors.h"

#define NAVBAR_HEIGHT 16

void RenderMapEditor(void)
{
    DrawText("Map Editor", 10, NAVBAR_HEIGHT + 40, 20, GREEN);
}
