package interpreter

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func newEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(token *token.Token) interface{} {
	value, ok := e.values[token.Lexeme]
	if ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(token)
	}

	panic(&RuntimeError{
		Line:    token.Line,
		message: fmt.Sprintf("Undefined variable '%d'.", token.Line),
	})
}
