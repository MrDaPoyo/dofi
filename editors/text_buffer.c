#include "text_buffer.h"

#include <stdio.h>
#include <stdlib.h>

struct TextBuffer createBuffer(size_t bufferSize) {
    const unsigned int newBufferSize = bufferSize <= MAX_TOTAL_BUFFER_SIZE ? bufferSize : MAX_TOTAL_BUFFER_SIZE;
    char* newBuffer = malloc(newBufferSize);
    newBuffer[0] = '\0';
    const struct TextBuffer newTextBuffer = {
        .buffer = newBuffer,
        .bufferSize = newBufferSize,
        .totalChars = 0,
    };
    printf("Buffer created size: %u \n", newBufferSize);
    return newTextBuffer;
};

void modifyBufferCapacity(struct TextBuffer *buffer, size_t newBufferSize)
{
    char *oldBuffer = buffer->buffer;
    char *newBuffer = malloc(newBufferSize + 1);
    for (unsigned int i = 0; i < buffer->bufferSize && i < newBufferSize; i++)
    {
        newBuffer[i] = oldBuffer[i];
    }
    newBuffer[newBufferSize] = '\0';
    free(oldBuffer);
    buffer->buffer = newBuffer;
    buffer->bufferSize = newBufferSize;
    if (buffer->totalChars > buffer->bufferSize)
        buffer->totalChars = buffer->bufferSize - 1;
};
