NAME = dofi

CC = gcc

SRC = $(wildcard src/*.c src/editors/*.c)

BIN = $(NAME)

CFLAGS = -Wall -Wextra -O3 -g -Werror -Wno-error=unused

ASSETS = src/assets

OUTPUT_DIR = .

ifeq ($(OS),Windows_NT)
	LIBS = -lraylib -lopengl32 -lgdi32 -lwinmm -llua5.4
	EXT = .exe
else
	LIBS = -lraylib -lm -ldl -lpthread -lGL -lX11 -llua5.4
	EXT =
endif

all: $(BIN)
	cd $(OUTPUT_DIR)

$(BIN): $(SRC)
	mkdir -p $(OUTPUT_DIR)
	$(CC) $(CFLAGS) -o $(OUTPUT_DIR)/$(BIN)$(EXT) $(SRC) $(LIBS)
	cp -r $(ASSETS) $(OUTPUT_DIR)

run: $(BIN)
	cd $(OUTPUT_DIR) && ./$(BIN)$(EXT)

clean:
	rm -f $(BIN) $(BIN).exe

.PHONY: all run clean
