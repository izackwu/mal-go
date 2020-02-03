package main

import (
	"fmt"
	"github.com/keithnull/mal-go/printer"
	"github.com/keithnull/mal-go/reader"
	"github.com/keithnull/mal-go/readline"
	"github.com/keithnull/mal-go/types"
)

func READ(in string) types.MalType {
	ast, err := reader.ReadStr(in)
	if err != nil {
		return nil
	}
	return ast
}

func EVAL(ast types.MalType, env string) types.MalType {
	return ast
}

func PRINT(exp types.MalType) string {
	return printer.PrintStr(exp)
}

func rep(in string) string {
	return PRINT(EVAL(READ(in), ""))
}

func main() {
	defer readline.Close()
	for { // infinite REPL loop
		input, err := readline.PromptAndRead("user> ")
		if err != nil { // EOF or something unexpected
			break
		}
		fmt.Println(rep(input))
	}
}
