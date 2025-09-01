#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "text_editor.h"

#define VISIBLE_LINES 14
#define CURSOR_SCROLL_GAP 2
#define GAP 2

int fontSize = 5;

void RenderLine(const char *str, int index)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = index * fontSize + GAP * (index + 1) + navY;

    if (index % 2)
    {
        DrawRectangle(0, y - 1, GetScreenWidth(), fontSize + GAP, Fade(systemPalette[1], 1));
    }

    RenderString(str, GAP, y);
}

struct TextEditor editors[10] = {
    {
        .buffer = "test1234\n5678",
        .cursorLine = 0,
        .cursorChar = 0,
        .scrollLineOffset = 0,
        .scrollCharOffset = 0
    }
};

void RenderTextEditor(void)
{
    struct TextEditor *editor = &editors[0];
}

// there's a 14 line limit btw