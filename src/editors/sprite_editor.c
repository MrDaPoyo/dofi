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
static int selectedSpriteRow = 0;
static int selectedSpriteCol = 0;

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

static void InitSpritesheet(struct Spritesheet* sheet) {
    for (size_t row = 0; row < TOTAL_MAX_SPRITES_SQRT; row++) {
        for (size_t col = 0; col < TOTAL_MAX_SPRITES_SQRT; col++) {
            sheet->sprites[row][col] = InitSprite();
        }
    }

    struct Sprite* s0 = &sheet->sprites[0][0];
    Color Y = (Color){255, 255, 0, 255};
    Color K = (Color){0, 0, 0, 255};
    for (size_t y = 0; y < s0->height; y++) {
        for (size_t x = 0; x < s0->width; x++) {
            s0->data[y][x] = Y;
        }
    }
    s0->data[2][2] = K; s0->data[2][5] = K;
    s0->data[5][2] = K; s0->data[5][3] = K; s0->data[5][4] = K; s0->data[5][5] = K;
}

#define DRAWING_EDITOR_X 2
#define DRAWING_EDITOR_Y (2 + NAVBAR_HEIGHT)
#define PALETTE_SCALE (EDITOR_SCALE * 2)

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

struct CursorPixel getRelativePixel(Vector2 mousePos, size_t regionWidth, int offsetX, int offsetY, float pixelSize) {
    float logicalX = mousePos.x / (float)SCALE;
    float logicalY = mousePos.y / (float)SCALE;

    float localX = logicalX - (float)offsetX;
    float localY = logicalY - (float)offsetY;

    struct CursorPixel cursor = { regionWidth, regionWidth, 0 };

    if (localX < 0 || localY < 0) {
        return cursor;
    }

    float spriteDrawWidth = (float)regionWidth * pixelSize;
    float spriteDrawHeight = (float)regionWidth * pixelSize;

    if (localX >= spriteDrawWidth || localY >= spriteDrawHeight) {
        return cursor;
    }

    size_t px = (size_t)floorf(localX / pixelSize);
    size_t py = (size_t)floorf(localY / pixelSize);

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

    int paletteX = WIDTH - DRAWING_EDITOR_X - PALETTE_SCALE * sqrt(sizeof(palette) / sizeof(palette[0]));
    int paletteY = DRAWING_EDITOR_Y;

    Vector2 mouse = GetMousePositionForEditors();
    struct CursorPixel paletteCursor = getRelativePixel(mouse, gridSize, paletteX, paletteY, (float)PALETTE_SCALE);

    if (paletteCursor.x < gridSize && paletteCursor.y < gridSize && IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        size_t index = paletteCursor.index;
        if (index < paletteSize) {
            CurrentColor = palette[index];
        }
    }

    size_t index = 0;
    size_t y_height = 0;
    size_t x_width = 0;

    for (y_height = 0; y_height < gridSize; y_height++) {
        for (x_width = 0; x_width < gridSize; x_width++) {
            if (index >= paletteSize) break;

            Color c = palette[index];
            DrawRectangle(
                paletteX + (int)(x_width * PALETTE_SCALE),
                paletteY + (int)(y_height * PALETTE_SCALE),
                PALETTE_SCALE,
                PALETTE_SCALE,
                c
            );

            if (colorsEqual(CurrentColor, c))
                DrawRectangleLines(
                    paletteX + (int)(x_width * PALETTE_SCALE),
                    paletteY + (int)(y_height * PALETTE_SCALE),
                    PALETTE_SCALE,
                    PALETTE_SCALE,
                    systemPalette[2]
                );

            index++;
        }
    }

    DrawRectangleLines(
    paletteX,
        paletteY,
        (int)x_width * PALETTE_SCALE,
        (int)y_height * PALETTE_SCALE,
        systemPalette[2]
        );
}

void renderSpritesheet(struct Spritesheet* sheet) {
    size_t spriteSize = MAX_SPRITE_WIDTH;
    int spriteDrawSize = (int)(spriteSize * SPRITESHEET_SCALE);
    Vector2 rawMouse = GetMousePositionForEditors();
    Vector2 mouse = (Vector2){ rawMouse.x / (float)SCALE, rawMouse.y / (float)SCALE };

    for (size_t row = 0; row < TOTAL_MAX_SPRITES_SQRT; row++) {
        for (size_t col = 0; col < TOTAL_MAX_SPRITES_SQRT; col++) {
            int x = SPRITESHEET_X + (int)(col * spriteDrawSize);
            int y = SPRITESHEET_Y + (int)(row * spriteDrawSize);

            struct Sprite* sprite = &sheet->sprites[row][col];

            DrawRectangle(x, y, spriteDrawSize, spriteDrawSize, BLACK);

            for (size_t sy = 0; sy < sprite->height; sy++) {
                for (size_t sx = 0; sx < sprite->width; sx++) {
                    Color c = sprite->data[sy][sx];
                    DrawRectangle(
                        x + (int)(sx * SPRITESHEET_SCALE),
                        y + (int)(sy * SPRITESHEET_SCALE),
                        SPRITESHEET_SCALE,
                        SPRITESHEET_SCALE,
                        c
                    );
                }
            }

            bool hovered = CheckCollisionPointRec(mouse, (Rectangle){x, y, (float)spriteDrawSize, (float)spriteDrawSize});
            bool isSelected = (int)row == selectedSpriteRow && (int)col == selectedSpriteCol;

            if (hovered) {
                DrawRectangleLines(x, y, spriteDrawSize, spriteDrawSize, systemPalette[2]);
                if (IsMouseButtonPressed(MOUSE_BUTTON_LEFT)) {
                    if (sprite->width == 0 || sprite->height == 0 || sprite->width > MAX_SPRITE_WIDTH || sprite->height > MAX_SPRITE_HEIGHT) {
                        *sprite = InitSprite();
                    }
                    selectedSpriteRow = (int)row;
                    selectedSpriteCol = (int)col;
                    CurrentSprite = *sprite;
                }
            }
            else if (isSelected) {
                DrawRectangleLines(x, y, spriteDrawSize, spriteDrawSize, systemPalette[3]);
            }
        }
    }
}

void RenderSpriteEditor(void) {
    if (!displayNavbar) {
        ShowNavbar();
    }

    if (!spriteEditorInitialized) {
        CurrentSprite = InitSprite();
        CurrentColor = palette[0];
        InitSpritesheet(&spritesheet);
        spriteEditorInitialized = true;
    }

    renderImageEditor(CurrentSprite);
    renderColorPalette();
    renderSpritesheet(&spritesheet);

    if (IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        Vector2 rawMouse = GetMousePositionForEditors();
        struct CursorPixel pixelCursor = getRelativePixel(rawMouse, CurrentSprite.width, DRAWING_EDITOR_X, DRAWING_EDITOR_Y, (float)EDITOR_SCALE);
        if (pixelCursor.x < CurrentSprite.width && pixelCursor.y < CurrentSprite.height) {
            Color* existing = &CurrentSprite.data[pixelCursor.y][pixelCursor.x];
            if (!colorsEqual(*existing, CurrentColor)) {
                SetSpritePixel(&CurrentSprite, pixelCursor.x, pixelCursor.y, CurrentColor);
                spritesheet.sprites[selectedSpriteRow][selectedSpriteCol] = CurrentSprite;
            }
        }
    }
}
