#ifndef TEXT_EDITOR_H
#define TEXT_EDITOR_H

#include "text_buffer.h"

struct TextEditor
{
    struct TextBuffer buffer;
    int cursorLine;
    int cursorChar;
};

void InitTextEditors(void);

struct TextEditor createTextEditor(size_t startingBufferSize);

#endif
