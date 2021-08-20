package main

import (
	"fmt"
	"github.com/albertlockett/crafting-interpreters-go/lox/lox"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(lox.CODE_INVALID_USAGE)
	} else if len(os.Args) == 2 {
		if err := lox.RunFile(os.Args[1]); err != nil {
			panic(err) // TODO exit more smarterly
		}
	} else {
		if err := lox.RunPrompt(); err != nil {
			panic(err)
		}
	}

	if lox.HasError {
		os.Exit(lox.CODE_ERROR)
	}

	if lox.HasRuntimeError {
		os.Exit(lox.CODE_RUNTIME_ERROR)
	}
}
