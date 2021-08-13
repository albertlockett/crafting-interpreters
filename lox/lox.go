package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			panic(err) // TODO exit more smarterly
		}
	} else {
		if err := runPrompt(); err != nil {
			panic(err)
		}
	}
}

func runFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return run(string(data))
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
	}
}

func run(source string) error {
	fmt.Printf("%s\n", source)
	return nil
}