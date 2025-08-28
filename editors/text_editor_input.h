#ifndef TEXT_EDITOR_INPUT_H
#define TEXT_EDITOR_INPUT_H

#include <stddef.h>

#define VISIBLE_LINES 14

extern char textBuffer[4096];
extern int cursorLine;
extern int cursorChar;
extern int scrollLineOffset;
extern int scrollCharOffset;

const char *GetLineText(const char *str, int index);

void TextEditor_HandleInput(void);
void ClampScrollOffsets(void);
void EnsureCursorVisible(void);

int GetTotalLines(const char *s);
int LineLength(const char *s, int line);
int FindAbsolutePos(const char *s, int line, int charIndex);

#endif
