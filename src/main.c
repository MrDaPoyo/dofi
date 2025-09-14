#include "raylib.h"
#include "colors.h"
#include "editors/editors.h"
#include "fonts.h"
#include "display.h"
#include "icons.h"
#include "editors/text_editor.h"
#include "intro.h"
#include "audio.h"
#include "editors/play_editor.h"

SuperTab currentSuperTab = SUPER_CLI;
EditorTab currentEditorTab = EDITOR_TEXT;

#define SCALED_NAV_HEIGHT (NAVBAR_HEIGHT * SCALE)

static inline Vector2 GetMousePositionForEditors(void) {
    Vector2 m = GetMousePosition();
    return m;
}

Color NavbarBGColor;
Color NavbarBorderColor;
Color EditorBGColor;
Color CliBGColor;
Color CliTextColor;

RenderTexture2D gRenderTex;

int main() {
    NavbarBGColor = systemPalette[2];
    NavbarBorderColor = systemPalette[3];
    EditorBGColor = systemPalette[0];
    CliBGColor = systemPalette[1];
    CliTextColor = systemPalette[4];

    InitWindow(WIDTH * SCALE, HEIGHT * SCALE, "Dofi v0.0");

    SetExitKey(KEY_F10);
    InitAudioDevice();

    RenderIntro();

    SetConfigFlags(0);
    SetTargetFPS(60);

    LoadFonts();

    gRenderTex = LoadRenderTexture(WIDTH, HEIGHT);

    LoadIcons();
    LoadAudios();

    Texture2D cursorTexture = LoadTexture("assets/icons/mouse_with_shadow.png");
    SetMouseOffset(-(cursorTexture.width * SCALE) / 2, -(cursorTexture.height * SCALE) / 2);
    HideCursor();

    // initialize uhhhhhhh editors
    InitTextEditors();
    HideNavbar();

    while (!WindowShouldClose()) {
        if (IsKeyPressed(KEY_ESCAPE)) {
            if (currentEditorTab == EDITOR_PLAY) {
                ShowNavbar();
                currentEditorTab = EDITOR_TEXT;
            } else {
                SwitchSuperTab();
            }
        }

        BeginDrawing();

        BeginTextureMode(gRenderTex);
        ClearBackground(EditorBGColor);

        if (currentSuperTab == SUPER_CLI) {
            RenderCLI();
        } else if (currentEditorTab == EDITOR_PLAY) {
            HideNavbar();
            RenderPlayEditor();
        } else {
            const Vector2 mousePos = GetMousePositionForEditors();
            ShowNavbar();
            RenderEditors();
            DrawTextureEx(cursorTexture,
                (Vector2){ mousePos.x / (float)SCALE, mousePos.y / (float)SCALE },
                0.0f, 1.0f, WHITE);
            if (IsMouseButtonPressed(MOUSE_LEFT_BUTTON) || IsMouseButtonPressed(MOUSE_RIGHT_BUTTON)) {
                PlaySound(audios[0]);
            }
        }

        EndTextureMode();

        SetTextureFilter(gRenderTex.texture, TEXTURE_FILTER_POINT);
        DrawTexturePro(gRenderTex.texture,
            (Rectangle){ 0, 0, (float)gRenderTex.texture.width, -(float)gRenderTex.texture.height },
            (Rectangle){ 0, 0, WIDTH * SCALE, HEIGHT * SCALE },
            (Vector2){ 0, 0 }, 0.0f, WHITE);

        EndDrawing();
    }

    UnloadIcons();

    UnloadTexture(cursorTexture);
    UnloadRenderTexture(gRenderTex);

    CloseWindow();
    UnloadAudios();
    UnloadFonts();
    return 0;
}

void SwitchSuperTab() {
    if (currentSuperTab == SUPER_CLI)
        currentSuperTab = SUPER_EDITORS;
    else
        currentSuperTab = SUPER_CLI;
}

void RenderCLI() {
    if (displayNavbar) {
        HideNavbar();
    }
    DrawRectangle(0, 0, WIDTH, HEIGHT, CliBGColor);

}

void RenderEditors(void) {
    int navH = NAVBAR_HEIGHT;
    int btnSize = NAVBAR_BUTTON_SIZE;
    int gap = NAVBAR_BUTTON_GAP;

    // navbar
    DrawRectangle(0, 0, WIDTH, navH, NavbarBGColor);

    // buttons
    for (int i = 0; i < EDITOR_COUNT; i++) {
        int x = i * (btnSize + gap);

        Rectangle btn;
        if (i == EDITOR_PLAY) {
            btn = (Rectangle){ (float)(WIDTH - gap - btnSize), 2.0f, (float)btnSize, (float)btnSize };
        } else {
            btn = (Rectangle){ (float)(gap + x), 2.0f, (float)btnSize, (float)btnSize };
        }

        DrawRectangleRec(btn, (i == (int)currentEditorTab) ? (Color){ 255, 169, 133, 255 } : systemPalette[3]);

        Vector2 m = GetMousePositionForEditors();
        Vector2 mLog = (Vector2){ m.x / (float)SCALE, m.y / (float)SCALE };
        if (CheckCollisionPointRec(mLog, btn) && IsMouseButtonPressed(MOUSE_LEFT_BUTTON)) {
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
            float iconSize = (float)(btnSize - 4);
            float dstX;

            if (i == EDITOR_PLAY) {
                dstX = (float)(WIDTH - gap - btnSize) + 2.0f;
            } else {
                dstX = (float)(gap + x) + 2.0f;
            }

            float dstY = (float)btn.y + ((float)btn.height - iconSize) * 0.5f;
            DrawTexturePro(
                icon,
                (Rectangle){ 0, 0, icon.width, icon.height },
                (Rectangle){ dstX, dstY, iconSize, iconSize },
                (Vector2){ 0, 0 },
                0.0f,
                WHITE);
        }
    }

    int editorY = NAVBAR_HEIGHT;
    int editorH = GetEditorCanvasHeight();

    DrawRectangle(0, editorY, WIDTH, editorH, EditorBGColor);

    BeginScissorMode(0, editorY, WIDTH, editorH);

    if (currentEditorTab == EDITOR_TEXT)
        RenderTextEditor();
    if (currentEditorTab == EDITOR_SPRITE)
        RenderSpriteEditor();
    if (currentEditorTab == EDITOR_MAP)
        RenderMapEditor();
    if (currentEditorTab == EDITOR_SOUND)
        RenderSoundEditor();

    EndScissorMode();
}