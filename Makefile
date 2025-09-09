NAME = dofi

CC = gcc

SRC = $(wildcard src/*.c src/editors/*.c)

BIN = $(NAME)

CFLAGS = -Wall -Wextra -O3 -g -Werror -Wno-error=unused

ASSETS = src/assets

OUTPUT_DIR = out

ifeq ($(OS),Windows_NT)
	LIBS = -lraylib -lopengl32 -lgdi32 -lwinmm
	RUN = .\$(BIN).exe
else
	LIBS = -lraylib -lm -ldl -lpthread -lGL -lX11
	RUN = ./$(BIN)
endif

all: $(BIN)

$(BIN): $(SRC)
	mkdir -p $(OUTPUT_DIR)
	$(CC) $(CFLAGS) -o $(OUTPUT_DIR)/$(BIN) $(SRC) $(LIBS)
	cp -r $(ASSETS) $(OUTPUT_DIR)

run: $(BIN)
	cd $(OUTPUT_DIR) && $(RUN)

clean:
	rm -f $(BIN) $(BIN).exe
	rm -rf $(OUTPUT_DIR)

.PHONY: all run clean
