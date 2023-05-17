package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type codegenArgs struct {
	path           string
	visitorName    string
	visitorArgName string
}

func main() {
	args := []codegenArgs{
		{
			path:           "ast/expr.go",
			visitorName:    "ExprVisitor",
			visitorArgName: "expr",
		},
		{
			path:           "ast/stmt.go",
			visitorName:    "StmtVisitor",
			visitorArgName: "stmt",
		},
	}

	for _, arg := range args {
		codegenFile(arg)
	}
}

func codegenFile(args codegenArgs) {
	source, err := os.ReadFile(args.path)
	if err != nil {
		panic(fmt.Sprintf("Failed to open %s", args.path))
	}
	re := regexp.MustCompile(`type ([a-zA-Z]+) struct\s?{`)
	var types []string
	for _, match := range re.FindAllStringSubmatch(string(source), -1) {
		types = append(types, match[1])
	}

	template := `package ast

type {{interface_name}} interface {
{{visitor_methods}}
}

{{type_accept_methods}}
`

	// TODO: Consider using parser package instead of regex
	re = regexp.MustCompile(`{{interface_name}}`)
	template = re.ReplaceAllString(template, args.visitorName)

	var visitorMethods []string
	for _, t := range types {
		visitorMethods = append(visitorMethods, fmt.Sprintf("    Visit%s(%s *%s) (interface{}, error)", t, args.visitorArgName, t))
	}
	re = regexp.MustCompile(`{{visitor_methods}}`)
	template = re.ReplaceAllString(template, strings.Join(visitorMethods, "\n"))

	var typeAcceptMethods []string
	for _, t := range types {
		f := `func (%s *%s) Accept(visitor %s) (interface{}, error) {
    return visitor.Visit%s(%s)
}`
		typeAcceptMethods = append(typeAcceptMethods, fmt.Sprintf(f, args.visitorArgName, t, args.visitorName, t, args.visitorArgName))
	}
	re = regexp.MustCompile(`{{type_accept_methods}}`)
	template = re.ReplaceAllString(template, strings.Join(typeAcceptMethods, "\n"))

	output := strings.TrimSuffix(args.path, ".go") + "_gen.go"
	os.WriteFile(output, []byte(template), 0644)
}
