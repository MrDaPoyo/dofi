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

struct TextBuffer checkBufferSpace(struct TextBuffer buffer) {
    const size_t availableBufferSize = retrieveAvailableBufferSize();

    if (buffer.totalChars + 2 >= buffer.bufferSize) {
        if (buffer.bufferSize + BUFFER_SIZE_INCREMENT <= availableBufferSize) {
            buffer = modifyBufferCapacity(buffer, buffer.bufferSize + BUFFER_SIZE_INCREMENT);
        } else if (availableBufferSize > buffer.bufferSize) {
            buffer = modifyBufferCapacity(buffer, availableBufferSize);
        }
    }

    return buffer;
}

struct TextBuffer appendCharacter(struct TextBuffer buffer, char character) {
    buffer = checkBufferSpace(buffer);

    buffer.buffer[buffer.totalChars] = character;
    buffer.buffer[buffer.totalChars + 1] = '\0';
    buffer.totalChars++;

    return buffer;
}

size_t getIndexFromEditorBuffer(struct TextEditor* editor) {
    const char* work_buffer = editor->buffer.buffer;

    const size_t row = editor->cursorLine;
    const size_t col = editor->cursorChar;

    size_t currentIndex = 0;
    size_t currentLine = 0;

    while (work_buffer[currentIndex] != '\0') {
        if (currentLine == row) {
            size_t lineStart = currentIndex;

            for (size_t i = 0; ; i++) {
                char c = work_buffer[lineStart + i];
                if (c == '\n' || c == '\0') {
                    return lineStart + i;
                }
                if (i == col) {
                    return lineStart + i;
                }
            }
        }

        if (work_buffer[currentIndex] == '\n') {
            currentLine++;
        }
        currentIndex++;
    }

    return currentIndex;
}

struct TextEditor insertCharacter(struct TextEditor* editor, char character) {
    struct TextBuffer buffer = checkBufferSpace(editor->buffer);
    const size_t position = getIndexFromEditorBuffer(editor);

    // this just shifts the chars so that i can make room for the inserted char
    for (size_t i = buffer.totalChars; i > position; i--) {
        buffer.buffer[i] = buffer.buffer[i - 1];
    }

    buffer.buffer[position] = character;
    buffer.totalChars++;
    buffer.buffer[buffer.totalChars] = '\0';

    editor->buffer = buffer;

    return *editor;
}

struct TextEditor removeCharacter(struct TextEditor* editor) {
    struct TextBuffer buffer = checkBufferSpace(editor->buffer);
    size_t position = getIndexFromEditorBuffer(editor);

    if (position < buffer.totalChars) {
        for (size_t i = position; i < buffer.totalChars - 1; i++) {
            buffer.buffer[i] = buffer.buffer[i + 1];
        }
        buffer.buffer[buffer.totalChars - 1] = '\0';
        buffer.totalChars--;
    }

    editor->buffer = buffer;
    return *editor;
}

// ALWAYS REMEMBER TO FREE UP THIS FUNCTION'S RETURNED CHAR
char* GetLineText(struct TextBuffer buffer, size_t lineIndex) {
    size_t currentLine = 0;
    size_t startIndex = 0;
    size_t endIndex = 0;

    for (size_t i = 0; buffer.buffer[i] != '\0'; i++) {
        if (currentLine == lineIndex) {
            startIndex = i;
            break;
        }
        if (buffer.buffer[i] == '\n') {
            currentLine++;
        }
    }

    if (currentLine != lineIndex) {
        return NULL;
    }

    for (endIndex = startIndex;
         buffer.buffer[endIndex] != '\0' && buffer.buffer[endIndex] != '\n';
         endIndex++);

    size_t lineLength = endIndex - startIndex;
    char* lineText = malloc(lineLength + 1);
    if (!lineText) return NULL;

    int index = 0;
    for (size_t i = 0; i < lineLength; i++) {
        if (buffer.buffer[startIndex + i] != '\n') {
            lineText[i] = buffer.buffer[startIndex + index];
            index++;
        }
    }
    lineText[lineLength] = '\0';

    return lineText;
}