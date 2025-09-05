#ifndef TEXT_EDITOR_H
#define TEXT_EDITOR_H

#include "text_buffer.h"

struct TextEditor
{
    struct TextBuffer buffer;
    size_t cursorLine;
    size_t cursorChar;
};

#define TOTAL_TEXT_EDITORS 10

extern struct TextEditor editors[TOTAL_TEXT_EDITORS];
extern struct TextEditor* editor;

void InitTextEditors(void);
struct TextEditor createTextEditor(size_t startingBufferSize);

#endif
