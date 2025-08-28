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

static char textBuffer[4096] = "hello world\nthis is a test\nthree lines";
static int cursorLine = 0;
static int cursorChar = 0;

// Add key repeat state and helper
#define REPEAT_INITIAL_DELAY 0.40   // seconds before first repeat
#define REPEAT_INTERVAL      0.06   // seconds between repeats

enum { KR_LEFT = 0, KR_RIGHT, KR_UP, KR_DOWN, KR_BACKSPACE, KR_COUNT };
static double keyNextRepeat[KR_COUNT] = { 0.0 };

static int BufferLength(void)
{
    return (int)strlen(textBuffer);
}

static int FindAbsolutePos(const char *s, int line, int charIndex)
{
    const char *p = s;
    int curLine = 0;
    int pos = 0;

    while (*p && curLine < line)
    {
        if (*p == '\n') curLine++;
        p++; pos++;
    }

    int c = 0;
    while (*p && *p != '\n' && c < charIndex)
    {
        p++; pos++; c++;
    }

    return pos;
}

static int LineLength(const char *s, int line)
{
    const char *ln = GetLineText(s, line);
    if (ln == NULL) return -1;
    return (int)strlen(ln);
}

static void InsertCharAt(int absPos, char ch)
{
    int len = (int)strlen(textBuffer);
    if (len + 1 >= (int)sizeof(textBuffer)) return;
    memmove(textBuffer + absPos + 1, textBuffer + absPos, len - absPos + 1);
    textBuffer[absPos] = ch;
}

static void DeleteCharAt(int absPos)
{
    int len = (int)strlen(textBuffer);
    if (absPos >= len) return;
    memmove(textBuffer + absPos, textBuffer + absPos + 1, len - absPos);
}

static void HandleBackspace(void)
{
    if (cursorChar > 0)
    {
        int pos = FindAbsolutePos(textBuffer, cursorLine, cursorChar);
        DeleteCharAt(pos - 1);
        cursorChar--;
    }
    else if (cursorLine > 0)
    {
        int pos = FindAbsolutePos(textBuffer, cursorLine, 0);
        if (pos > 0 && textBuffer[pos - 1] == '\n')
        {
            DeleteCharAt(pos - 1);
            cursorLine--;
            int prevLen = LineLength(textBuffer, cursorLine);
            cursorChar = (prevLen >= 0) ? prevLen : 0;
        }
    }
}

static void HandleDelete(void)
{
    int pos = FindAbsolutePos(textBuffer, cursorLine, cursorChar);
    DeleteCharAt(pos);
}

static void MoveCursorLeft(void)
{
    if (cursorChar > 0) cursorChar--;
    else if (cursorLine > 0)
    {
        cursorLine--;
        int len = LineLength(textBuffer, cursorLine);
        cursorChar = (len >= 0) ? len : 0;
    }
}

static void MoveCursorRight(void)
{
    int len = LineLength(textBuffer, cursorLine);
    if (len < 0) len = 0;
    if (cursorChar < len) cursorChar++;
    else
    {
        const char *next = GetLineText(textBuffer, cursorLine + 1);
        if (next != NULL)
        {
            cursorLine++;
            cursorChar = 0;
        }
    }
}

static void MoveCursorUp(void)
{
    if (cursorLine > 0)
    {
        cursorLine--;
        int len = LineLength(textBuffer, cursorLine);
        if (len < cursorChar) cursorChar = (len >= 0) ? len : 0;
    }
}

static void MoveCursorDown(void)
{
    const char *next = GetLineText(textBuffer, cursorLine + 1);
    if (next != NULL)
    {
        cursorLine++;
        int len = LineLength(textBuffer, cursorLine);
        if (len < cursorChar) cursorChar = (len >= 0) ? len : 0;
    }
}

static void ProcessKeyRepeat(int key, int idx, void (*action)(void))
{
    double t = GetTime();
    bool down = IsKeyDown(key);

    if (IsKeyPressed(key))
    {
        action();
        keyNextRepeat[idx] = t + REPEAT_INITIAL_DELAY;
        return;
    }

    if (down)
    {
        if (keyNextRepeat[idx] <= 0.0) keyNextRepeat[idx] = t + REPEAT_INITIAL_DELAY;
        if (t >= keyNextRepeat[idx])
        {
            action();
            keyNextRepeat[idx] = t + REPEAT_INTERVAL;
        }
    }
    else
    {
        keyNextRepeat[idx] = 0.0;
    }
}

static void HandleTextInput(void)
{
    int code = 0;
    while ((code = GetCharPressed()) > 0)
    {
        if (code >= 32)
        {
            int pos = FindAbsolutePos(textBuffer, cursorLine, cursorChar);
            InsertCharAt(pos, (char)code);
            cursorChar++;
        }
    }

    // use repeating handlers for arrows and backspace
    ProcessKeyRepeat(KEY_LEFT, KR_LEFT, MoveCursorLeft);
    ProcessKeyRepeat(KEY_RIGHT, KR_RIGHT, MoveCursorRight);
    ProcessKeyRepeat(KEY_UP, KR_UP, MoveCursorUp);
    ProcessKeyRepeat(KEY_DOWN, KR_DOWN, MoveCursorDown);
    ProcessKeyRepeat(KEY_BACKSPACE, KR_BACKSPACE, HandleBackspace);

    // single-press behaviour for other keys
    if (IsKeyPressed(KEY_DELETE)) HandleDelete();
    if (IsKeyPressed(KEY_ENTER))
    {
        int pos = FindAbsolutePos(textBuffer, cursorLine, cursorChar);
        InsertCharAt(pos, '\n');
        cursorLine++;
        cursorChar = 0;
    }
    if (IsKeyPressed(KEY_HOME))
    {
        cursorChar = 0;
    }
    if (IsKeyPressed(KEY_END))
    {
        int len = LineLength(textBuffer, cursorLine);
        cursorChar = (len >= 0) ? len : 0;
    }
}

void RenderCursor(int lineIndex, int charIndex)
{
    int navY = GetNavHeight();
    int fontSize = 5;

    int y = lineIndex * fontSize + GAP * (lineIndex + 1) + navY;
    const char *lineText = GetLineText(textBuffer, lineIndex);
    if (lineText == NULL) lineText = "";
    int x = MeasureText("-", fontSize) * (charIndex + 1) - 2;
    Color cursorColor = systemPalette[2];

    if (fmod(GetTime(), 1.0) < 0.5)
        cursorColor = systemPalette[3];
    

    DrawRectangle(x, y, 3, fontSize, cursorColor);
}


void RenderTextEditor(void)
{
    HandleTextInput();

    int lineIndex = 0;

    while (lineIndex < 14)
    {
        const char *line = GetLineText(textBuffer, lineIndex);
        if (line == NULL) break;
        RenderLine(line, lineIndex);
        lineIndex++;
    }

    RenderCursor(cursorLine, cursorChar);
}

// there's a 14 line limit btw