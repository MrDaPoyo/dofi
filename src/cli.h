#ifndef DOFI_CLI_H
#define DOFI_CLI_H

#include "raylib.h"

void AppendCLILine(const char* str);

struct CLILine {
    char* str;
    int len;
    Color color;
};

#endif //DOFI_CLI_H