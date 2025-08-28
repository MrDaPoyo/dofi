#include "display.h"

int GetDisplayScale(void)
{
    return SCALE;
}

int GetScaledNavHeight(void)
{
    return NAVBAR_HEIGHT * SCALE;
}

int GetEditorCanvasHeight(void)
{
    return HEIGHT;
}

int GetScaledEditorHeight(void)
{
    return HEIGHT * SCALE;
}

int GetNavHeight(void)
{
    return NAVBAR_HEIGHT;
}

int GetEditorCanvasWidth(void)
{
    return WIDTH;
}

void RenderCLI(void);
void RenderEditors(void);
void switchSuperTab(void);