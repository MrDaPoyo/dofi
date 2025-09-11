#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include "text_editor.h"
#include "text_buffer.h"

#include <string.h>
#include <math.h>
#include <stdlib.h>

#define VISIBLE_LINES 14
#define CURSOR_SCROLL_GAP 2
#define GAP 2
#define FONT_SIZE 5
#define FONT_WIDTH 4
#define LINE_CHAR_WIDTH (WIDTH /(FONT_WIDTH + GAP))

void RenderLine(const char* str, int index, size_t realIndex) {
    int navY = GetNavHeight();

    int y = index * FONT_SIZE + GAP * (index + 1) + navY;

    if (realIndex % 2) {
        DrawRectangle(0, y - 1, GetScreenWidth(), FONT_SIZE + GAP, systemPalette[1]);
    }

    RenderString(str, GAP, y);
};

void RenderStatBar(struct TextEditor editor) {
    DrawRectangle(0, HEIGHT - FONT_SIZE - GAP, GetScreenWidth(), FONT_SIZE + GAP, systemPalette[2]);
    char buffer[64];

    sprintf(buffer, "line: %li/%li; char: %li; scroll %li/%li;", editor.cursorLine + 1, editor.buffer.totalLines + 1, editor.cursorChar, editor.scrollOffsetX, editor.scrollOffsetY);
    RenderString(buffer, GAP / 2, HEIGHT - FONT_SIZE - GAP / 2);
}

void RenderCursor(struct TextEditor editor) {
    int navY = GetNavHeight();

    const size_t lineIndex = editor.cursorLine - editor.scrollOffsetY;

    int x = GAP + (editor.cursorChar - editor.scrollOffsetX) * FONT_WIDTH;
    int y = lineIndex * FONT_SIZE + GAP * (lineIndex + 1) + navY;

    Color cursorColor = systemPalette[2];

    if (fmod(GetTime(), 1.0) < 0.5)
        cursorColor = systemPalette[3];

    DrawRectangle(x, y, 3, FONT_SIZE, cursorColor);
}

struct LineCollection retrieveAllLines(const struct TextBuffer* buffer, size_t y_index) {
    struct LineCollection result = { .count = 0 };
    const char* work_buffer = buffer->buffer;

    const char* lineStart = work_buffer;
    size_t currentLine = 0;

    for (unsigned int i = 0; ; i++) {
        char c = work_buffer[i];

        if (c == '\n' || c == '\0') {
            if (currentLine >= y_index && result.count < VISIBLE_LINES) {
                result.lines[result.count] = lineStart;
                result.count++;
            }

            currentLine++;
            lineStart = &work_buffer[i + 1];

            if (c == '\0')
                break;
        }
    }
    return result;
}

size_t retrieveLinesCount(struct TextBuffer buffer) {
    const struct LineCollection lines = retrieveAllLines(&buffer, 0);
    return lines.count;
}

void RenderBuffer(struct TextEditor editor) {
    struct TextBuffer buffer = editor.buffer;
    size_t y_index = editor.scrollOffsetY;
    size_t x_index = editor.scrollOffsetX;

    struct LineCollection lc = retrieveAllLines(&buffer, y_index);

    for (int i = 0; i < lc.count; i++) {
        const char* line = lc.lines[i];
        if (line == NULL)
            continue;

        const char* line_end = line;
        while (*line_end != '\n' && *line_end != '\0')
            line_end++;

        size_t line_len = line_end - line;

        if (x_index > line_len)
            x_index = line_len;

        const char* start = line + x_index;
        const char* end = start;
        while (*end != '\n' && *end != '\0')
            end++;

        unsigned long len = end - start;
        char temp[1025]; // max 1024 chars per line
        if (len >= sizeof(temp))
            len = sizeof(temp) - 1;
        memcpy(temp, start, len);
        temp[len] = '\0';

        RenderLine(temp, i, i + y_index);
    }
}


struct TextEditor editors[10];
struct TextEditor* editor = &editors[0];

