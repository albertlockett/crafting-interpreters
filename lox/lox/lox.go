package lox

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/interpreter"
	"github.com/albertlockett/crafting-interpreters-go/lox/parser"
	"github.com/albertlockett/crafting-interpreters-go/lox/scanner"
	"github.com/albertlockett/crafting-interpreters-go/lox/token"
	"io/ioutil"
	"os"
	"strings"
)

var HasError = false
var HasRuntimeError = false

var Interpreter = interpreter.NewInterpreter()

const CODE_INVALID_USAGE = 64
const CODE_ERROR = 65
const CODE_RUNTIME_ERROR = 70

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
	fmt.Printf("~~~ SOURCE ~~~\n")
	lines := strings.Split(source, "\n")
	for i := range lines {
		fmt.Printf("%d  %s\n", i+1, lines[i])
	}
	fmt.Printf("~~~~~~~~~~~~~~\n\n")


	// scan the tokens
	s := scanner.NewScanner(source, Lerror)
	tokens := s.ScanTokens()

	fmt.Printf("~~ TOKENS ~~\n")
	for _, t := range tokens {
		fmt.Printf("%s\n", t)
	}
	fmt.Printf("~~~~~~~~~~~~\n\n")

	// parse the tokens
	p := parser.NewParser(tokens, Terror)
	stmts, _ := p.Parse()
	if HasError {
		return nil
	}

	// handle panics
	defer func() {
		if r := recover(); r != nil {
			// handle runtime error
			if runtimeError, ok := r.(*interpreter.RuntimeError); ok {
				fmt.Printf(
					"RuntimeError[line %d]: %s\n",
					runtimeError.Line,
					runtimeError.Error(),
				)
				HasRuntimeError = true
			} else {
				panic(r)
			}
		}
	}()

	fmt.Printf("~~~ OUTPUT ~~~\n")
	val := Interpreter.Interpret(stmts)
	fmt.Printf("~~~~~~~~~~~~~~\n")
	fmt.Printf("\n~~~ INTERPRETER ~~~\n retval: %v\n~~~~~~~~~~~~~~~~~~~\n", val)

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