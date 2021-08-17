package parser

import (
	"errors"
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func (p *Parser) match(t ...token.TokenType) bool {
	for _, ttype := range t {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(ttype token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Tokentype == ttype
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Tokentype == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	e := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e
}

func (p *Parser) comparison() expr.Expr {
	e := p.term()

	for p.match(token.GREATER, token.EQUAL, token.LESS_EQUAL, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e
}

func (p *Parser) term() expr.Expr {
	e := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e
}

func (p *Parser) factor() expr.Expr {
	e := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e
}

// unary 	-> ( "!" | "-" ) unary
//				|  primary
func (p *Parser) unary() expr.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &expr.Unary{Token: operator, Right: right}
	}
	return p.primary()
}

// primary  -> NUMBER | STRING | "true" | "false" | "nil"
//					|  "(" expression ")"
func (p *Parser) primary() expr.Expr {
	if p.match(token.FALSE) {
		return &expr.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return &expr.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return &expr.Literal{Value: nil}
	}
	if p.match(token.NUMBER, token.STRING) {
		return &expr.Literal{Value: p.previous().Literal}
	}
	if p.match(token.LEFT_PAREN) {
		e := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		return &expr.Grouping{Expression: e}
	}
	panic(errors.New("FUC!!!!!"))
}

func (p *Parser) consume(ttype token.TokenType, s string) {

}
