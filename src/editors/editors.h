#ifndef DOFI_EDITORS_H
#define DOFI_EDITORS_H

#include "../fonts.h"

#include <stddef.h>

void RenderTextEditor(void);
void RenderSpriteEditor(void);
void RenderMapEditor(void);
void RenderSoundEditor(void);
void RenderPlayEditor(void);

void RenderLine(char *str, int index, size_t realIndex);

// TEXT EDITOR STUFF
extern const char *testString;

#endif
