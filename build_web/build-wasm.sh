#!/bin/bash
set -euo pipefail

source ~/emsdk/emsdk_env.sh

shopt -s globstar

emcc ../src/**/*.c \
    -I../src \
    -I../raylib/src \
    ../raylib/src/libraylib.web.a \
    -s USE_GLFW=3 \
    -s ASYNCIFY \
    --shell-file ../shell.html \
    --preload-file ../src/assets@/assets \
    -s TOTAL_STACK=64MB \
    -s INITIAL_MEMORY=128MB \
    -s ASSERTIONS \
    -DPLATFORM_WEB \
    -o ../out/index.html

cd ../out
emrun --no_browser index.html