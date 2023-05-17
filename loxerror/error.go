package loxerror

import (
	"fmt"

	"github.com/jordanwebster/golox/token"
)

var hadError = false
var hadRuntimeError = false

type RuntimeError struct {
    message string
	token token.Token
}

func (e *RuntimeError) Error() string {
    return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}

func NewRuntimeError(t token.Token, message string) *RuntimeError {
	return &RuntimeError{
        message: message,
		token: t,
	}
}

func ReportRuntimeError(e *RuntimeError) {
    fmt.Println(e.Error())
    hadRuntimeError = true
}

func HadRuntimeError() bool {
    return hadRuntimeError
}

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