void InitTextEditors(void) {
    for (size_t i = 0; i < sizeof(editors) / sizeof(editors[0]); i++) {
        editors[i] = (struct TextEditor){
            .buffer = createBuffer(0), // ts makes a buffer with zero bytes, but that's expanded when the buffer's first character is appended
            .cursorChar = 0,
            .cursorLine = 0,
            .scrollOffsetX = 0,
            .scrollOffsetY = 0,
        };
    }
}

void PrintBufferLength(void) {
    if (editor->buffer.bufferSize > 0) {
        printf("Buffer Length: %li\n", editor->buffer.bufferSize);
    }
}

char pressedKey;

void RenderTextEditor(void) {
    pressedKey = GetCharPressed();
    if (IsKeyPressed(KEY_RIGHT)) {
        if (editor->cursorChar < editor->buffer.totalChars && editor->cursorChar < strlen(GetLineText(editor->buffer, editor->cursorLine))) {
            editor->cursorChar++;
        } else {
            char* line = GetLineText(editor->buffer, editor->cursorLine + 1);
            if (line != NULL) {
                editor->cursorChar = 0;
                editor->cursorLine++;
            }
            free(line);
        }
    }
    if (IsKeyPressed(KEY_LEFT)) {
        if (editor->cursorChar > 0) {
            editor->cursorChar--;
        } else if (editor->cursorLine > 0) {
            char* line = GetLineText(editor->buffer, editor->cursorLine - 1);
            if (line) {
                editor->cursorLine--;
                editor->cursorChar = strlen(line);
                free(line);
            } else {
                editor->cursorChar = 0;
            }
        } else {
            editor->cursorChar = 0;
        }
    }
    if (IsKeyPressed(KEY_DOWN)) {
        if (editor->cursorLine < retrieveLinesCount(editor->buffer) - 1) {
            editor->cursorLine++;

            char* line = GetLineText(editor->buffer, editor->cursorLine);
            if (line != NULL) {
                size_t lineLen = strlen(line);
                if (editor->cursorChar > lineLen) {
                    editor->cursorChar = lineLen;
                }
            }
            free(line);
        }
    }
    if (IsKeyPressed(KEY_UP)) {
        if (editor->cursorLine > 0) {
            editor->cursorLine--;

            char* line = GetLineText(editor->buffer, editor->cursorLine);
            if (line != NULL) {
                size_t lineLen = strlen(line);
                if (editor->cursorChar > lineLen) {
                    editor->cursorChar = lineLen;
                }
            }
            free(line);
        }
    }
    if (IsKeyPressed(KEY_ENTER)) {
        pressedKey = '\n';
    }
    if (IsKeyPressed(KEY_BACKSPACE)) {
        *editor = removeCharacter(editor);
    }

    if (pressedKey != 0) {
        *editor = insertCharacter(editor, pressedKey);
        if (pressedKey == '\n') {
            editor->cursorChar = 0;
            editor->cursorLine++;
        } else {
            editor->cursorChar++;
        }
        pressedKey = GetCharPressed();

        if (editor->cursorChar >= editor->scrollOffsetX + LINE_CHAR_WIDTH - GAP) {
            editor->scrollOffsetX = editor->cursorChar - (LINE_CHAR_WIDTH - GAP) + 1;
        } else if (editor->cursorChar < editor->scrollOffsetX + GAP) {
            if (editor->cursorChar > GAP)
                editor->scrollOffsetX = editor->cursorChar - GAP;
            else
                editor->scrollOffsetX = 0;
        }

        if (editor->cursorLine >= editor->scrollOffsetY + VISIBLE_LINES - GAP && editor->scrollOffsetY < editor->buffer.totalLines) {
            editor->scrollOffsetY++;
        } else if (editor->cursorLine < editor->scrollOffsetY + GAP && editor->scrollOffsetY > 0) {
            editor->scrollOffsetY--;
        }
    }

    RenderBuffer(*editor);
    RenderCursor(*editor);
    RenderStatBar(*editor);
}

// there's a 14 line limit btw