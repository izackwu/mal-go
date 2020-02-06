package main

import (
	"fmt"
	"github.com/keithnull/mal-go/environment"
	"github.com/keithnull/mal-go/printer"
	"github.com/keithnull/mal-go/reader"
	"github.com/keithnull/mal-go/readline"
	. "github.com/keithnull/mal-go/types" // not recommended but convenient
)

func READ(in string) (MalType, error) {
	ast, err := reader.ReadStr(in)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func evalAST(ast MalType, env MalEnv) (MalType, error) {
	switch t := ast.(type) {
	case MalSymbol:
		elem, err := env.Get(t)
		if err != nil {
			return nil, err
		}
		return elem, nil
	case MalList:
		evaluatedList := make(MalList, 0)
		for _, origin := range t {
			if evaluated, err := EVAL(origin, env); err == nil {
				evaluatedList = append(evaluatedList, evaluated)
			} else {
				return nil, err
			}
		}
		return evaluatedList, nil
	default:
		return ast, nil
	}
}

func EVAL(ast MalType, env MalEnv) (MalType, error) {
	switch t := ast.(type) {
	case MalList:
		if len(t) == 0 {
			return t, nil
		}
		evaluatedList, err := evalAST(t, env)
		if err != nil {
			return nil, err
		}
		f, ok := evaluatedList.(MalList)[0].(func(...MalType) (MalType, error))
		if !ok {
			return nil, fmt.Errorf("invalid function calling")
		}
		return f(evaluatedList.(MalList)[1:]...)
	default:
		return evalAST(ast, env)
	}
}

func PRINT(exp MalType) string {
	return printer.PrintStr(exp)
}

func rep(in string, env MalEnv) string {
	ast, err := READ(in)
	if err != nil {
		return fmt.Sprint(err)
	}
	exp, err := EVAL(ast, env)
	if err != nil {
		return fmt.Sprint(err)
	}
	output := PRINT(exp)
	return output
}

func main() {
	defer readline.Close()
	replEnv := environment.GetInitEnv()
	for { // infinite REPL loop
		input, err := readline.PromptAndRead("user> ")
		if err != nil { // EOF or something unexpected
			break
		}
		fmt.Println(rep(input, replEnv))
	}
}
