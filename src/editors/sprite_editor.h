#ifndef DOFI_SPRITE_EDITOR_H
#define DOFI_SPRITE_EDITOR_H

#include "raylib.h"
#include <stddef.h>
#include <math.h>

#define TOTAL_MAX_SPRITES 64
#define TOTAL_MAX_SPRITES_SQRT 8
#define MAX_SPRITE_WIDTH 16
#define MAX_SPRITE_HEIGHT 16
#define EDITOR_SCALE 4
#define EDITOR_WIDTH (MAX_SPRITE_WIDTH * EDITOR_SCALE)
#define EDITOR_HEIGHT (MAX_SPRITE_HEIGHT * EDITOR_SCALE)

struct Sprite {
    Color data[MAX_SPRITE_HEIGHT][MAX_SPRITE_WIDTH];
    size_t height;
    size_t width;
};

struct Spritesheet {
    struct Sprite sprites[TOTAL_MAX_SPRITES_SQRT][TOTAL_MAX_SPRITES_SQRT];
};

extern struct Sprite CurrentSprite;

void RenderSpriteEditor(void);

#endif // DOFI_SPRITE_EDITOR_H
