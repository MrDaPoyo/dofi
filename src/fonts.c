#include "raylib.h"
#include "fonts.h"
#include "colors.h"
#include "display.h"

#include <stddef.h>
#include <stdlib.h>
#include <string.h>

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
    // empty line
    if (*line == '\0' || *line == '\n') {
        return -1;
    }
    // start at the 2nd char
    size_t index = 1;
    const char* current = line + 1;
    // walk until end of string or end of line
    while (*current != '\0' && *current != '\n') {
        const char* prev = current - 1;
        // check for "--"
        if (*prev == '-' && *current == '-') {
            return index;
        }
        ++current;
        ++index;
    }
    return -1;
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