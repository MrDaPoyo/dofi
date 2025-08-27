NAME = dofi

CC = gcc

SRC = main.c

BIN = $(NAME)

CFLAGS = -Wall -Wextra -O2

ifeq ($(OS),Windows_NT)
	LIBS = -lraylib -lopengl32 -lgdi32 -lwinmm
	RUN = .\$(BIN).exe
else
	LIBS = -lraylib -lm -ldl -lpthread -lGL -lX11
	RUN = ./$(BIN)
endif

all: $(BIN)

$(BIN): $(SRC)
	$(CC) $(CFLAGS) -o $(BIN) $(SRC) $(LIBS)

run: $(BIN)
	$(RUN)

clean:
	-rm -f $(BIN) $(BIN).exe

.PHONY: all run clean
