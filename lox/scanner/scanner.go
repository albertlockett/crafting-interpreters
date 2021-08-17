package scanner

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
	"strconv"
)

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"fun":    token.FUN,
	"for":    token.FOR,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source  string
	tokens  []*token.Token
	start   int
	current int
	line    int
	onError func(line int, message string)
}

func NewScanner(source string, onError func(line int, message string)) *Scanner {
	return &Scanner{
		line:   1,
		source: source,
		tokens: make([]*token.Token, 0),
		onError: onError,
	}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '[':
		s.addToken(token.LEFT_BRACE, nil)
	case ']':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '=':
		tokenType := token.EQUAL
		if s.matches('=') {
			tokenType = token.EQUAL_EQUAL
		}
		s.addToken(tokenType, nil)
	case '!':
		tokenType := token.BANG
		if s.matches('=') {
			tokenType = token.BANG_EQUAL
		}
		s.addToken(tokenType, nil)
	case '<':
		tokenType := token.LESS
		if s.matches('=') {
			tokenType = token.LESS_EQUAL
		}
		s.addToken(tokenType, nil)
	case '>':
		tokenType := token.GREATER
		if s.matches('=') {
			tokenType = token.GREATER_EQUAL
		}
		s.addToken(tokenType, nil)
	case '/':
		if s.matches('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.onError(s.line, "Unexpected character")
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() uint8 {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(tokenType token.TokenType, linteral interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, linteral, s.line))
}

func (s *Scanner) matches(char uint8) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] == char {
		s.current++
		return true
	}

	return false
}

func (s *Scanner) peek() uint8 {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() uint8 {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.onError(s.line, "Unterminated string.")
		return
	}
	s.advance() // the closing "

	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *Scanner) isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
	}

	for s.isDigit(s.peek()) {
		s.advance()
	}

	v := s.source[s.start:s.current]
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		s.onError(s.line, fmt.Sprint("lexer invalid number %s", v))
	}
	s.addToken(token.NUMBER, f)
}

func (s *Scanner) isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c uint8) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType := keywords[text]
	if tokenType == "" {
		tokenType = token.IDENTIFIER
	}
	s.addToken(tokenType, nil)
}
