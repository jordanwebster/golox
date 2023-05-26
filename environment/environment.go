package environment

import (
	"fmt"

	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/token"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
    return &Environment{
        values: make(map[string]interface{}),
    }
}

func (environment *Environment) Define(name string, value interface{}) {
	environment.values[name] = value
}

func (environment *Environment) Get(name token.Token) (interface{}, error) {
	if value, isPresent := environment.values[name.Lexeme]; isPresent {
		return value, nil
	}

	return nil, loxerror.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
