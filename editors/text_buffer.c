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

struct TextBuffer checkIfMoreSpaceIsNeeded(struct TextBuffer buffer) {
    if (buffer.totalChars + 2 >= buffer.bufferSize) {
        const size_t availableBufferSize = retrieveAvailableBufferSize();

        if (buffer.bufferSize + BUFFER_SIZE_INCREMENT <= availableBufferSize) {
            buffer = modifyBufferCapacity(buffer, buffer.bufferSize + BUFFER_SIZE_INCREMENT);
        } else if (availableBufferSize > buffer.bufferSize) {
            buffer = modifyBufferCapacity(buffer, availableBufferSize);
        }
    }
    return buffer;
}

struct TextBuffer appendCharacter(struct TextBuffer buffer, char character) {
    buffer = checkIfMoreSpaceIsNeeded(buffer);

    buffer.buffer[buffer.totalChars] = character;
    buffer.buffer[buffer.totalChars + 1] = '\0';
    buffer.totalChars++;

    return buffer;
}

struct TextBuffer insertCharacter(struct TextBuffer buffer, char character, size_t position) {
    buffer = checkIfMoreSpaceIsNeeded(buffer);

    // this just shifts the chars so that i can make room for the inserted char
    for (size_t i = buffer.totalChars; i > position; i--) {
        buffer.buffer[i] = buffer.buffer[i - 1];
    }

    buffer.buffer[position] = character;
    buffer.totalChars++;
    buffer.buffer[buffer.totalChars] = '\0';

    return buffer;
}

