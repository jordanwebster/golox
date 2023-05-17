package error

import "fmt"

var hadError = false

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
