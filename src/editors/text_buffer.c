#include "text_buffer.h"
#include "text_editor.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

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
        .totalLines = 0,
    };
    printf("Buffer created size: %li \n", bufferSize);
    return newTextBuffer;
};

struct TextBuffer modifyBufferCapacity(struct TextBuffer buffer, size_t newBufferSize) {
    char* newBuffer = malloc(newBufferSize);

    size_t copyLen = buffer.totalChars < newBufferSize - 1 ? buffer.totalChars : newBufferSize - 1;
    memcpy(newBuffer, buffer.buffer, copyLen);

    newBuffer[copyLen] = '\0';

    free(buffer.buffer);
    buffer.buffer = newBuffer;
    buffer.bufferSize = newBufferSize;
    buffer.totalChars = copyLen;

    return buffer;
}

struct TextBuffer checkBufferSpace(struct TextBuffer buffer) {
    const size_t availableBufferSize = retrieveAvailableBufferSize();

    if (buffer.totalChars + buffer.totalLines + 10 >= buffer.bufferSize) {
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

    if (retrieveAvailableBufferSize() > 2) {
        // this just shifts the chars so that i can make room for the inserted char
        for (size_t i = buffer.totalChars; i > position; i--) {
            buffer.buffer[i] = buffer.buffer[i - 1];
        }

        buffer.buffer[position] = character;
        buffer.totalChars++;
        buffer.buffer[buffer.totalChars] = '\0';

        if (character == '\n') {
            buffer.totalLines++;
        }
        editor->buffer = buffer;
    }

    return *editor;
}

struct TextEditor removeCharacter(struct TextEditor* editor) {
    struct TextBuffer buffer = checkBufferSpace(editor->buffer);

    size_t index = getIndexFromEditorBuffer(editor);

    if (index == 0) {
        editor->buffer = buffer;
        return *editor;
    }

    if (index > buffer.totalChars) index = buffer.totalChars;

    size_t delPos = index - 1;

    char deleted = buffer.buffer[delPos];

    size_t prevLineLen = 0;
    if (deleted == '\n') {
        size_t prevStart = delPos;
        while (prevStart > 0 && buffer.buffer[prevStart - 1] != '\n') {
            prevStart--;
        }
        prevLineLen = delPos - prevStart;
    }

    for (size_t i = delPos; i + 1 < buffer.totalChars; i++) {
        buffer.buffer[i] = buffer.buffer[i + 1];
    }

    if (buffer.totalChars > 0) {
        buffer.totalChars--;
        buffer.buffer[buffer.totalChars] = '\0';
    }

    if (deleted == '\n') {
        if (editor->cursorLine > 0) {
            editor->cursorLine--;
            editor->cursorChar = prevLineLen;
        } else {
            editor->cursorChar = 0;
        }
    } else {
        if (editor->cursorChar > 0) editor->cursorChar--;
        else editor->cursorChar = 0;
    }

    if (deleted == '\n') {
        buffer.totalLines--;
    }
    editor->buffer = buffer;
    return *editor;
}


// ALWAYS REMEMBER TO FREE THIS FUNCTION'S RETURNED CHAR
char* GetLineText(struct TextBuffer buffer, size_t lineIndex) {
    size_t currentLine = 0;
    size_t startIndex = 0;

    if (lineIndex == 0) {
        startIndex = 0;
    } else {
        for (size_t i = 0; buffer.buffer[i] != '\0'; i++) {
            if (buffer.buffer[i] == '\n') {
                currentLine++;
                if (currentLine == lineIndex) {
                    startIndex = i + 1;
                    break;
                }
            }
        }
        if (currentLine != lineIndex) return NULL;
    }

    size_t endIndex = startIndex;
    while (buffer.buffer[endIndex] != '\0' && buffer.buffer[endIndex] != '\n') {
        endIndex++;
    }

    size_t lineLength = endIndex - startIndex;

    char* lineText = (char*)malloc(lineLength + 1);
    if (!lineText) return NULL;

    if (lineLength > 0) {
        memcpy(lineText, buffer.buffer + startIndex, lineLength);
    }
    lineText[lineLength] = '\0';
    return lineText;
}
