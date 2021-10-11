package stmt

import (
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
)

type Statement interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitVar(*Var) interface{}
	VisitBlock(*Block) interface{}
	VisitFunction(*Function) interface{}
	VisitIfStmt(*IfStmt) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface{}
	VisitPrint(*Print) interface{}
	VisitWhile(*While) interface{}
}

type Var struct {
	Name        *token.Token
	Initializer expr.Expr
}

func (s *Var) Accept(v Visitor) interface{} {
	return v.VisitVar(s)
}

// Block
type Block struct {
	Statements []Statement
}

func (b *Block) Accept(v Visitor) interface{} {
	return v.VisitBlock(b)
}

// Expression
type ExpressionStmt struct {
	Expression expr.Expr
}

func (e *ExpressionStmt) Accept(v Visitor) interface{} {
	return v.VisitExpressionStmt(e)
}

// Function
type Function struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Statement
}

func (f *Function) Accept(v Visitor) interface{} {
	return v.VisitFunction(f)
}

// If
type IfStmt struct {
	Condition  expr.Expr
	ThenBranch Statement
	ElseBranch Statement
}

func (i *IfStmt) Accept(v Visitor) interface{} {
	return v.VisitIfStmt(i)
}

// Print
type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(v Visitor) interface{} {
	return v.VisitPrint(p)
}

// While
type While struct {
	Condition expr.Expr
	Body      Statement
}

func (w *While) Accept(v Visitor) interface{} {
	return v.VisitWhile(w)
}
