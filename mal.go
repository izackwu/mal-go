package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	reader := bufio.NewReader(os.Stdin)
	for { // infinite REPL loop
		fmt.Print("user> ")
		line, err := reader.ReadString('\n')
		if err != nil { // EOF or something unexpected
			break
		}
		line = strings.TrimRight(line, "\n")
		fmt.Println(rep(line))
	}
}
