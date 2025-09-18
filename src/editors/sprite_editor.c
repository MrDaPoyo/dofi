#include "raylib.h"
#include "editors.h"
#include "../display.h"
#include "sprite_editor.h"

#include "../colors.h"

struct Spritesheet spritesheet;
struct Sprite CurrentSprite;

void SetSpritePixel(size_t spriteId, size_t x, size_t y, Color color) {
    spritesheet.sprites[spriteId]->data[x][y] = color;
}

#define DRAWING_TILES_X 2
#define DRAWING_TILES_Y 2

bool spriteEditorInitialized = false;

/*
void renderImageEditor(struct Sprite sprite) {
    DrawRectangle(0, 0, EDITOR_WIDTH, EDITOR_HEIGHT, BLACK);
}
*/

struct Spritesheet newSprite(struct Spritesheet spritesheet, size_t id) {
    struct Sprite s = {
        .width = MAX_SPRITE_WIDTH,
        .height = MAX_SPRITE_HEIGHT,
    };

    for (size_t y = 0; y < MAX_SPRITE_HEIGHT; y++) {
        for (size_t x = 0; x < MAX_SPRITE_WIDTH; x++) {
            s.data[y][x] = (Color){0, 0, 0, 255};
        }
    }

    size_t row = id / TOTAL_MAX_SPRITES_SQRT;
    size_t col = id % TOTAL_MAX_SPRITES_SQRT;

    spritesheet.sprites[row][col] = s;

    return spritesheet;
}

void RenderSpriteEditor(void)
{
    if (!displayNavbar) {
        ShowNavbar();
    }
}
