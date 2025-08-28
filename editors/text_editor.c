#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../scale.h"

#define GAP 2

void RenderLine(const char *str, int index)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = index * fontSize + GAP * index + navY;

    char y_str[32];
    sprintf(y_str, "%d", index);

    RenderString(y_str, GAP, y);
}

void RenderTextEditor(void)
{
    RenderLine("hey there babe 0", 0);
    RenderLine("hey there babe 1", 1);
    RenderLine("hey there babe 2", 2);
    RenderLine("hey there babe 2", 3);
    RenderLine("hey there babe 2", 4);
    RenderLine("hey there babe 2", 5);
    RenderLine("hey there babe 2", 6);
    RenderLine("hey there babe 2", 7);
    RenderLine("hey there babe 2", 8);
    RenderLine("hey there babe 2", 9);
    RenderLine("hey there babe 2", 10);
    RenderLine("hey there babe 2", 11);
    RenderLine("hey there babe 2", 12);
    RenderLine("hey there babe 2", 13);
    RenderLine("hey there babe 14", 14);
    RenderLine("hey there babe 15", 15);
    RenderLine("hey there babe 16", 16);
    RenderLine("hey there babe 17", 17);
    RenderLine("hey there babe 18", 18);
    RenderLine("hey there babe 19", 19);
}
