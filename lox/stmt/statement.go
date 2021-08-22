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
	VisitPrint(*Print) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface{}
	VisitBlock(*Block) interface{}
}

type Var struct {
	Name        *token.Token
	Initializer expr.Expr
}

func (s *Var) Accept(v Visitor) interface{} {
	return v.VisitVar(s)
}

// Print
type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(v Visitor) interface{} {
	return v.VisitPrint(p)
}

// Expression
type ExpressionStmt struct {
	Expression expr.Expr
}

func (e *ExpressionStmt) Accept(v Visitor) interface{} {
	return v.VisitExpressionStmt(e)
}

// Block
type Block struct {
	Statements []Statement
}

func (b *Block) Accept(v Visitor) interface{} {
	return v.VisitBlock(b)
}