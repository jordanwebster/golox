package loxerror

import (
	"fmt"

	"github.com/jordanwebster/golox/token"
)

var hadError = false

func ParseError(t token.Token, message string) {
	if t.Type == token.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error ", where, ": ", message)
	hadError = true
}

func ClearError() {
	hadError = false
}

func HadError() bool {
	return hadError
}
