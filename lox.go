package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/interpreter"
	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/loxio"
	"github.com/jordanwebster/golox/parser"
	"github.com/jordanwebster/golox/scanner"
	"github.com/jordanwebster/golox/token"
)

//go:generate go run ./ast/cmd/gen.go
//go:generate go fmt ./ast

var globalInterpreter *interpreter.Interpreter = interpreter.NewInterpreter()

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
	read_writer := loxio.NewChannelReadWriter()
	go readLines(read_writer)

	tokens := make(chan token.Token)
	scanner := scanner.NewScanner(read_writer, tokens)
	go scanner.ScanTokens()

	statements := make(chan ast.Stmt)
	parser := parser.NewParser(tokens, statements)
	go parser.Parse()

	for stmt := range statements {
		globalInterpreter.Interpret([]ast.Stmt{stmt})
	}

}

func readLines(writer *loxio.ChannelReadWriter) {
	line_scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for line_scanner.Scan() {
		line := line_scanner.Bytes()
		writer.Write(line)
		fmt.Print("> ")
	}

	writer.Close()
}

func runFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	tokens := make(chan token.Token)
	scanner := scanner.NewScanner(strings.NewReader(string(source)), tokens)
	go scanner.ScanTokens()

	statements := make(chan ast.Stmt)
	parser := parser.NewParser(tokens, statements)
	go parser.Parse()

	stmts := make([]ast.Stmt, 0, 64)
	for stmt := range statements {
		stmts = append(stmts, stmt)
	}

	if loxerror.HadError() {
		os.Exit(65)
	}

	globalInterpreter.Interpret(stmts)

	if loxerror.HadRuntimeError() {
		os.Exit(70)
	}
}
