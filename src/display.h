#ifndef DOFI_DISPLAY_H
#define DOFI_DISPLAY_H

#include "raylib.h"

int GetDisplayScale(void);
int GetScaledNavHeight(void);
int GetNavHeight(void);
int GetEditorCanvasWidth(void);
int GetEditorCanvasHeight(void);
int GetScaledEditorHeight(void);

#define WIDTH 128
#define HEIGHT 128
#define SCALE 4

#define NAVBAR_HEIGHT 16
#define NAVBAR_BUTTON_SIZE 12
#define NAVBAR_BUTTON_GAP 2

extern bool displayNavbar;
extern Color framebuffer[WIDTH][HEIGHT];

typedef enum {
    SUPER_CLI,
    SUPER_EDITORS
} SuperTab;

typedef enum {
    EDITOR_TEXT,
    EDITOR_SPRITE,
    EDITOR_MAP,
    EDITOR_SOUND,
    EDITOR_PLAY,
    EDITOR_COUNT
} EditorTab;

extern SuperTab currentSuperTab;
extern EditorTab currentEditorTab;

void HideNavbar(void);
void ShowNavbar(void);

void switchSuperTab(void);
void RenderCLI(void);
void RenderEditors(void);

#endif
