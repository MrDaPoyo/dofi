#include "text_buffer.h"
#include "text_editor.h"

#include <stdio.h>
#include <stdlib.h>

size_t retrieveUsedBufferSize() {
    size_t totalBufferSize = 0;
    for (unsigned int i = 0; i < TOTAL_TEXT_EDITORS; i++) {
        totalBufferSize += editors[i].buffer.bufferSize;
    }
    return totalBufferSize;
}

size_t retrieveAvailableBufferSize() {
    return MAX_TOTAL_BUFFER_SIZE - retrieveUsedBufferSize();
}

struct TextBuffer createBuffer(size_t bufferSize) {
    const size_t totalBufferSize = retrieveAvailableBufferSize();
    const size_t availableBufferSize = MAX_TOTAL_BUFFER_SIZE - totalBufferSize;
    if (availableBufferSize < bufferSize) {
        bufferSize = availableBufferSize;
    }

    char* newBuffer = malloc(bufferSize);
    newBuffer[0] = '\0';
    const struct TextBuffer newTextBuffer = {
        .buffer = newBuffer,
        .bufferSize = bufferSize,
        .totalChars = 0,
    };
    printf("Buffer created size: %li \n", bufferSize);
    return newTextBuffer;
};

struct TextBuffer modifyBufferCapacity(struct TextBuffer buffer, size_t newBufferSize) {
    char* newBuffer = malloc(newBufferSize);
    for (unsigned int i = 0; i < buffer.bufferSize && i < newBufferSize; i++) {
        newBuffer[i] = buffer.buffer[i];
    }
    newBuffer[0] = '\0';

    free(buffer.buffer);
    buffer.buffer = newBuffer;
    buffer.bufferSize = newBufferSize;
    if (buffer.totalChars > buffer.bufferSize)
        buffer.totalChars = buffer.bufferSize;
    return buffer;
};

struct TextBuffer appendCharacter(struct TextBuffer buffer, char character) {
    if (buffer.totalChars + 1 >= buffer.bufferSize) {
        const size_t availableBufferSize = retrieveAvailableBufferSize();

        if (buffer.bufferSize + BUFFER_SIZE_INCREMENT <= availableBufferSize) {
            buffer = modifyBufferCapacity(buffer, buffer.bufferSize + BUFFER_SIZE_INCREMENT);
        } else if (availableBufferSize > buffer.bufferSize) {
            buffer = modifyBufferCapacity(buffer, availableBufferSize);
        } else {
            return buffer;
        }
    }

    buffer.buffer[buffer.totalChars - 1] = character;
    buffer.buffer[buffer.totalChars] = '\0';
    buffer.totalChars++;

    printf("Buffer Size After Input '%c': %li", character, buffer.bufferSize);

    return buffer;
}
