#ifndef DOFI_EDITORS_H
#define DOFI_EDITORS_H

#include "../fonts.h"

void RenderTextEditor(void);
void RenderSpriteEditor(void);
void RenderMapEditor(void);
void RenderSoundEditor(void);
void RenderPlayEditor(void);

void RenderLine(const char *str, int index);

// TEXT EDITOR STUFF
extern const char *testString;

#endif
