#include "raylib.h"
#include <stdint.h>

#define WIDTH 128
#define HEIGHT 128
#define SCALE 4

uint8_t framebuffer[WIDTH * HEIGHT];

Color palette[16] = {
    {0, 0, 0, 255},       // 0  Black
    {16, 8, 32, 255},     // 1  Dark Purple
    {81, 30, 67, 255},    // 2  Violet
    {204, 116, 83, 255},  // 3  Burnt Umber
    {174, 181, 189, 255}, // 4  French Gray
    {255, 255, 255, 255}, // 5  White
    {233, 56, 65, 255},   // 6  Imperial Red
    {235, 108, 130, 255}, // 7  Bright Pink
    {241, 137, 45, 255},  // 8  Tangerine
    {255, 209, 157, 255}, // 9  Sunset
    {255, 233, 71, 255},  // 10 Maize
    {30, 138, 76, 255},   // 11 Sea Green
    {52, 235, 164, 255},  // 12 Aquamarine
    {5, 68, 148, 255},    // 13 Polynesian Blue
    {77, 128, 201, 255},  // 14 Glaucous
    {125, 62, 191, 255}   // 15 Blue Violet
};

void pset(int x, int y, uint8_t color)
{
    if (x < 0 || y < 0 || x >= WIDTH || y >= HEIGHT)
        return;
    framebuffer[y * WIDTH + x] = color;
}

typedef enum
{
    SUPER_CLI,
    SUPER_EDITORS
} SuperTab;

typedef enum
{
    EDITOR_SPRITE,
    EDITOR_MAP,
    EDITOR_SOUND,
    EDITOR_COUNT
} EditorTab;

SuperTab currentSuperTab = SUPER_CLI;
EditorTab currentEditorTab = EDITOR_SPRITE;

Texture2D iconSprite;
Texture2D iconMap;
Texture2D iconMusic;
Texture2D iconPlay;
Texture2D iconCode;

#define NAVBAR_HEIGHT 16
#define BUTTON_SIZE 12
#define GAP 2

Color NavbarBGColor;
Color NavbarBorderColor;
Color CliBGColor;
Color CliTextColor;

void RenderCLI(Texture2D tex);
void RenderEditors(void);
void RenderSpriteEditor(void);
void RenderMapEditor(void);
void RenderSoundEditor(void);

int main()
{
    NavbarBGColor = palette[3];
    NavbarBorderColor = palette[7];
    CliBGColor = palette[15];
    CliTextColor = palette[5];

    InitWindow(WIDTH * SCALE, HEIGHT * SCALE, "Dofi v0.0");
    SetConfigFlags(0);
    SetTargetFPS(30);

    Image img = GenImageColor(WIDTH, HEIGHT, BLACK);
    Texture2D tex = LoadTextureFromImage(img);
    UnloadImage(img);

    iconSprite = LoadTexture("assets/icons/brush.png");
    iconMap = LoadTexture("assets/icons/tile.png");
    iconMusic = LoadTexture("assets/icons/music.png");
    iconPlay = LoadTexture("assets/icons/play.png");
    iconCode = LoadTexture("assets/icons/code.png");

    while (!WindowShouldClose())
    {
        if (IsKeyPressed(KEY_F1))
            currentSuperTab = SUPER_CLI;
        if (IsKeyPressed(KEY_F2))
            currentSuperTab = SUPER_EDITORS;

        BeginDrawing();
        ClearBackground(RAYWHITE);

        if (currentSuperTab == SUPER_CLI)
        {
            RenderCLI(tex);
        }
        else
        {
            RenderEditors();
        }

        EndDrawing();
    }

    UnloadTexture(iconSprite);
    UnloadTexture(iconMap);
    UnloadTexture(iconMusic);
    UnloadTexture(iconPlay);
    UnloadTexture(iconCode);

    UnloadTexture(tex);
    CloseWindow();
    return 0;
}

void RenderCLI(Texture2D tex)
{
    for (int i = 0; i < WIDTH; i++)
    {
        pset(i, i, (i / 8) % 16);
    }

    Color pixels[WIDTH * HEIGHT];
    for (int i = 0; i < WIDTH * HEIGHT; i++)
    {
        pixels[i] = palette[framebuffer[i]];
    }
    UpdateTexture(tex, pixels);

    DrawTextureEx(tex, (Vector2){0, 0}, 0, SCALE, WHITE);

    DrawText("CLI Mode (F2 = Editors)", 10, 10, 10, YELLOW);
}

void RenderEditors(void)
{
    int scaledNavHeight = NAVBAR_HEIGHT * SCALE;
    int scaledBtnSize = BUTTON_SIZE * SCALE;
    int scaledGap = GAP * SCALE;

    // navbar
    DrawRectangle(0, 0, GetScreenWidth(), scaledNavHeight, NavbarBGColor);
    DrawRectangleLines(0, 0, GetScreenWidth(), scaledNavHeight, NavbarBorderColor);

    // buttons
    for (int i = 0; i < EDITOR_COUNT; i++)
    {
        int x = i * (scaledBtnSize + scaledGap);
        Rectangle btn = {scaledGap + x, 4 * SCALE / 2, scaledBtnSize, scaledBtnSize};
        DrawRectangleRec(btn, (i == currentEditorTab) ? (Color){255, 169, 133, 255} : (Color){154, 56, 63, 255});

        if (CheckCollisionPointRec(GetMousePosition(), btn) && IsMouseButtonPressed(MOUSE_LEFT_BUTTON))
        {
            currentEditorTab = (EditorTab)i;
        }

        Texture2D icon;
        if (i == EDITOR_SPRITE)
            icon = iconSprite;
        else if (i == EDITOR_MAP)
            icon = iconMap;
        else if (i == EDITOR_SOUND)
            icon = iconMusic;

        {
            float iconSize = (float)(scaledBtnSize - 4 * SCALE);
            float dstX = (float)(scaledGap + x) + 2.0f * SCALE;
            float dstY = (float)btn.y + ((float)btn.height - iconSize) * 0.5f;
            DrawTexturePro(
                icon,
                (Rectangle){0, 0, icon.width, icon.height},
                (Rectangle){dstX, dstY, iconSize, iconSize},
                (Vector2){0, 0},
                0.0f,
                WHITE);
        }
    }

    DrawText("Editors Mode (F1 = CLI)", 10 * SCALE, scaledNavHeight + 10 * SCALE, 10 * SCALE, BLACK);

    if (currentEditorTab == EDITOR_SPRITE)
        RenderSpriteEditor();
    if (currentEditorTab == EDITOR_MAP)
        RenderMapEditor();
    if (currentEditorTab == EDITOR_SOUND)
        RenderSoundEditor();
}

void RenderSpriteEditor(void)
{
    DrawText("Sprite Editor", 10, NAVBAR_HEIGHT + 40, 20, RED);
}

void RenderMapEditor(void)
{
    DrawText("Map Editor", 10, NAVBAR_HEIGHT + 40, 20, GREEN);
}

void RenderSoundEditor(void)
{
    DrawText("Sound Editor", 10, NAVBAR_HEIGHT + 40, 20, BLUE);
}

void RenderTextEditor(void)
{
    DrawText("Text Editor", 10, NAVBAR_HEIGHT + 40, 20, BLACK);
}

void RenderImageEditor(void)
{
    DrawText("Image Editor", 10, NAVBAR_HEIGHT + 40, 20, PURPLE);
}