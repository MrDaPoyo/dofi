#include "colors.h"
#include "raylib.h"
#include "fonts.h"
#include "display.h"

#include <stdio.h>

#define INTRO_TIME 1.75f // seconds
#define TILE_SIZE 4

struct Tile {
    int x;
    int y;
    float delay;
    float height;
};

struct Tile newTile(const int x, const int y) {
    float delay = ((x / TILE_SIZE) + (y / TILE_SIZE)) * 0.05f;

    struct Tile tile = {
        .x = x,
        .y = y,
        .delay = delay,
        .height = TILE_SIZE
    };
    return tile;
}

struct Tile updateTile(struct Tile tile, float progress) {
    float localProgress = progress - tile.delay;
    if (localProgress < 0) localProgress = 0;
    if (localProgress > 1) localProgress = 1;

    tile.height = TILE_SIZE * (1.0f - localProgress);

    return tile;
}

void renderTile(struct Tile tile) {
    DrawRectangle(tile.x, tile.y, TILE_SIZE, tile.height, systemPalette[1]);
}

#define MAX_TILES ((WIDTH / TILE_SIZE) * (HEIGHT / TILE_SIZE))
struct Tile tiles[MAX_TILES];

void RenderIntro() {
    float elapsed = 0.0f;
    SetTargetFPS(60);
    RenderTexture2D gRenderTex = LoadRenderTexture(WIDTH, HEIGHT);
    Texture2D logo = LoadTexture("assets/logo_small_nobg.png");
    Vector2 logoVec = { WIDTH / 2 - logo.width * 2, HEIGHT / 2 - logo.height * 2 };

    LoadFonts();

    int heightDelay = 1;
    int totalTiles = 0;

    for (int y = 0; y < HEIGHT; y += TILE_SIZE) {
        int widthDelay = heightDelay;
        for (int x = 0; x < WIDTH; x += TILE_SIZE) {
            tiles[totalTiles] = newTile(x, y);
            widthDelay++;
            totalTiles++;
        }
        heightDelay++;
    }

    printf("Size of total tiles: %d\n", totalTiles);

    while (elapsed < INTRO_TIME && !WindowShouldClose()) {
        float dt = GetFrameTime();
        elapsed += dt;
        float progress = elapsed / INTRO_TIME * 10; // time / max_time * speed multiplier

        BeginDrawing();
        BeginTextureMode(gRenderTex);

        ClearBackground(systemPalette[3]);

        DrawTextureEx(logo, logoVec, 0, 4, WHITE);

        for (int tileIndex = 0; tileIndex < totalTiles; tileIndex++) {
            tiles[tileIndex] = updateTile(tiles[tileIndex], progress);
            renderTile(tiles[tileIndex]);
        }

        EndTextureMode();

        DrawTexturePro(gRenderTex.texture, (Rectangle){ 0, 0, (float)gRenderTex.texture.width, -(float)gRenderTex.texture.height }, (Rectangle){ 0, 0, WIDTH * SCALE, HEIGHT * SCALE }, (Vector2){ 0, 0 }, 0.0f, WHITE);

        EndDrawing();
    }

    UnloadFonts();
    UnloadRenderTexture(gRenderTex);
}
