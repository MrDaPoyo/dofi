#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include <string.h>

#define GAP 2

void RenderLine(const char *str, int index)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = index * fontSize + GAP * (index + 1) + navY;

    RenderString(str, GAP, y);
}


const char *testString = "hello world\nthis is a test\nthree lines";

void RenderTextEditor(void)
{
    char buffer[256];
    strncpy(buffer, testString, sizeof(buffer));
    buffer[sizeof(buffer) - 1] = '\0';

    char *line = strtok(buffer, "\n");
    int lineIndex = 0;

    while (line != NULL && lineIndex < 14)
    {
        RenderLine(line, lineIndex);
        line = strtok(NULL, "\n");
        lineIndex++;
    }
}

// there's a 14 line limit btw