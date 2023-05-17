package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO: Consider using parser package instead of regex
func main() {
	source, err := os.ReadFile("ast/expr.go")
	if err != nil {
		panic("Failed to open ast/expr.go")
	}
	re := regexp.MustCompile(`type ([a-zA-Z]+) struct\s?{`)
	var types []string
	for _, match := range re.FindAllStringSubmatch(string(source), -1) {
		types = append(types, match[1])
	}

	template := `package ast

type ExprVisitor interface {
{{visitor_methods}}
}

{{type_accept_methods}}
`

	var visitorMethods []string
	for _, t := range types {
		visitorMethods = append(visitorMethods, fmt.Sprintf("    Visit%s(expr *%s) interface{}", t, t))
	}
	re = regexp.MustCompile(`{{visitor_methods}}`)
	template = re.ReplaceAllString(template, strings.Join(visitorMethods, "\n"))

	var typeAcceptMethods []string
	for _, t := range types {
		f := `func (expr *%s) Accept(visitor ExprVisitor) interface{} {
    return visitor.Visit%s(expr)
}`
		typeAcceptMethods = append(typeAcceptMethods, fmt.Sprintf(f, t, t))
	}
	re = regexp.MustCompile(`{{type_accept_methods}}`)
	template = re.ReplaceAllString(template, strings.Join(typeAcceptMethods, "\n"))

	os.WriteFile("ast/expr_gen.go", []byte(template), 0644)
}
