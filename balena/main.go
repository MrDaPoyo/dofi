package main

import (
    "fmt"
    "os"
	scanner "github.com/mrdapoyo/dofi/balena/scanner"
	parser "github.com/mrdapoyo/dofi/balena/parser"
	"github.com/mrdapoyo/dofi/balena/token"
	ast "github.com/mrdapoyo/dofi/balena/tool"
)

func main() {
    args := os.Args[1:]
    if len(args) > 1 {
        fmt.Println("Usage: go run . [script]")
        os.Exit(64)
    } else if len(args) == 1 {
        runFile(args[0])
    } else {
        runPrompt()
    }
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(65)
	}
	if hadRuntimeError {
		fmt.Fprintln(os.Stderr, "Runtime error occurred. Exiting.")
		os.Exit(70)
	}
	run(string(bytes))
}

func runPrompt() {
	reader := os.Stdin
	buf := make([]byte, 1024)
	for {
		fmt.Print("> ")
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		line := string(buf[:n])
		if len(line) == 0 {
			break
		}
		run(line)
	}
}

func run(script string) {
	scanner := scanner.NewScanner(script)
	tokens := scanner.ScanTokens()

	parser := parser.NewParser(tokens)
	expression := parser.Parse()

	if hadError {
		return
	}

	fmt.Println(ast.AstPrinter(expression))
}

var hadError = false
var hadRuntimeError = false;

func ErrorToken(badToken token.Token, message string) {
	if badToken.Type == token.EOF {
		Report(badToken.Line, " at end", message)
	} else {
		Report(badToken.Line, " at '"+badToken.Lexeme+"'", message)
	}
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}