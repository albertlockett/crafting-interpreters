package main

type Scanner struct {
	source  string
	tokens  []*Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: make([]*Token, 0),
	}
}

func (s *Scanner) scanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch string(c) {
	case "(": s.addToken(LEFT_PAREN, nil)
	case ")": s.addToken(RIGHT_PAREN, nil)
	case "[": s.addToken(LEFT_BRACE, nil)
	case "]": s.addToken(RIGHT_BRACE, nil)
	case ",": s.addToken(COMMA, nil)
	case ".": s.addToken(DOT, nil)
	case "-": s.addToken(MINUS, nil)
	case "+": s.addToken(PLUS, nil)
	case ";": s.addToken(SEMICOLON, nil)
	case "*": s.addToken(STAR, nil)
	default:
		lerror(s.line, "Unexpected character")
	}
}

func (s *Scanner) advance() uint8 {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(tokenType TokenType, linteral interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, linteral, s.line))
}