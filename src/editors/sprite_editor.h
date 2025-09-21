#ifndef DOFI_SPRITE_EDITOR_H
#define DOFI_SPRITE_EDITOR_H

#include "raylib.h"
#include <stddef.h>
#include <math.h>

#define TOTAL_MAX_SPRITES 64
#define TOTAL_MAX_SPRITES_SQRT 8
#define MAX_SPRITE_WIDTH 8
#define MAX_SPRITE_HEIGHT 8
#define EDITOR_SCALE 8
#define EDITOR_WIDTH (MAX_SPRITE_WIDTH * EDITOR_SCALE)
#define EDITOR_HEIGHT (MAX_SPRITE_HEIGHT * EDITOR_SCALE)
#define SPRITESHEET_SCALE (EDITOR_SCALE / 4)
#define SPRITESHEET_X 0
#define SPRITESHEET_Y (DRAWING_EDITOR_Y + EDITOR_HEIGHT + FONT_GAP)


struct Sprite {
    Color data[MAX_SPRITE_HEIGHT][MAX_SPRITE_WIDTH];
    size_t height;
    size_t width;
};

struct Spritesheet {
    struct Sprite sprites[TOTAL_MAX_SPRITES_SQRT][TOTAL_MAX_SPRITES_SQRT];
};

extern struct Sprite CurrentSprite;
extern struct Spritesheet spritesheet;

static inline struct Sprite* GetSpriteById(int id) {
    if (id < 0 || id >= (int)TOTAL_MAX_SPRITES) return NULL;
    int row = id / TOTAL_MAX_SPRITES_SQRT;
    int col = id % TOTAL_MAX_SPRITES_SQRT;
    return &spritesheet.sprites[row][col];
}

void RenderSpriteEditor(void);

#endif // DOFI_SPRITE_EDITOR_H
