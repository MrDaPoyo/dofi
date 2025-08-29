#ifndef TEXT_EDITOR_H
#define TEXT_EDITOR_H

struct TextEditor
{
    const char *buffer;
    int scrollLineOffset;
    int scrollCharOffset;
    int cursorLine;
    int cursorChar;
};

extern struct TextEditor editors[10];

int scrollLineOffset;
int scrollCharOffset;
int cursorLine;
int cursorChar;
int fontSize;

#endif
