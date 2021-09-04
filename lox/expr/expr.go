package expr

import "github.com/albertlockett/crafting-interpreters-go/lox/token"

type Expr interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitAssignment(*Assignment) interface{}
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitLogical(*Logical) interface{}
	VisitUnary(*Unary) interface{}
	VisitVarExpr(*Variable) interface{}
}

// Assignment
type Assignment struct {
	Value Expr
	Name  *token.Token
}

func (a *Assignment) Accept(v Visitor) interface{} {
	return v.VisitAssignment(a)
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

func (g *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGrouping(g)
}

// Literal
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteral(l)
}

// Logical
type Logical struct {
	Left     Expr
	Right    Expr
	Operator *token.Token
}

func (l *Logical) Accept(v Visitor) interface{} {
	return v.VisitLogical(l)
}

// Unary
type Unary struct {
	Operator *token.Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnary(u)
}

// Variable
type Variable struct {
	Name *token.Token
}

func (e *Variable) Accept(v Visitor) interface{} {
	return v.VisitVarExpr(e)
}
