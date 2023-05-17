package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	_ "github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/loxerror"
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
}

func run(source string) {
	fmt.Println("Executing", source)
	scanner := NewScanner(source)
	for _, token := range scanner.ScanTokens() {
		fmt.Println(token)
	}
}
