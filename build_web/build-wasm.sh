#!/bin/bash
set -euo pipefail

source ~/emsdk/emsdk_env.sh

shopt -s globstar

mkdir -p ../out

SRC_FILES=(../src/**/*.c)

emcc "${SRC_FILES[@]}" \
    -I../src \
    -I../raylib/src \
    -I../src/lua \
    ../raylib/src/libraylib.web.a \
    -s USE_GLFW=3 \
    -s ASYNCIFY \
    --shell-file ../shell.html \
    --preload-file ../src/assets@/assets \
    -s TOTAL_STACK=64MB \
    -s INITIAL_MEMORY=128MB \
    -s ASSERTIONS=0 \
    -s WARN_UNALIGNED=0 \
    -s EXIT_RUNTIME=0 \
    -s LLD_REPORT_UNDEFINED=0 \
    -DNDEBUG \
    -DPLATFORM_WEB \
    -o ../out/index.html \
    -Oz \
    -g0 \
    -flto

cp ../src/assets/favicon.ico ../out
cp ../src/assets/fonts/dofi.ttf ../out
