package lox

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"github.com/albertlockett/crafting-interpreters-go/lox/stmt"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) PrintStmts(stmts []stmt.Statement) interface{} {
	vals := make([]string, 0)
	for i := range stmts {
		val := a.PrintStmt(stmts[i])
		text, ok := val.(string)
		if !ok {
			panic(fmt.Sprintf("Expected %v / %s to be string", text, text))
		}
		vals = append(vals, text)
	}
	return strings.Join(vals, "\n")
}

func (a *AstPrinter) PrintStmt(s stmt.Statement) interface{} {
	return s.Accept(a)
}

func (a *AstPrinter) Print(e expr.Expr) interface{} {
	return e.Accept(a)
}

func (a *AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	strs := make([]string, 0)
	strs = append(strs, fmt.Sprintf("(%s", name))
	for _, e := range exprs {
		strs = append(strs, fmt.Sprintf(" %v", e.Accept(a)))
	}
	strs = append(strs, ")")
	return strings.Join(strs, "")
}

func (a *AstPrinter) withSemicolon(text interface{}) string {
	return fmt.Sprintf("%s;", text)
}

// stmt.Visitor interface

func (a *AstPrinter) VisitVar(v *stmt.Var) interface{} {
	return a.withSemicolon(a.parenthesize(fmt.Sprintf("var %s =", v.Name.Lexeme), v.Initializer))
}

func (a *AstPrinter) VisitExpressionStmt(s *stmt.ExpressionStmt) interface {} {
	return a.withSemicolon(s.Expression.Accept(a));
}

func (a *AstPrinter ) VisitPrint(s *stmt.Print) interface{} {
	return a.withSemicolon(a.parenthesize("print", s.Expression))
}

func (a *AstPrinter ) VisitBlock(b *stmt.Block) interface{} {
	vals := make([]string, 0)
	vals = append(vals, "{")

	for i := range b.Statements {
		val := a.PrintStmt(b.Statements[i])
		vals = append(vals, val.(string))
	}
	vals = append(vals, "}")
	return strings.Join(vals, "")
}

// expr.Visitor interface

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
	if str, ok := e.Value.(string); ok {
		return fmt.Sprintf(`"%s"`, str)
	}
	return fmt.Sprintf("%v", e.Value)
}

func (a *AstPrinter) VisitUnary(e *expr.Unary) interface{} {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}

func (a *AstPrinter) VisitVarExpr(e *expr.Variable) interface{} {
	return "var[" + e.Name.Lexeme + "]"
}