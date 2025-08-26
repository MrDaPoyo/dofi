NAME = dofi

CC = gcc

SRC = main.c

BIN = $(NAME)

CFLAGS = -Wall -Wextra -O2

LIBS = -lraylib -lm -ldl -lpthread -lGL -lX11

all: $(BIN)

$(BIN): $(SRC)
	$(CC) $(CFLAGS) -o $(BIN) $(SRC) $(LIBS)

run: $(BIN)
	./$(BIN)

clean:
	rm -f $(BIN)

.PHONY: all run clean
