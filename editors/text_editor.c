#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include "text_editor.h"
#include "text_buffer.h"

#include <string.h>

#define VISIBLE_LINES (14 + 1)
#define CURSOR_SCROLL_GAP 2
#define GAP 2

int fontSize = 5;

void RenderLine(const char* str, int index) {
    int navY = GetNavHeight();

    int y = index * fontSize + GAP * (index + 1) + navY;

    if (index % 2) {
        DrawRectangle(0, y - 1, GetScreenWidth(), fontSize + GAP, Fade(systemPalette[1], 1));
    }

    RenderString(str, GAP, y);
};

struct LineCollection retrieveAllLines(const struct TextBuffer* buffer, int y_index) {
    struct LineCollection result = { .count = 0 };
    const char* work_buffer = buffer->buffer;
    const int maxLines = 14;

    const char* lineStart = work_buffer;
    int currentLine = 0;

    for (unsigned int i = 0; ; i++) {
        char c = work_buffer[i];

        if (c == '\n' || c == '\0') {
            if (currentLine >= y_index && result.count < maxLines) {
                result.lines[result.count] = lineStart;
                result.count++;
            }

            currentLine++;
            lineStart = &work_buffer[i + 1];

            if (c == '\0') break;
        }
    }
    return result;
}

size_t retrieveLinesCount(struct TextBuffer buffer) {
    const struct LineCollection lines = retrieveAllLines(&buffer, 0);
    return lines.count;
}

void RenderBuffer(const struct TextBuffer buffer, int y_index) {
    struct LineCollection lc = retrieveAllLines(&buffer, y_index);

    for (int i = 0; i < lc.count; i++) {
        const char* start = lc.lines[i];
        const char* end = start;
        while (*end != '\n' && *end != '\0') end++;

        unsigned long len = end - start;
        char temp[1025]; // max 1024 chars per line
        if (len >= sizeof(temp)) len = sizeof(temp) - 1;
        memcpy(temp, start, len);
        temp[len] = '\0';

        RenderLine(temp, i);
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
    RenderBuffer(editor->buffer, editor->scrollOffsetY);

    pressedKey = GetCharPressed();
    if (IsKeyPressed(KEY_RIGHT)) {
        if (editor->cursorChar < editor->buffer.totalChars) {
            editor->cursorChar++;
        }
    }
    if (IsKeyPressed(KEY_LEFT)) {
        if (editor->cursorChar > 0) {
            editor->cursorChar--;
        }
    }
    if (IsKeyPressed(KEY_DOWN)) {
        if (editor->cursorLine < retrieveLinesCount(editor->buffer)) {
            editor->cursorLine++;
        }
    }
    if (IsKeyPressed(KEY_ENTER)) {
        pressedKey = '\n';
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
    }
}

// there's a 14 line limit btw