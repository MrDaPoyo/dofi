#include "raylib.h"
#include "editors.h"
#include "../display.h"
#include "sprite_editor.h"

struct Spritesheet spritesheet;
struct Sprite CurrentSprite;

void SetSpritePixel(size_t spriteId, size_t x, size_t y, Color color) {
    spritesheet.sprites[spriteId]->data[x][y] = color;
}

void RenderSpriteEditor(void)
{
    if (!displayNavbar) {
        ShowNavbar();
    }
    int navY = GetNavHeight();
    DrawText("Sprite Editor", 10, navY + 40, 20, RED);
}
