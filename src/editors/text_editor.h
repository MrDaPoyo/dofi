#ifndef TEXT_EDITOR_H
#define TEXT_EDITOR_H

#include "text_buffer.h"
#include "stdbool.h"

struct TextEditor
{
    struct TextBuffer buffer;
    size_t cursorLine;
    size_t cursorChar;
    size_t scrollOffsetY;
    size_t scrollOffsetX;
};

struct LineCollection {
    const char* lines[14];
    int count;
};

#define TOTAL_TEXT_EDITORS 10

extern struct TextEditor editors[TOTAL_TEXT_EDITORS];
extern struct TextEditor* editor;

extern bool hasCodeChanged;

void InitTextEditors(void);
char* RetrieveAllCodeFromTextEditors(struct TextEditor editor[TOTAL_TEXT_EDITORS]);
struct TextEditor createTextEditor(size_t startingBufferSize);

#endif
