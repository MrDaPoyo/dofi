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

int scrollLineOffset = 0;
int scrollCharOffset = 0;
#define VISIBLE_LINES 14
#define CURSOR_SCROLL_GAP
static const char *textBuffer = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1234\naaaaaa\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n20";
int cursorLine = 0;
int cursorChar = 25;

void processInput(void)
{
    if (IsKeyPressed(KEY_UP))
    {
        cursorLine--;
        if (cursorLine < 0)
            cursorLine = 0;
    }
    else if (IsKeyPressed(KEY_DOWN))
    {
        cursorLine++;
        if (cursorLine >= GetTotalLines(textBuffer))
            cursorLine = GetTotalLines(textBuffer) - 1;
        int lineLength = GetLineText(textBuffer, cursorLine) ? (int)strlen(GetLineText(textBuffer, cursorLine)) : 0;
        if (cursorChar > lineLength) {
            cursorChar = lineLength;
        }
    }
    else if (IsKeyPressed(KEY_LEFT))
    {
        cursorChar--;
        if (cursorChar < 0)
        {
            cursorLine--;
            if (cursorLine < 0)
                cursorLine = 0;
            cursorChar = 0;
        }
    }
    else if (IsKeyPressed(KEY_RIGHT))
    {
        cursorChar++;
        const char *lineText = GetLineText(textBuffer, cursorLine);
        if (lineText != NULL && cursorChar > (int)strlen(lineText))
            cursorChar = (int)strlen(lineText);
    }
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

int GetTotalLines(const char *str)
{
    if (str == NULL)
        return 0;

    int lines = 1;
    const char *p = str;
    while (*p != '\0')
    {
        if (*p == '\n')
            lines++;
        p++;
    }
    return lines;
}

void RenderTextEditor(void)
{
    processInput();

    int totalLines = GetTotalLines(textBuffer);
    int fontSize = 5;
    int charWidth = MeasureText("-", fontSize);
    int screenW = GetScreenWidth();

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