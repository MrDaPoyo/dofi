#ifndef TEXT_BUFFER_H
#define TEXT_BUFFER_H

#include <stddef.h>

#define MAX_TOTAL_BUFFER_SIZE 65536
#define BUFFER_SIZE_INCREMENT 1024

struct TextBuffer
{
    char* buffer;
    size_t bufferSize;
    size_t totalChars;
};

struct TextBuffer createBuffer(size_t bufferSize);
struct TextBuffer modifyBufferCapacity(struct TextBuffer buffer, size_t newBufferSize);
struct TextBuffer appendCharacter(struct TextBuffer buffer, char character);

#endif