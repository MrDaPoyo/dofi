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

    if (index % 2)
    {
        DrawRectangle(0, y - 1, GetScreenWidth(), fontSize + GAP, Fade(systemPalette[1], 1));
    }

    RenderString(str, GAP, y);
}

#include "text_editor_input.h"

const char *GetLineText(const char *str, int index)
{
    static char linebuf[256];
    if (str == NULL)
        return NULL;

    const char *p = str;
    for (int i = 0; i < index; i++)
    {
        p = strchr(p, '\n');
        if (p == NULL)
            return NULL;
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

void RenderCursor(int lineIndex, int charIndex)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int visibleRow = lineIndex - scrollLineOffset;
    if (visibleRow < 0 || visibleRow >= VISIBLE_LINES)
        return;

    const char *lineText = GetLineText(textBuffer, lineIndex);
    if (lineText == NULL)
        lineText = "";

    if (charIndex < scrollCharOffset)
        return;

    int charWidth = MeasureText("-", fontSize);
    int x = GAP + charWidth * (charIndex - scrollCharOffset);
    int y = visibleRow * fontSize + GAP * (visibleRow + 1) + navY;

    Color cursorColor = systemPalette[2];
    if (fmod(GetTime(), 1.0) < 0.5)
        cursorColor = systemPalette[3];

    DrawRectangle(x, y, 3, fontSize, cursorColor);
}

void RenderTextEditor(void)
{
    TextEditor_HandleInput();

    int totalLines = GetTotalLines(textBuffer);
    int fontSize = 5;
    int charWidth = MeasureText("-", fontSize);
    int screenW = GetScreenWidth();
    int visibleCols = (screenW - GAP * 2) / (charWidth > 0 ? charWidth : 1);

    int row = 0;
    for (row = 0; row < VISIBLE_LINES; row++)
    {
        int lineIndex = scrollLineOffset + row;
        if (lineIndex >= totalLines)
            break;
        const char *line = GetLineText(textBuffer, lineIndex);
        if (line == NULL)
            break;
        int len = (int)strlen(line);
        const char *display = line;
        if (len > scrollCharOffset)
            display = line + scrollCharOffset;
        else
            display = "";
        RenderLine(display, row);
    }

    RenderCursor(cursorLine, cursorChar);
}

// there's a 14 line limit btw