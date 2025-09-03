#ifndef TEXT_BUFFER_H
#define TEXT_BUFFER_H

#include <stddef.h>

#define MAX_TOTAL_BUFFER_SIZE 65536

struct TextBuffer
{
    char* buffer;
    size_t bufferSize;
    size_t totalChars;
};

struct TextBuffer createBuffer(size_t bufferSize);
void modifyBufferCapacity(struct TextBuffer* buffer, size_t newBufferSize);
void appendCharacter(struct TextBuffer* buffer, char character);

#endif