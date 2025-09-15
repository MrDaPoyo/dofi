#include "raylib.h"
#include "../colors.h"
#include "editors.h"
#include "../fonts.h"
#include "../display.h"
#include <stdio.h>
#include "text_editor.h"
#include "text_buffer.h"

#include <ctype.h>
#include <string.h>
#include <math.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

#define VISIBLE_LINES 14
#define CURSOR_SCROLL_FONT_GAP 1

void RenderLine(char* str, int index, size_t realIndex) {
    int navY = GetNavHeight();

    int y = index * FONT_HEIGHT + FONT_GAP * (index + 1) + navY;

    if (realIndex % 2) {
        DrawRectangle(0, y - 1, GetScreenWidth(), FONT_HEIGHT + FONT_GAP, systemPalette[1]);
    }

    RenderString(str, FONT_GAP, y);
};

void RenderStatBar(struct TextEditor editor) {
    DrawRectangle(0, NAVBAR_HEIGHT + VISIBLE_LINES * (FONT_HEIGHT + FONT_GAP) + 1, GetScreenWidth(), HEIGHT - NAVBAR_HEIGHT - VISIBLE_LINES * (FONT_HEIGHT + FONT_GAP), systemPalette[2]);
    char buffer[64];

    sprintf(buffer, "pos: %li/%li;scroll %li/%li;", editor.cursorChar + 1, editor.cursorLine + 1, editor.scrollOffsetX, editor.scrollOffsetY);
    RenderString(buffer, FONT_GAP / 2, HEIGHT - FONT_HEIGHT - FONT_GAP / 2);
    const float progress = (float)editor.buffer.totalChars / MAX_TOTAL_BUFFER_SIZE * WIDTH;
    DrawRectangle(0, NAVBAR_HEIGHT + VISIBLE_LINES * (FONT_HEIGHT + FONT_GAP) + 1, progress, 2, systemPalette[3]);
}

void RenderCursor(struct TextEditor editor) {
    int navY = GetNavHeight();
    const size_t lineIndex = editor.cursorLine - editor.scrollOffsetY;

    int cellWidth = FONT_WIDTH + FONT_SPACING;
    int x = FONT_SPACING + (editor.cursorChar - editor.scrollOffsetX) * cellWidth + 1;
    int y = lineIndex * FONT_HEIGHT + FONT_GAP * (lineIndex + 1) + navY;

    Color cursorColor = systemPalette[2];
    if (fmod(GetTime(), 1.0) < 0.5)
        cursorColor = systemPalette[3];

    DrawRectangle(x, y, FONT_WIDTH, FONT_HEIGHT, cursorColor);
}

struct LineCollection retrieveAllLines(const struct TextBuffer* buffer, size_t y_index) {
    struct LineCollection result = { .count = 0 };
    const char* work_buffer = buffer->buffer;

    const char* lineStart = work_buffer;
    size_t currentLine = 0;

