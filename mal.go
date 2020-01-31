package main

import (
	"fmt"
	"github.com/keithnull/mal-go/readline"
)

func READ(in string) string {
	return in
}

func EVAL(ast string, env string) string {
	return ast
}

func PRINT(exp string) string {
	return exp
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
