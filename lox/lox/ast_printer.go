package lox

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(e expr.Expr) interface{} {
	return e.Accept(a)
}

func (a *AstPrinter) parenthesize(name string, exprs ... expr.Expr) string {
	strs := make([]string, 0)
	strs = append(strs, fmt.Sprintf("(%s", name))
	for _, e := range exprs {
		strs = append(strs, fmt.Sprintf(" %v", e.Accept(a)))
	}
	strs = append(strs, ")")
	return strings.Join(strs, "")
}

func (a *AstPrinter) VisitAssign(e *expr.Assign) interface{} {
	return a.parenthesize(e.Name, e.Value)
}

func (a *AstPrinter) VisitBinary(e *expr.Binary) interface{} {
	return a.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}
func (a *AstPrinter) VisitGrouping(e *expr.Grouping) interface{} {
	return a.parenthesize("group", e.Expression)
}
func (a *AstPrinter) VisitLiteral(e *expr.Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}
func (a *AstPrinter) VisitUnary(e *expr.Unary) interface{} {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}