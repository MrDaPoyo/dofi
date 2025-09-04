#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include "text_editor.h"
#include "text_buffer.h"

#define VISIBLE_LINES 14
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

struct TextEditor editors[10];
struct TextEditor* editor = &editors[0];

void InitTextEditors(void) {
    for (size_t i = 0; i < sizeof(editors) / sizeof(editors[0]); i++) {
        editors[i] = (struct TextEditor){
            .buffer = createBuffer(0),
            .cursorChar = 0,
            .cursorLine = 0,
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
    RenderLine(editor->buffer.buffer, 1);
    pressedKey = GetCharPressed();
    if (pressedKey != 0) {
        appendCharacter(editor->buffer, pressedKey);
        pressedKey = GetCharPressed();
    }
}

// there's a 14 line limit btw