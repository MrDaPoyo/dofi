#include "raylib.h"
#include "fonts.h"
#include "colors.h"
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

bool IsComment(const char* string) {
    const size_t len = strlen(string);
    char* tempBuffer = malloc(len + 1);
    if (!tempBuffer) return false;

    strcpy(tempBuffer, string);
    trim(tempBuffer);

    bool isComment = tempBuffer[0] == '-' && tempBuffer[1] == '-';
    free(tempBuffer);

    return isComment;
}

void RenderString(char* str, int x, int y) {
    int offsetX = 0;

    char buf[256]; // temp buf for single words

    bool isComment = IsComment(str);

    while (*str) {
        int i = 0;
        while (*str && *str != ' ' && i < (int)(sizeof(buf) - 1)) {
            buf[i++] = *str++;
        }
        buf[i] = '\0';

        if (!isComment) {
            bool isKeyword = IsReservedKeyword(buf);

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
                    isKeyword ? systemPalette[2] : systemPalette[4] // Orange / White
                    );
                offsetX += FONT_SPACING + FONT_WIDTH;
            }
        } else {
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
                    systemPalette[6] // Grey
                    );
                offsetX += FONT_SPACING + FONT_WIDTH;
            }
        }

        if (*str == ' ') {
            ;
            offsetX += FONT_SPACING + FONT_WIDTH;
            str++;
        }
    }
}
