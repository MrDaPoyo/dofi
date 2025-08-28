#include "raylib.h"
#include <stdio.h>
#include <string.h>
#include <math.h>

#include "text_editor_input.h"
#include "../fonts.h"
#include "../display.h"

#define GAP 2

char textBuffer[4096] = "hello world\nthis is a test\nthree lines";
int cursorLine = 0;
int cursorChar = 0;

int scrollLineOffset = 0;
int scrollCharOffset = 0;

#define REPEAT_INITIAL_DELAY 0.40
#define REPEAT_INTERVAL 0.06

enum
{
    KR_LEFT = 0,
    KR_RIGHT,
    KR_UP,
    KR_DOWN,
    KR_BACKSPACE,
    KR_COUNT
};
static double keyNextRepeat[KR_COUNT] = {0.0};

static int BufferLength(void)
{
    return (int)strlen(textBuffer);
}

int FindAbsolutePos(const char *s, int line, int charIndex)
{
    const char *p = s;
    int curLine = 0;
    int pos = 0;

    while (*p && curLine < line)
    {
        if (*p == '\n')
            curLine++;
        p++;
        pos++;
    }

    int c = 0;
    while (*p && *p != '\n' && c < charIndex)
    {
        p++;
        pos++;
        c++;
    }

    return pos;
}

int LineLength(const char *s, int line)
{
    const char *ln = GetLineText(s, line);
    if (ln == NULL)
        return -1;
    return (int)strlen(ln);
}

static void InsertCharAt(int absPos, char ch)
{
    int len = (int)strlen(textBuffer);
    if (len + 1 >= (int)sizeof(textBuffer))
        return;
    memmove(textBuffer + absPos + 1, textBuffer + absPos, len - absPos + 1);
    textBuffer[absPos] = ch;
}

static void DeleteCharAt(int absPos)
{
    int len = (int)strlen(textBuffer);
    if (absPos >= len)
        return;
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
    if (cursorChar > 0)
        cursorChar--;
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
    if (len < 0)
        len = 0;
    if (cursorChar < len)
        cursorChar++;
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
        if (len < cursorChar)
            cursorChar = (len >= 0) ? len : 0;
    }
}

static void MoveCursorDown(void)
{
    const char *next = GetLineText(textBuffer, cursorLine + 1);
    if (next != NULL)
    {
        cursorLine++;
        int len = LineLength(textBuffer, cursorLine);
        if (len < cursorChar)
            cursorChar = (len >= 0) ? len : 0;
    }
}

int GetTotalLines(const char *s)
{
    if (s == NULL)
        return 0;
    int lines = 0;
    const char *p = s;
    if (*p != '\0')
        lines = 1;
    while (*p)
    {
        if (*p == '\n')
            lines++;
        p++;
    }
    return lines;
}

void ClampScrollOffsets(void)
{
    int total = GetTotalLines(textBuffer);
    if (scrollLineOffset < 0)
        scrollLineOffset = 0;
    int maxV = total - VISIBLE_LINES;
    if (maxV < 0)
        maxV = 0;
    if (scrollLineOffset > maxV)
        scrollLineOffset = maxV;

    int fontSize = 5;
    int charWidth = MeasureText("-", fontSize);
    int screenW = GetScreenWidth();
    int visibleCols = (screenW - GAP * 2) / (charWidth > 0 ? charWidth : 1);

    if (scrollCharOffset < 0)
        scrollCharOffset = 0;

    int maxLineLen = 0;
    for (int i = 0; i < total; i++)
    {
        int len = LineLength(textBuffer, i);
        if (len > maxLineLen)
            maxLineLen = len;
    }
    int maxH = maxLineLen - visibleCols;
    if (maxH < 0)
        maxH = 0;
    if (scrollCharOffset > maxH)
        scrollCharOffset = maxH;
}

void EnsureCursorVisible(void)
{
    if (cursorLine < scrollLineOffset)
        scrollLineOffset = cursorLine;
    else if (cursorLine >= scrollLineOffset + VISIBLE_LINES)
        scrollLineOffset = cursorLine - VISIBLE_LINES + 1;

    int fontSize = 5;
    int charWidth = MeasureText("-", fontSize);
    int screenW = GetScreenWidth();
    int visibleCols = (screenW - GAP * 2) / (charWidth > 0 ? charWidth : 1);

    if (cursorChar < scrollCharOffset)
        scrollCharOffset = cursorChar;
    else if (cursorChar >= scrollCharOffset + visibleCols)
        scrollCharOffset = cursorChar - visibleCols + 1;

    ClampScrollOffsets();
}

static void ScrollLeft(void)
{
    if (scrollCharOffset > 0)
        scrollCharOffset--;
    ClampScrollOffsets();
}
static void ScrollRight(void)
{
    scrollCharOffset++;
    ClampScrollOffsets();
}
static void ScrollUp(void)
{
    scrollLineOffset = (scrollLineOffset > 0) ? scrollLineOffset - 1 : 0;
    ClampScrollOffsets();
}
static void ScrollDown(void)
{
    scrollLineOffset++;
    ClampScrollOffsets();
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
        if (keyNextRepeat[idx] <= 0.0)
            keyNextRepeat[idx] = t + REPEAT_INITIAL_DELAY;
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

void TextEditor_HandleInput(void)
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

    bool ctrlDown = IsKeyDown(KEY_LEFT_CONTROL) || IsKeyDown(KEY_RIGHT_CONTROL);

    if (ctrlDown)
    {
        ProcessKeyRepeat(KEY_LEFT, KR_LEFT, ScrollLeft);
        ProcessKeyRepeat(KEY_RIGHT, KR_RIGHT, ScrollRight);
        ProcessKeyRepeat(KEY_UP, KR_UP, ScrollUp);
        ProcessKeyRepeat(KEY_DOWN, KR_DOWN, ScrollDown);
    }
    else
    {
        ProcessKeyRepeat(KEY_LEFT, KR_LEFT, MoveCursorLeft);
        ProcessKeyRepeat(KEY_RIGHT, KR_RIGHT, MoveCursorRight);
        ProcessKeyRepeat(KEY_UP, KR_UP, MoveCursorUp);
        ProcessKeyRepeat(KEY_DOWN, KR_DOWN, MoveCursorDown);
    }

    ProcessKeyRepeat(KEY_BACKSPACE, KR_BACKSPACE, HandleBackspace);

    if (IsKeyPressed(KEY_DELETE))
        HandleDelete();

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

    if (IsKeyPressed(KEY_PAGE_UP))
    {
        scrollLineOffset -= VISIBLE_LINES;
        if (scrollLineOffset < 0)
            scrollLineOffset = 0;
    }
    if (IsKeyPressed(KEY_PAGE_DOWN))
    {
        scrollLineOffset += VISIBLE_LINES;
    }

    float wheel = GetMouseWheelMove();
    if (wheel != 0.0f)
    {
        scrollLineOffset -= (int)wheel;
    }

    ClampScrollOffsets();
    EnsureCursorVisible();
}
