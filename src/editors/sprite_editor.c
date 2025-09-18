#include "raylib.h"
#include "editors.h"
#include "../display.h"
#include "sprite_editor.h"
#include "../colors.h"

struct Spritesheet spritesheet;
struct Sprite CurrentSprite;
bool spriteEditorInitialized = false;

static struct Sprite InitSprite(void) {
    struct Sprite s = {
        .width = MAX_SPRITE_WIDTH,
        .height = MAX_SPRITE_HEIGHT,
    };

    for (size_t y = 0; y < MAX_SPRITE_HEIGHT; y++) {
        for (size_t x = 0; x < MAX_SPRITE_WIDTH; x++) {
            s.data[y][x] = (Color){0, 0, 255, 255};
        }
    }
    return s;
}

void SetSpritePixel(size_t spriteId, size_t x, size_t y, Color color) {
    size_t row = spriteId / TOTAL_MAX_SPRITES_SQRT;
    size_t col = spriteId % TOTAL_MAX_SPRITES_SQRT;
    spritesheet.sprites[row][col].data[y][x] = color;
}

#define DRAWING_EDITOR_X 2
#define DRAWING_EDITOR_Y (2 + NAVBAR_HEIGHT)

void renderImageEditor(struct Sprite sprite) {
    DrawRectangle(
        DRAWING_EDITOR_X,
        DRAWING_EDITOR_Y,
        sprite.width * EDITOR_SCALE,
        sprite.height * EDITOR_SCALE,
        BLACK
    );

    for (size_t y = 0; y < sprite.height; y++) {
        for (size_t x = 0; x < sprite.width; x++) {
            Color c = sprite.data[y][x];
            DrawRectangle(
                DRAWING_EDITOR_X + (int)(x * EDITOR_SCALE),
                DRAWING_EDITOR_Y + (int)(y * EDITOR_SCALE),
                EDITOR_SCALE,
                EDITOR_SCALE,
                c
            );
        }
    }

    DrawRectangleLines(
        DRAWING_EDITOR_X,
        DRAWING_EDITOR_Y,
        sprite.width * EDITOR_SCALE,
        sprite.height * EDITOR_SCALE,
        systemPalette[2]
    );
}


struct Spritesheet newSprite(struct Spritesheet spritesheet, size_t id) {
    struct Sprite s = InitSprite();

    size_t row = id / TOTAL_MAX_SPRITES_SQRT;
    size_t col = id % TOTAL_MAX_SPRITES_SQRT;

    spritesheet.sprites[row][col] = s;

    return spritesheet;
}

void RenderSpriteEditor(void) {
    if (!displayNavbar) {
        ShowNavbar();
    }

    if (!spriteEditorInitialized) {
        CurrentSprite = InitSprite();
        spriteEditorInitialized = true;
    }

    renderImageEditor(CurrentSprite);
}
