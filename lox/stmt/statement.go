package stmt

import "github.com/albertlockett/crafting-interpreters-go/lox/expr"

type Statement interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitPrint(*Print) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface {}
}

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(v Visitor) interface{} {
	return v.VisitPrint(p)
}


type ExpressionStmt struct {
	Expression expr.Expr
}

func (e *ExpressionStmt) Accept(v Visitor) interface{} {
	return v.VisitExpressionStmt(e)
}