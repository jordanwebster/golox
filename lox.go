package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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
        fmt.Print("> ")
    }
}

func runFile(path string) {
    source, err := os.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }
    run(string(source))
}

func run(source string) {
    fmt.Println("Executing", source)
    scanner := NewScanner()
    for _, token := range scanner.ScanTokens() {
        fmt.Println(token)
    }
}
