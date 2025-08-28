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

    int y = index * fontSize + GAP * (index + 1) + navY;

    char y_str[32];
    sprintf(y_str, "%d", index);

    RenderString(y_str, GAP, y);
}

void RenderTextEditor(void)
{
    RenderLine("hey there babe 0", 0);
    RenderLine("hey there babe 1", 1);
    RenderLine("hey there babe 2", 2);
    RenderLine("hey there babe 3", 3);
    RenderLine("hey there babe 4", 4);
    RenderLine("hey there babe 5", 5);
    RenderLine("hey there babe 6", 6);
    RenderLine("hey there babe 7", 7);
    RenderLine("hey there babe 8", 8);
    RenderLine("hey there babe 9", 9);
    RenderLine("hey there babe 10", 10);
    RenderLine("hey there babe 11", 11);
    RenderLine("hey there babe 12", 12);
    RenderLine("hey there babe 13", 13);
}

// there is a 14 line limit. I still gotta implement scrolling and cursors and actual text, etc etc

