package interpreter

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"github.com/albertlockett/crafting-interpreters-go/lox/stmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
	"strings"
)

type Interpreter struct {
	Env *Environment
}

func NewInterpreter() *Interpreter {
	var clock Callable
	clock = &Clock{}

	globals := newEnvironment(nil)
	globals.define("clock", clock)
	return &Interpreter{ Env: globals, }
}

type RuntimeError struct {
	Line    int
	message string
}

func (e *RuntimeError) Error() string {
	return e.message
}

func (i *Interpreter) Interpret(statements []stmt.Statement) interface{} {
	for j := range statements {
		i.execute(statements[j])
	}
	return nil
}

func (i *Interpreter) execute(s stmt.Statement) {
	s.Accept(i)
}

func (i *Interpreter) executeBlock(b *stmt.Block, env *Environment) {
	// set scope to next scope while evaluating, then unset
	previous := env
	defer func() {
		i.Env = previous
	}()
	i.Env = env

	for j := range b.Statements {
		i.execute(b.Statements[j])
	}
}

func (i *Interpreter) evaluate(e expr.Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	if ok {
		return b
	}
	return true
}

func (i *Interpreter) stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}

	if _, ok := v.(float64); ok {
		text := fmt.Sprintf("%v", v)
		if strings.HasSuffix(text, ".0") {
			return text[0 : len(text)-2]
		}
		return text
	}

	return fmt.Sprintf("%v", v)
}

// stmt.Visitor interface

func (i *Interpreter) VisitVar(v *stmt.Var) interface{} {
	var value interface{} = nil
	if v.Initializer != nil {
		value = i.evaluate(v.Initializer)
	}
	i.Env.define(v.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlock(b *stmt.Block) interface{} {
	i.executeBlock(b, newEnvironment(i.Env))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(s *stmt.ExpressionStmt) interface{} {
	i.evaluate(s.Expression)
	return nil
}

func (i *Interpreter) VisitIfStmt(s *stmt.IfStmt) interface{} {
	condition := i.evaluate(s.Condition)
	if i.isTruthy(condition) {
		i.execute(s.ThenBranch)
	} else if s.ElseBranch != nil {
		i.execute(s.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrint(s *stmt.Print) interface{} {
	val := i.evaluate(s.Expression)
	fmt.Printf("%s\n", i.stringify(val))
	return nil
}

func (i *Interpreter) VisitWhile(s *stmt.While) interface{} {
	for {
		condition := i.evaluate(s.Condition)
		if !i.isTruthy(condition) {
			break
		}

		i.execute(s.Body)
	}
	return nil
}

// expr.Visitor interface:

func (i *Interpreter) VisitAssignment(e *expr.Assignment) interface{} {
	value := i.evaluate(e.Value)
	i.Env.assign(e.Name, value)
	return value
}

func (i *Interpreter) VisitBinary(e *expr.Binary) interface{} {
	left := i.evaluate(e.Left)
	right := i.evaluate(e.Right)

	switch e.Operator.Tokentype {
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)

	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)

	case token.GREATER:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl > dr

	case token.GREATER_EQUAL:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl >= dr

	case token.LESS:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl < dr

	case token.LESS_EQUAL:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl <= dr

	case token.MINUS:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl - dr

	case token.PLUS:
		sl, okl := left.(string)
		sr, okr := right.(string)
		if okl && okr {
			return fmt.Sprintf("%s%s", sl, sr)
		}
		dl, okl := left.(float64)
		dr, okr := right.(float64)
		if okl && okr {
			return dl + dr
		}
		panic(&RuntimeError{
			message: "Operands must be two numbers or two strings",
			Line:    e.Operator.Line,
		})

	case token.SLASH:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl / dr

	case token.STAR:
		dl, _ := left.(float64)
		dr, _ := right.(float64)
		// TODO handle thing
		return dl * dr
	}
	return nil // TODO
}

func (i *Interpreter) VisitCallExpr(e *expr.Call) interface{} {
	callee := i.evaluate(e.Callee)

	args := make([]interface{}, 0)
	for _, arg := range e.Arguments {
		args = append(args, i.evaluate(arg))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic(&RuntimeError{
			message: "Can only call functions and classes.",
			Line:    e.Paren.Line,
		})
	}

	if len(args) != function.Arity() {
		panic(&RuntimeError{
			message: fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(args)),
			Line:    e.Paren.Line,
		})
	}

	return function.Call(i, args)
}

func (i *Interpreter) VisitGrouping(e *expr.Grouping) interface{} {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitLiteral(e *expr.Literal) interface{} {
	return e.Value
}

func (i *Interpreter) VisitLogical(e *expr.Logical) interface{} {
	left := i.evaluate(e.Left)

	if e.Operator.Tokentype == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(e.Right)
}

func (i *Interpreter) VisitUnary(e *expr.Unary) interface{} {
	right := i.evaluate(e.Right)

	switch e.Operator.Tokentype {
	case token.MINUS:
		d, _ := right.(float64)
		// TODO handle not OK
		return -d
	case token.BANG:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitVarExpr(e *expr.Variable) interface{} {
	return i.Env.get(e.Name)
}
