package interpreter

import (
	"time"

	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/environment"
)

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
}

type ClockCallable struct{}

func (callable *ClockCallable) Arity() int {
	return 0
}

func (callable *ClockCallable) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return (time.Now().UnixMilli() / 1000.0), nil
}

type LoxFunction struct {
	declaration *ast.FunctionStmt
	closure     *environment.Environment
}

func NewFunction(declaration *ast.FunctionStmt, closure *environment.Environment) LoxCallable {
	return &LoxFunction{
		declaration,
		closure,
	}
}

func (function *LoxFunction) Arity() int {
	return len(function.declaration.Parameters)
}

func (function *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := environment.NewEnvironment(function.closure)
	for i, param := range function.declaration.Parameters {
		env.Define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(function.declaration.Body, env)
	switch v := err.(type) {
	case *Return:
		return v.Value, nil
	default:
		return nil, err
	}
}
