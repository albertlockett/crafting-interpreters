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

// program -> declaration * EOF
func (p *Parser) Parse() ([]stmt.Statement, error) {
	statements := make([]stmt.Statement, 0)
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

// declaration -> varDeclaration
//							| statement
func (p *Parser) declaration() (stmt.Statement, error) {
	var s stmt.Statement
	var err error

	if p.match(token.VAR) {
		s, err = p.varDeclaration()
	} else {
		s, err = p.statement()
	}

	if err != nil {
		if _, ok := err.(ParseError); ok {
			p.synchronize()
		} else {
			return nil, err
		}
	}
	return s, nil
}

// varDeclaration -> "var" IDENTIFIER ( "=" expression)? ";"
func (p *Parser) varDeclaration() (stmt.Statement, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected variable name.")
	if err != nil {
		return nil, err
	}

	var initializer expr.Expr = nil
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(token.SEMICOLON, "Expected ';' after variable declaration.")
	return &stmt.Var{
		Name:        name,
		Initializer: initializer,
	}, nil
}

// statement -> exprStmt
//            | ifStmt
//						| printStmt
//						| block
func (p *Parser) statement() (stmt.Statement, error) {
	if p.match(token.IF) {
		return p.ifStmt()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	return p.expressionStmt()
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStmt() (stmt.Statement, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if cnondition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch stmt.Statement = nil
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &stmt.IfStmt{
		Condition: condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
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

// block -> "{" declaration "}"
func (p *Parser) block() (stmt.Statement, error) {
	statements := make([]stmt.Statement, 0)
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		d, err := p.declaration()
		statements = append(statements, d)
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return &stmt.Block{
		Statements: statements,
	}, nil
}

// expressionStmt -> expression ";"
func (p *Parser) expressionStmt() (stmt.Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Exprected ';' after value.")
	if err != nil {
		return nil, err
	}

	return &stmt.ExpressionStmt{Expression: value}, nil
}

// expression -> assignment
func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment
//						|	equality
func (p *Parser) assignment() (expr.Expr, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := e.(*expr.Variable); ok {
			name := v.Name
			return &expr.Assignment{
				Name: name,
				Value: value,
			}, nil
		}

		return nil, p.error(equals, "Invalid assignment target.")
	}

	return e, nil
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
// 					| IDENTIFIER
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

	if p.match(token.IDENTIFIER) {
		return &expr.Variable{Name: p.previous()}, nil
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

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Tokentype == token.SEMICOLON {
			return
		}

		switch p.peek().Tokentype {
		case token.CLASS:
		case token.FOR:
		case token.FUN:
		case token.IF:
		case token.PRINT:
		case token.RETURN:
		case token.VAR:
		case token.WHILE:
			return
		}
		p.advance()
	}
}
