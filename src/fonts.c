#include "raylib.h"
#include "fonts.h"
#include "colors.h"
#include "display.h"

#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include "editors/text_buffer.h"

Font GeneralFont;
const char* ReservedKeywords[] = { "and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while" };

void LoadFonts(void) {
    if (FileExists("assets/fonts/dofi.ttf")) {
        GeneralFont = LoadFont("assets/fonts/dofi.ttf");
    } else {
        GeneralFont = GetFontDefault();
    }
}

void UnloadFonts(void) {
    UnloadFont(GeneralFont);
}

void trim(char* s) {
    int i = 0, j = 0;
    while (s[i] == ' ')
        i++;
    while ((s[j++] = s[i++]));
}

bool IsReservedKeyword(const char* word) {
    for (size_t i = 0; i < sizeof(ReservedKeywords) / sizeof(ReservedKeywords[0]); i++) {
        if (strcmp(word, ReservedKeywords[i]) == 0) {
            return true;
        }
    }
    return false;
}

int IsComment(const char* line) {
    if (line == NULL) return -1;
    if (*line == '\0' || *line == '\n') return -1;

    for (int i = 1; line[i] != '\0' && line[i] != '\n'; i++) {
        if (line[i - 1] == '-' && line[i] == '-') {
            return i - 1; // index of first '-'
        }
    }
    return -1;
}

void RenderTextEditorString(char* str, size_t length, int scrollX, int y) {
    if (!str || length == 0) return;

    int offsetX = 0;
    char buf[256];
    const int commentPos = IsComment(str);
    size_t pos = 0;
    long int index = 0;

    int baseX = FONT_GAP + -(FONT_WIDTH + FONT_SPACING) * scrollX;
    y = y * FONT_HEIGHT + FONT_GAP * (y + 1) + NAVBAR_HEIGHT;
    length += scrollX + FONT_GAP;

    if (y % 2) {
        DrawRectangle(0, y  - 1, GetScreenWidth(), FONT_HEIGHT + FONT_GAP, systemPalette[1]);
    }

    while (pos < length && str[pos] != '\0') {
        int i = 0;
        while (pos < length && str[pos] != '\0' && str[pos] != ' ' && i < (int)(sizeof(buf) - 1)) {
            buf[i++] = str[pos++];
            index++;
        }
        buf[i] = '\0';

        bool inComment = (commentPos >= 0 && index > commentPos);

        for (int j = 0; j < i; j++) {
            char c[2] = { buf[j], '\0' };
            Vector2 drawPos = (Vector2){ (float)(baseX + offsetX), (float)y };
            DrawTextPro(
                GeneralFont,
                c,
                drawPos,
                (Vector2){ 0, 0 },
                0.0f,
                FONT_HEIGHT,
                FONT_SPACING,
                inComment ? systemPalette[6] : IsReservedKeyword(buf) ? systemPalette[2] : systemPalette[4]
            );
            offsetX += FONT_SPACING + FONT_WIDTH;
        }

        if (pos < length && str[pos] == ' ') {
            offsetX += FONT_SPACING + FONT_WIDTH;
            pos++;
            index++;
        }
    }
}

void RenderString(char* str, int x, int y) {
    int offsetX = 0;
    char buf[256]; // temp buf for single words

    const int commentPos = IsComment(str); // -1 if no comment
    long int index = 0;

    while (*str) {
        int i = 0;

        while (*str && *str != ' ' && i < (int)(sizeof(buf) - 1)) {
            buf[i++] = *str++;
            index++;
        }
        buf[i] = '\0';

        bool inComment = (commentPos >= 0 && index > commentPos);

        for (int j = 0; j < i; j++) {
            char c[2] = { buf[j], '\0' };
            Vector2 pos = (Vector2){ (float)(x + offsetX), (float)y };
            DrawTextPro(
                GeneralFont,
                c,
                pos,
                (Vector2){ 0, 0 },
                0.0f,
                FONT_HEIGHT,
                FONT_SPACING,
                inComment ? systemPalette[6] : IsReservedKeyword(buf) ? systemPalette[2] : systemPalette[4] // Grey. Orange, White.
            );
            offsetX += FONT_SPACING + FONT_WIDTH;
        }

        if (*str == ' ') {
            offsetX += FONT_SPACING + FONT_WIDTH;
            str++;
            index++;
        }
    }
}

void RenderStringWrap(char* str, int x, int y) {
    int offsetX = 0;
    int offsetY = 0;
    char buf[256]; // temp buffer for single words

    const int commentPos = IsComment(str); // -1 if no comment
    long int index = 0;

    while (*str) {
        int i = 0;

        while (*str && *str != ' ' && i < (int)(sizeof(buf) - 1)) {
            buf[i++] = *str++;
            index++;
        }
        buf[i] = '\0';

        bool inComment = (commentPos >= 0 && index > commentPos);

        if (offsetX + i > LINE_CHAR_WIDTH * (FONT_WIDTH + FONT_SPACING)) {
            offsetX = 0;
            offsetY += FONT_HEIGHT + FONT_SPACING;
        }

        for (int j = 0; j < i; j++) {
            char c[2] = { buf[j], '\0' };
            Vector2 pos = (Vector2){ (float)(x + offsetX), (float)(y + offsetY) };
            DrawTextPro(
                GeneralFont,
                c,
                pos,
                (Vector2){ 0, 0 },
                0.0f,
                FONT_HEIGHT,
                FONT_SPACING,
                inComment ? systemPalette[6] : IsReservedKeyword(buf) ? systemPalette[2] : systemPalette[4]
            );
            offsetX += FONT_WIDTH + FONT_SPACING;
        }

        if (*str == ' ') {
            offsetX += FONT_WIDTH + FONT_SPACING;
            str++;
            index++;
        }

        if (offsetX > LINE_CHAR_WIDTH * (FONT_WIDTH + FONT_SPACING)) {
            offsetX = 0;
            offsetY += FONT_HEIGHT + FONT_SPACING;
        }
    }
}