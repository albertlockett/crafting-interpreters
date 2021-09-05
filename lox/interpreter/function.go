package interpreter

type Callable interface {
	Arity() int
	Call(i *Interpreter, arguments []interface{}) interface{}
}


