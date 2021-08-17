package lox

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/parser"
	"github.com/albertlockett/crafting-interpreters-go/lox/scanner"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
	"io/ioutil"
	"os"
)

var HasError = false

const CODE_INVALID_USAGE = 64
const CODE_ERROR = 65



func RunFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := run(string(data)); err != nil {
		return err
	}
	return nil
}

func RunPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		data, err, _ := reader.ReadLine()
		if err {
			return errors.New("Error reading input")
		}
		line := string(data)
		if line == "" {
			return nil
		}
		if err := run(line); err != nil {
			return err
		}
		HasError = false
	}
}

func run(source string) error {

	// scan the tokens
	s := scanner.NewScanner(source, Lerror)
	tokens := s.ScanTokens()
	for _, t := range tokens {
		fmt.Printf("%s\n", t)
	}

	// parse the tokens
	p := parser.NewParser(tokens, Terror)
	ast, _ := p.Parse()
	if HasError {
		return nil
	}

	printer := &AstPrinter{}
	fmt.Printf("%s", printer.Print(ast))

	return nil
}

func Lerror(line int, message string) {
	report(line, "", message)
}

func Terror(t *token.Token, message string) {
	if t.Tokentype == token.EOF {
		report(t.Line, "at end", message);
	} else {
		report(t.Line, fmt.Sprintf("'%s'", t.Lexeme), message)
	}
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HasError = true
}