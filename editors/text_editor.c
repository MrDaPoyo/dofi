#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include <string.h>
#include <math.h>

#define GAP 2

void RenderLine(const char *str, int index)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = index * fontSize + GAP * (index + 1) + navY;

    RenderString(str, GAP, y);
}

const char *GetLineText(const char *str, int index)
{
    static char linebuf[256];
    if (str == NULL) return NULL;

    const char *p = str;
    for (int i = 0; i < index; i++)
    {
        p = strchr(p, '\n');
        if (p == NULL) return NULL;
        p++;
    }

    size_t len = 0;
    while (p[len] != '\0' && p[len] != '\n' && len < sizeof(linebuf) - 1)
    {
        len++;
    }
    memcpy(linebuf, p, len);
    linebuf[len] = '\0';
    return linebuf;
}

const char *testString = "hello world\nthis is a test\nthree lines";

void RenderCursor(int lineIndex, int charIndex)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = lineIndex * fontSize + GAP * (lineIndex + 1) + navY;
    const char *lineText = GetLineText(testString, lineIndex);
    if (lineText == NULL) lineText = "";
    int x = MeasureText("-", fontSize) * (charIndex + 1) - 2;
    Color cursorColor = systemPalette[2];

    if (fmod(GetTime(), 1.0) < 0.5)
        cursorColor = systemPalette[3];
    

    DrawRectangle(x, y, 3, fontSize, cursorColor);
}


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

    RenderCursor(0, 0);
}

// there's a 14 line limit btw