    for (unsigned int i = 0; ; i++) {
        char c = work_buffer[i];

        if (c == '\n' || c == '\0') {
            if (currentLine >= y_index) {
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

    int maxLines = (lc.count > VISIBLE_LINES) ? VISIBLE_LINES : lc.count;

    for (int i = 0; i < maxLines; i++) {
        const char* line = lc.lines[i];
        if (line == NULL)
            continue;

        const char* line_end = line;
        while (*line_end != '\n' && *line_end != '\0')
            line_end++;

        size_t line_len = line_end - line;

        size_t current_x_index = x_index;
        if (current_x_index > line_len)
            current_x_index = line_len;

        const char* start = line + current_x_index;
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

int LoadEntireFile(const char* path, char** outBuf) {
    *outBuf = NULL;
    FILE* f = fopen(path, "rb");
    if (!f) return false;
    fseek(f, 0, SEEK_END);
    long sz = ftell(f);
    if (sz < 0) { fclose(f); return false; }
    rewind(f);
    char* buf = (char*)malloc((size_t)sz + 1);
    if (!buf) { fclose(f); return false; }
    size_t rd = fread(buf, 1, (size_t)sz, f);
    fclose(f);
    buf[rd] = '\0';

    char* src = buf;
    char* dst = buf;
    while (*src) {
        if (*src != '\r')
            *dst++ = *src;
        src++;
    }
    *dst = '\0';

    *outBuf = buf;
    return true;
}

void InitTextEditors(void) {
    editors[0] = (struct TextEditor){
        .buffer = createBuffer(1024),
        .cursorChar = 0,
        .cursorLine = 0,
        .scrollOffsetX = 0,
        .scrollOffsetY = 0,
    };

    char* tempString = NULL;
    if (LoadEntireFile("assets/examples/gradient.lua", &tempString)) {
        for (int i = strlen(tempString); i > 0; i--) {
            char c = tempString[i - 1];
            insertCharacter(&editors[0], c);
        }
        free(tempString);
    }


    for (size_t i = 1; i < sizeof(editors) / sizeof(editors[0]); i++) {
        editors[i] = (struct TextEditor){
            .buffer = createBuffer(0),
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

// FREE THE RETURNED VALUE FROM THIS FUNCTIOn, YOU SON OF A LOTUS BISCOFF BISCUIT!!!! :333 >:D --poyo
char* RetrieveAllCodeFromTextEditors(struct TextEditor editors[TOTAL_TEXT_EDITORS]) {
    size_t needed = 1;
    for (int i = 0; i < TOTAL_TEXT_EDITORS; i++) {
        needed += editors[i].buffer.bufferSize;
    }

    char* tempBuf = malloc(needed + 1);
    size_t index = 0;

    for (int i = 0; i < TOTAL_TEXT_EDITORS; i++) {
        for (size_t j = 0; j < editors[i].buffer.totalChars; j++) {
            tempBuf[index] = editors[i].buffer.buffer[j];
            index++;
        }
        tempBuf[index] = '\n';
        index++;
    }
    tempBuf[index] = '\0';

    return tempBuf;
}


char pressedKey;
bool hasCodeChanged = false;

void RenderTextEditor(void) {
    if (!displayNavbar) {
        ShowNavbar();
    }

    static float keyRepeatDelay = 0.3f;
    static float keyRepeatInterval = 0.05f;
    static float keyTimer = 0.0f;
    static int keyPressed = 0;

    float deltaTime = GetFrameTime();
    pressedKey = GetCharPressed();

    int currentKey = 0;
    if (IsKeyDown(KEY_RIGHT))
        currentKey = KEY_RIGHT;
    else if (IsKeyDown(KEY_LEFT))
        currentKey = KEY_LEFT;
    else if (IsKeyDown(KEY_DOWN))
        currentKey = KEY_DOWN;
    else if (IsKeyDown(KEY_UP))
        currentKey = KEY_UP;
    else if (IsKeyDown(KEY_BACKSPACE))
        currentKey = KEY_BACKSPACE;
    else if (IsKeyDown(KEY_TAB))
        currentKey = KEY_TAB;
    else {
        keyPressed = 0;
        keyTimer = 0.0f;
    }

    if (currentKey != 0) {
        if (currentKey != keyPressed) {
            keyPressed = currentKey;
            keyTimer = keyRepeatDelay;
        } else {
            keyTimer -= deltaTime;
            if (keyTimer > 0.0f) {
                currentKey = 0;
            } else {
                keyTimer = keyRepeatInterval;
            }
        }
    }

    if (currentKey == KEY_RIGHT) {
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
    } else if (currentKey == KEY_LEFT) {
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
    } else if (currentKey == KEY_DOWN) {
        if (editor->cursorLine < editor->buffer.totalLines) {
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
    } else if (currentKey == KEY_UP) {
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
    } else if (currentKey == KEY_BACKSPACE) {
        *editor = removeCharacter(editor);
        hasCodeChanged = true;
    } else if (currentKey == KEY_TAB) {
        for (int i = 0; i < 4; i++) {
            insertCharacter(editor, ' ');
            editor->cursorChar++;
        }
    }

    if (IsKeyPressed(KEY_ENTER)) {
        pressedKey = '\n';
    }
    if (IsKeyPressed(KEY_HOME)) {
        editor->cursorChar = 0;
        editor->scrollOffsetX = 0;
    }
    if (IsKeyPressed(KEY_END)) {
        editor->cursorChar = strlen(GetLineText(editor->buffer, editor->cursorLine));
        if (editor->cursorChar < editor->scrollOffsetX) {
            editor->scrollOffsetX = editor->cursorChar;
        } else if (editor->cursorChar >= editor->scrollOffsetX + LINE_CHAR_WIDTH) {
            editor->scrollOffsetX = editor->cursorChar - LINE_CHAR_WIDTH + 1;
        }
    }

    if (pressedKey != 0 && pressedKey != (char)KEY_TAB) {
        pressedKey = tolower(pressedKey);
        *editor = insertCharacter(editor, pressedKey);
        hasCodeChanged = true;
        if (pressedKey == '\n') {
            editor->cursorChar = 0;
            editor->cursorLine++;
        } else {
            editor->cursorChar++;
        }
    }

    size_t scrollX = editor->scrollOffsetX;
    size_t scrollY = editor->scrollOffsetY;
    size_t X = editor->cursorChar;
    size_t Y = editor->cursorLine;

    size_t viewWidth = LINE_CHAR_WIDTH;
    size_t viewHeight = VISIBLE_LINES;

    if (X < scrollX) {
        scrollX = X;
    } else if (X >= scrollX + viewWidth) {
        scrollX = X - viewWidth + 1;
    }

    if (Y < scrollY) {
        scrollY = Y;
    } else if (Y >= scrollY + viewHeight) {
        scrollY = Y - viewHeight + 1;
    }

    editor->scrollOffsetX = scrollX;
    editor->scrollOffsetY = scrollY;

    RenderBuffer(*editor);
    RenderCursor(*editor);
    RenderStatBar(*editor);
}

// there's a 14 line limit btw