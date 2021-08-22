package parser

import (
	"errors"
	"github.com/albertlockett/crafting-interpreters-go/lox/expr"
	"github.com/albertlockett/crafting-interpreters-go/lox/stmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
)

type Parser struct {
	tokens  []*token.Token
	onError func(t *token.Token, message string)
	current int
}

type ParseError error

func NewParser(tokens []*token.Token, onError func(t *token.Token, message string)) *Parser {
	return &Parser{tokens: tokens, onError: onError}
}

// program -> stmt* EOF
func (p *Parser) Parse() ([]stmt.Statement, error) {
	statements := make([]stmt.Statement, 0)
	for !p.isAtEnd() {
		statement, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

// statement -> exprStmt
//						| printStmt
func (p *Parser) statement() (stmt.Statement, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return nil, nil
}

// printStmt -> "print" expression ";"
func (p *Parser) printStatement() (stmt.Statement, error) {
	expression, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &stmt.Print{Expression: expression}, nil
}

// expression -> equality
func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

// equality -> comparison ( ( "==" | "!=" ) comparison)*
func (p *Parser) equality() (expr.Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (expr.Expr, error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e, nil
}

// term -> factor ( ( "+" | "-" ) factor)*
func (p *Parser) term() (expr.Expr, error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e, nil
}

// factor -> unary ( ( "/" | "*" ) unary) *
func (p *Parser) factor() (expr.Expr, error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{
			Left:     e,
			Operator: operator,
			Right:    right,
		}
	}
	return e, nil
}

// unary 	-> ( "!" | "-" ) unary
//				|  primary
func (p *Parser) unary() (expr.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

// primary  -> NUMBER | STRING | "true" | "false" | "nil"
//					|  "(" expression ")"
func (p *Parser) primary() (expr.Expr, error) {
	if p.match(token.FALSE) {
		return &expr.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return &expr.Literal{Value: true}, nil
	}
	if p.match(token.NIL) {
		return &expr.Literal{Value: nil}, nil
	}
	if p.match(token.NUMBER, token.STRING) {
		return &expr.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(token.LEFT_PAREN) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping{Expression: e}, nil
	}

	return nil, p.error(p.peek(), "Expect Expression")
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

func (p *Parser) error(token *token.Token, message string) ParseError {
	p.onError(token, message)
	return ParseError(errors.New(message))
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(ttype token.TokenType, message string) (*token.Token, error) {
	if p.check(ttype) {
		return p.advance(), nil
	}
	return nil, p.error(p.peek(), message)
}
