package interpreter

import "time"

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
