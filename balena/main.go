package main

import (
	"bufio"
	"fmt"
	"os"

	"runtime/debug"

	interpreter "github.com/mrdapoyo/dofi/balena/interpreter"
	parser "github.com/mrdapoyo/dofi/balena/parser"
	scanner "github.com/mrdapoyo/dofi/balena/scanner"
	"github.com/mrdapoyo/dofi/balena/token"
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
	run(string(bytes))
	if hadRuntimeError {
		fmt.Fprintln(os.Stderr, "Runtime error occurred. Exiting.")
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	interp := interpreter.NewInterpreter()
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		if line == "" {
			continue
		}
		runWithInterpreter(interp, line)
		hadError = false
		hadRuntimeError = false
	}
}

func run(script string) {
	s := scanner.NewScanner(script)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	statements := p.Parse()
	if hadError {
		return
	}

	interpreterInstance := interpreter.NewInterpreter()
	fmt.Println("Created interpreter")
	resolver := interpreter.NewResolver(interpreterInstance)
	fmt.Println("Created resolver")
	resolver.Resolve(statements)
	fmt.Println("Resolved statements")

	if hadError {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*interpreter.RuntimeError); ok {
				fmt.Fprintln(os.Stderr, runtimeErr.Error())
				hadRuntimeError = true
			} else if returnErr, ok := r.(interpreter.Return); ok {
				if returnErr.Value != nil {
					fmt.Println(returnErr.Value)
				}
				fmt.Fprintln(os.Stderr, "Unexpected return statement at top level.")
				hadRuntimeError = true
			} else {
				fmt.Fprintf(os.Stderr, "Unexpected panic: %v\n", r)
				debug.PrintStack()
				hadRuntimeError = true
			}
		}
	}()

	interpreterInstance.Execute(statements)
}

func runWithInterpreter(interpreterInstance *interpreter.Interpreter, script string) {
	s := scanner.NewScanner(script)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	statements := p.Parse()
	if hadError {
		return
	}

	resolver := interpreter.NewResolver(interpreterInstance)
	resolver.Resolve(statements)

	if hadError {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*interpreter.RuntimeError); ok {
				fmt.Fprintln(os.Stderr, runtimeErr.Error())
				hadRuntimeError = true
			} else if returnErr, ok := r.(interpreter.Return); ok {
				if returnErr.Value != nil {
					fmt.Println(returnErr.Value)
				}
				hadRuntimeError = false
			} else {
				fmt.Fprintf(os.Stderr, "Unexpected panic: %v\n", r)
				hadRuntimeError = true
			}
		}
	}()

	interpreterInstance.Execute(statements)
}

var hadError = false
var hadRuntimeError = false

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
