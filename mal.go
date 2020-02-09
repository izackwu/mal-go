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

// EVAL evaluates `ast` within `env` environment
// If any error occurs, the result will be `nil`
func EVAL(ast MalType, env MalEnv) (MalType, error) {
	switch t := ast.(type) {
	case MalList:
		if len(t) == 0 {
			return t, nil
		}
		// the apply phase in Lisp
		first := ""
		if symbol, ok := t[0].(MalSymbol); ok {
			first = symbol.Value
		}
		evaluatedList, err := evalAST(t, env)
		switch first {
		case "def!":
			if len(t) != 3 {
				return nil, fmt.Errorf("incorrect number of parameters for 'def!'")
			}
			k, ok := t[1].(MalSymbol)
			if !ok {
				return nil, fmt.Errorf("the first parameter is expected to be a symbol")
			}
			v, err := EVAL(t[2], env)
			if err != nil {
				return nil, err
			}
			err = env.Set(k, v)
			return v, err
		case "let*":
			if len(t) != 3 {
				return nil, fmt.Errorf("incorrect number of arguments for 'let*'")
			}
			bindings, ok := t[1].(MalList)
			if !ok || len(bindings)%2 != 0 {
				return nil, fmt.Errorf("the first parameter is expected to be a list of even length")
			}
			// do variable bindings in the temp environment
			tempEnv, _ := environment.CreateEnv(env, nil, nil)
			for i := 0; i < len(bindings); i += 2 {
				k, ok := bindings[i].(MalSymbol)
				if !ok {
					return nil, fmt.Errorf("invalid symbol(s) in variable bindings")
				}
				v, err := EVAL(bindings[i+1], tempEnv)
				if err != nil {
					return nil, err
				}
				err = tempEnv.Set(k, v)
				if err != nil {
					return nil, err
				}
			}
			return EVAL(t[2], tempEnv)
		default: // function calling or invalid cases
			if err != nil {
				return nil, err
			}
			f, ok := evaluatedList.(MalList)[0].(MalFunction)
			if !ok {
				return nil, fmt.Errorf("invalid function calling")
			}
			return f(evaluatedList.(MalList)[1:]...)
		}
	case MalVector:
		evaluatedList, err := evalAST(MalList(t), env)
		if err != nil {
			return nil, err
		}
		return MalVector(evaluatedList.(MalList)), nil
	case MalHashmap:
		result := make(MalHashmap)
		for k, v := range t {
			v, err := EVAL(v, env)
			if err != nil {
				return nil, err
			}
			result[k] = v
		}
		return result, nil
	default:
		return evalAST(ast, env)
	}
}

func PRINT(exp MalType) string {
	return printer.PrintStr(exp, true)
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
