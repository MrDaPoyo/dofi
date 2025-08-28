#include "raylib.h"
#include <stdint.h>
#include "colors.h"
#include "editors/editors.h"
#include "fonts.h"
#include "scale.h"

#define WIDTH 128
#define HEIGHT 128
#define SCALE 4

uint8_t framebuffer[WIDTH * HEIGHT];

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
    EDITOR_TEXT,
    EDITOR_SPRITE,
    EDITOR_MAP,
    EDITOR_SOUND,
    EDITOR_PLAY,
    EDITOR_COUNT
} EditorTab;

SuperTab currentSuperTab = SUPER_CLI;
EditorTab currentEditorTab = EDITOR_TEXT;

Texture2D iconSprite;
Texture2D iconMap;
Texture2D iconMusic;
Texture2D iconPlay;
Texture2D iconCode;
Texture2D iconUndefined;

#define NAVBAR_HEIGHT 16
#define BUTTON_SIZE 12
#define GAP 2

#define SCALED_NAV_HEIGHT (NAVBAR_HEIGHT * SCALE)

static inline Vector2 GetMousePositionForEditors(void)
{
    Vector2 m = GetMousePosition();
    return m;
}

int GetDisplayScale(void)
{
    return SCALE;
}

int GetScaledNavHeight(void)
{
    return SCALED_NAV_HEIGHT;
}

int GetEditorCanvasHeight(void)
{
    return HEIGHT;
}

int GetScaledEditorHeight(void)
{
    return HEIGHT * SCALE;
}

Color NavbarBGColor;
Color NavbarBorderColor;
Color EditorBGColor;
Color CliBGColor;
Color CliTextColor;

void RenderCLI(void);
void RenderEditors(void);
void switchSuperTab(void);

int main()
{
    NavbarBGColor = systemPalette[2];
    NavbarBorderColor = systemPalette[3];
    EditorBGColor = systemPalette[0];
    CliBGColor = systemPalette[5];
    CliTextColor = systemPalette[4];

    InitWindow(WIDTH * SCALE, HEIGHT * SCALE, "Dofi v0.0");
    SetConfigFlags(0);
    SetTargetFPS(45);

    LoadFonts();

    iconSprite = LoadTexture("assets/icons/brush.png");
    iconMap = LoadTexture("assets/icons/tile.png");
    iconMusic = LoadTexture("assets/icons/music.png");
    iconPlay = LoadTexture("assets/icons/play.png");
    iconCode = LoadTexture("assets/icons/code.png");
    iconUndefined = LoadTexture("assets/icons/undefined.png");

    Texture2D cursorTexture = LoadTexture("assets/icons/mouse_with_shadow.png");
    SetMouseOffset(-(cursorTexture.width * SCALE) / 2, -(cursorTexture.height * SCALE) / 2);
    HideCursor();

    while (!WindowShouldClose())
    {
        Vector2 mousePos = GetMousePositionForEditors();

        if (IsKeyPressed(KEY_F1))
            switchSuperTab();

        BeginDrawing();
        ClearBackground(EditorBGColor);

        if (currentSuperTab == SUPER_CLI)
        {
            RenderCLI();
        }
        else
        {
            RenderEditors();
            DrawTextureEx(cursorTexture, (Vector2){mousePos.x, mousePos.y}, 0.0f, SCALE, WHITE);
        }

        EndDrawing();
    }

    UnloadTexture(iconSprite);
    UnloadTexture(iconMap);
    UnloadTexture(iconMusic);
    UnloadTexture(iconPlay);
    UnloadTexture(iconCode);
    UnloadTexture(iconUndefined);
    UnloadTexture(cursorTexture);

    CloseWindow();
    UnloadFonts();
    return 0;
}

void switchSuperTab()
{
    if (currentSuperTab == SUPER_CLI)
        currentSuperTab = SUPER_EDITORS;
    else
        currentSuperTab = SUPER_CLI;
}

void RenderCLI()
{
    DrawRectangle(0, 0, WIDTH * SCALE, HEIGHT * SCALE, CliBGColor);

    DrawText("CLI Mode (F1 = Switch)", 10, 10, 10, YELLOW);
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

        Rectangle btn;
        if (i == EDITOR_PLAY)
        {
            btn = (Rectangle){GetScreenWidth() - scaledGap - scaledBtnSize, 4 * SCALE / 2, scaledBtnSize, scaledBtnSize};
        }
        else
        {
            btn = (Rectangle){scaledGap + x, 4 * SCALE / 2, scaledBtnSize, scaledBtnSize};
        }

        DrawRectangleRec(btn, (i == (int)currentEditorTab) ? (Color){255, 169, 133, 255} : systemPalette[3]);

        if (CheckCollisionPointRec(GetMousePositionForEditors(), btn) && IsMouseButtonPressed(MOUSE_LEFT_BUTTON))
        {
            currentEditorTab = (EditorTab)i;
        }

        Texture2D icon;
        if (i == EDITOR_TEXT)
            icon = iconCode;
        else if (i == EDITOR_SPRITE)
            icon = iconSprite;
        else if (i == EDITOR_MAP)
            icon = iconMap;
        else if (i == EDITOR_SOUND)
            icon = iconMusic;
        else if (i == EDITOR_PLAY)
            icon = iconPlay;
        else
            icon = iconUndefined;

        {
            float iconSize = (float)(scaledBtnSize - 4 * SCALE);
            float dstX;

            if (i == EDITOR_PLAY)
            {
                dstX = (float)(GetScreenWidth() - scaledGap - scaledBtnSize) + 2.0f * SCALE;
            }
            else
            {
                dstX = (float)(scaledGap + x) + 2.0f * SCALE;
            }

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

    int editorY = scaledNavHeight;
    int editorH = GetScaledEditorHeight();

    DrawRectangle(0, editorY, GetScreenWidth(), editorH, EditorBGColor);

    BeginScissorMode(0, editorY, GetScreenWidth(), editorH);

    if (currentEditorTab == EDITOR_TEXT)
        RenderTextEditor();
    if (currentEditorTab == EDITOR_SPRITE)
        RenderSpriteEditor();
    if (currentEditorTab == EDITOR_MAP)
        RenderMapEditor();
    if (currentEditorTab == EDITOR_SOUND)
        RenderSoundEditor();
    if (currentEditorTab == EDITOR_PLAY)
        RenderPlayEditor();

    EndScissorMode();
}