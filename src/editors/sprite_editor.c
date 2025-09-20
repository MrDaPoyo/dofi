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
Color CurrentColor;
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
            s.data[y][x] = (Color){ 0, 0, 0, 255 };
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

struct CursorPixel getRelativePixel(Vector2 mousePos, size_t regionWidth, int offsetX, int offsetY) {
    float logicalX = mousePos.x / (float)SCALE;
    float logicalY = mousePos.y / (float)SCALE;

    float localX = logicalX - (float)offsetX;
    float localY = logicalY - (float)offsetY;

    struct CursorPixel cursor = { regionWidth, regionWidth, 0 };

    if (localX < 0 || localY < 0) {
        return cursor;
    }

    float spritePixelSize = (float)EDITOR_SCALE;
    float spriteDrawWidth = (float)regionWidth * spritePixelSize;
    float spriteDrawHeight = (float)regionWidth * spritePixelSize;

    if (localX >= spriteDrawWidth || localY >= spriteDrawHeight) {
        return cursor;
    }

    size_t px = (size_t)floorf(localX / spritePixelSize);
    size_t py = (size_t)floorf(localY / spritePixelSize);

    if (px >= regionWidth) px = regionWidth - 1;
    if (py >= regionWidth) py = regionWidth - 1;

    cursor.x = px;
    cursor.y = py;
    cursor.index = py * regionWidth + px;
    return cursor;
}

bool colorsEqual(Color a, Color b) {
    return (a.r == b.r) &&
           (a.g == b.g) &&
           (a.b == b.b) &&
           (a.a == b.a);
}

void renderColorPalette() {
    size_t paletteSize = sizeof(palette) / sizeof(palette[0]);
    size_t gridSize = (size_t)sqrt((double)paletteSize);

    int paletteX = DRAWING_EDITOR_X + (int)(CurrentSprite.width * EDITOR_SCALE) + 10;
    int paletteY = DRAWING_EDITOR_Y;

    Vector2 mouse = GetMousePositionForEditors();
    struct CursorPixel paletteCursor = getRelativePixel(mouse, gridSize, paletteX, paletteY);

    if (paletteCursor.x < gridSize && paletteCursor.y < gridSize && IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        size_t index = paletteCursor.index;
        if (index < paletteSize) {
            CurrentColor = palette[index];
        }
    }

    size_t index = 0;
    size_t y_height = 0;
    size_t x_height = 0;

    for (y_height = 0; y_height < gridSize; y_height++) {
        for (x_height = 0; x_height < gridSize; x_height++) {
            if (index >= paletteSize) break;

            Color c = palette[index];
            DrawRectangle(
                paletteX + (int)(x_height * EDITOR_SCALE),
                paletteY + (int)(y_height * EDITOR_SCALE),
                EDITOR_SCALE,
                EDITOR_SCALE,
                c
            );

            if (colorsEqual(CurrentColor, c))
                DrawRectangleLines(
                    paletteX + (int)(x_height * EDITOR_SCALE),
                    paletteY + (int)(y_height * EDITOR_SCALE),
                    EDITOR_SCALE,
                    EDITOR_SCALE,
                    systemPalette[2]
                );

            index++;
        }
    }

    DrawRectangleLines(
    paletteX,
        paletteY,
        x_height * EDITOR_SCALE,
        y_height * EDITOR_SCALE,
        systemPalette[2]
        );
}


void RenderSpriteEditor(void) {
    if (!displayNavbar) {
        ShowNavbar();
    }

    if (!spriteEditorInitialized) {
        CurrentSprite = InitSprite();
        CurrentColor = palette[0];
        spriteEditorInitialized = true;
    }

    renderImageEditor(CurrentSprite);
    renderColorPalette(CurrentSprite);

    if (IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        Vector2 mouse = GetMousePositionForEditors();
        struct CursorPixel pixelCursor = getRelativePixel(mouse, CurrentSprite.width, DRAWING_EDITOR_X, DRAWING_EDITOR_Y);

        if (pixelCursor.x < CurrentSprite.width && pixelCursor.y < CurrentSprite.height) {
            SetSpritePixel(&CurrentSprite, pixelCursor.x, pixelCursor.y, CurrentColor);
        }
    }
}
