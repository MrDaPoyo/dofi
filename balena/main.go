package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			fmt.Fprintln(out, "Goodbye!")
			return
		}

		l := NewLexer(line)
		p := NewParser(l)
		program := p.ParseProgram()

		if errors := p.Errors(); len(errors) != 0 {
			printParserErrors(out, errors)
			continue
		}

		fmt.Fprintln(out, program.String())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "Parser errors:")
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}

func main() {
	fmt.Printf("Welcome to the Balena programming language REPL!\n")
	fmt.Printf("Type 'exit' or 'quit' to exit.\n")
	Start(os.Stdin, os.Stdout)
}