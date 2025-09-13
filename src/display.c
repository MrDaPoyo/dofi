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

bool displayNavbar = false;

void HideNavbar(void) {
    displayNavbar = false;
}

void ShowNavbar(void) {
    displayNavbar = true;
}