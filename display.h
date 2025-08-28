#ifndef DOFI_DISPLAY_H
#define DOFI_DISPLAY_H


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

void switchSuperTab(void);
void RenderCLI(void);
void RenderEditors(void);

#endif
