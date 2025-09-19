#include "raylib.h"
#include "editors.h"
#include "../display.h"
#include "sprite_editor.h"
#include "../colors.h"
#include "../main.h"
#include <math.h>
#include <stdio.h>

struct Spritesheet spritesheet;
struct Sprite CurrentSprite;
bool spriteEditorInitialized = false;

struct CursorPixel {
    size_t x;
    size_t y;
    size_t index;
};

static struct Sprite InitSprite(void) {
    struct Sprite s = {
        .width = MAX_SPRITE_WIDTH,
        .height = MAX_SPRITE_HEIGHT,
    };

    for (size_t y = 0; y < MAX_SPRITE_HEIGHT; y++) {
        for (size_t x = 0; x < MAX_SPRITE_WIDTH; x++) {
            s.data[y][x] = (Color){ 0, 0, 255, 255 };
        }
    }
    return s;
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

struct Sprite* SetSpritePixel(struct Sprite* sprite, size_t x, size_t y, Color color) {
    if (x < sprite->width && y < sprite->height) {
        sprite->data[y][x] = color;
    }
    return sprite;
}

struct CursorPixel getRelativePixel(Vector2 mousePos, size_t editorWidth) {
    float logicalX = mousePos.x / (float)SCALE;
    float logicalY = mousePos.y / (float)SCALE;

    float localX = logicalX - (float)DRAWING_EDITOR_X;
    float localY = logicalY - (float)DRAWING_EDITOR_Y;

    struct CursorPixel cursor = { 0, 0, 0 };

    if (localX < 0 || localY < 0) {
        cursor.x = cursor.y = editorWidth;
        cursor.index = 0;
        return cursor;
    }

    float spritePixelSize = (float)EDITOR_SCALE;
    float spriteDrawWidth = (float)editorWidth * spritePixelSize;
    float spriteDrawHeight = (float)editorWidth * spritePixelSize;

    if (localX >= spriteDrawWidth || localY >= spriteDrawHeight) {
        cursor.x = cursor.y = editorWidth;
        cursor.index = 0;
        return cursor;
    }

    size_t px = (size_t)floorf(localX / spritePixelSize);
    size_t py = (size_t)floorf(localY / spritePixelSize);
    if (px >= editorWidth) px = editorWidth - 1;
    if (py >= editorWidth) py = editorWidth - 1;

    cursor.x = px;
    cursor.y = py;
    cursor.index = py * editorWidth + px;
    return cursor;
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

    if (IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        Vector2 mouse = GetMousePositionForEditors();
        const struct CursorPixel pixelCursor = getRelativePixel(mouse, CurrentSprite.width);
        if (pixelCursor.x < CurrentSprite.width && pixelCursor.y < CurrentSprite.height) {
            SetSpritePixel(&CurrentSprite, pixelCursor.x, pixelCursor.y, RED);
        }
    }
}
