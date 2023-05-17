package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	_ "github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/interpreter"
	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/parser"
)

//go:generate go run ./ast/cmd/gen.go
//go:generate go fmt ./ast

func main() {
	switch numArgs := len(os.Args); numArgs {
	case 1:
		runPrompt()
	case 2:
		runFile(os.Args[1])
	default:
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		run(line)
		loxerror.ClearError()
		fmt.Print("> ")
	}
}

func runFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	run(string(source))

	if loxerror.HadError() {
		os.Exit(65)
	}

    if loxerror.HadRuntimeError() {
        os.Exit(70)
    }
}

func run(source string) {
	scanner := NewScanner(source)
    tokens := scanner.ScanTokens()
    parser := parser.NewParser(tokens) 
    expr := parser.Parse()

    if loxerror.HadError() {
        return
    }
    // TODO: 
    // 1. Move scanner to own module
    // 2. Ensure parse is not swallowing errors
    // 3. Make error module API consisent

    interpreter := interpreter.NewInterpreter()
    interpreter.Interpret(expr)
}
