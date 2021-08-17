package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var hasError = false

const CODE_INVALID_USAGE = 64
const CODE_ERROR = 65

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(CODE_INVALID_USAGE)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			panic(err) // TODO exit more smarterly
		}
	} else {
		if err := runPrompt(); err != nil {
			panic(err)
		}
	}

	if hasError {
		os.Exit(CODE_ERROR)
	}
}

func runFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := run(string(data)); err != nil {
		return err
	}
	return nil
}

func runPrompt() error {
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
		hasError = false
	}
}

func run(source string) error {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Printf("%s\n", token)
	}
	return nil
}

func lerror(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	hasError = true
}