#ifndef DOFI_EDITORS_H
#define DOFI_EDITORS_H

void RenderTextEditor(void);
void RenderSpriteEditor(void);
void RenderMapEditor(void);
void RenderSoundEditor(void);
void RenderPlayEditor(void);

void RenderLine(const char *str, int y);

void RenderString(const char *str, int x, int y);

#endif
