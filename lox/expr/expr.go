package expr

import "github.com/albertlockett/crafting-interpreters-go/lox/token"

type Expr interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitAssign(*Assign) interface{}
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitUnary(*Unary) interface{}
}

// Assign
type Assign struct {
	Name  string
	Value Expr
}

func (a *Assign) Accept(v Visitor) interface{} {
	return v.VisitAssign(a)
}

// Binary
type Binary struct {
	Right    Expr
	Operator *token.Token
	Left     Expr
}

func (b *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinary(b)
}

// Grouping
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept (v Visitor) interface{}{
	return v.VisitGrouping(g)
}

// Literal
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteral(l)
}

// Unary
type Unary struct {
	Token *token.Token
	Right Expr
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnary(u)
}
