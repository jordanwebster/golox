package loxerror

import (
	"fmt"

	"github.com/jordanwebster/golox/token"
)

var hadError = false
var hadRuntimeError = false

type RuntimeError struct {
	message string
	token   token.Token
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}

func NewRuntimeError(t token.Token, message string) *RuntimeError {
	return &RuntimeError{
		message: message,
		token:   t,
	}
}

func ReportRuntimeError(e *RuntimeError) {
	fmt.Println(e.Error())
	hadRuntimeError = true
}

func HadRuntimeError() bool {
	return hadRuntimeError
}

type ParseError struct {
	message string
	token   token.Token
}

func (e *ParseError) Error() string {
	var where string
	if e.token.Type == token.EOF {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%s'", e.token.Lexeme)
	}

	return fmt.Sprintf("[line %d] Error %s where: %s", e.token.Line, where, e.message)
}

func NewParseError(t token.Token, message string) *ParseError {
	return &ParseError{
		message: message,
		token:   t,
	}
}

type SyntaxError struct {
    message string
    line int
}

func (e *SyntaxError) Error() string {
    return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

func NewSyntaxError(line int, message string) *SyntaxError {
    return &SyntaxError{
        message: message,
        line: line,
    }
}

func ReportError(e error) {
	fmt.Println(e.Error())
	hadError = true
}

func ClearError() {
	hadError = false
}

func HadError() bool {
	return hadError
}
