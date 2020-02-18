package main

import (
	"fmt"
	"github.com/keithnull/mal-go/core"
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
		return ast, nil
	}
}

// EVAL evaluates `ast` within `env` environment
// If any error occurs, the result will be `nil`
func EVAL(ast MalType, env MalEnv) (MalType, error) {
	// infinite loop for tail call optimization (TCO)
	for {
		// only MalList is handled here, other types will be passed to evalAST() directly
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
				ast, env = t[2], tempEnv // tail call
			case "do":
				for _, exp := range t[1 : len(t)-1] {
					_, err := EVAL(exp, env)
					if err != nil {
						return nil, err
					}
				}
				ast = t[len(t)-1] // tail call
			case "if":
				if len(t) == 3 { // by default, the False branch is nil
					t = append(t, MalNil)
				} else if len(t) != 4 {
					return nil, fmt.Errorf("incorrect number of arguments for 'if'")
				}
				condition, err := EVAL(t[1], env)
				if err != nil {
					return nil, err
				}
				// tail call
				ast = t[2]                                        // True
				if condition == MalFalse || condition == MalNil { // False
					ast = t[3]
				}
			case "fn*":
				if len(t) != 3 {
					return nil, fmt.Errorf("incorrect number of arguments for 'fn*'")
				}
				// make sure parameters are a list of MalSymbols
				// maybe this is unnecessary as it will be checked during future function calling
				// but this makes me feel securer
				params, ok := t[1].(MalList)
				if !ok {
					return nil, fmt.Errorf("the first argument should be function parameter list")
				}
				for i, v := range params {
					if _, ok := v.(MalSymbol); !ok {
						return nil, fmt.Errorf("parameter %d is not a valid symbol", i)
					}
				}
				// it's so good that Golang supports closure, love it~
				closure := func(args ...MalType) (MalType, error) {
					wrappedEnv, err := environment.CreateEnv(env, params, args)
					if err != nil {
						return nil, err
					}
					return EVAL(t[2], wrappedEnv)
				}
				return MalFunctionTCO{
					AST:      t[2],
					Params:   params,
					Env:      env,
					Function: closure,
				}, nil
			default: // function calling or invalid cases
				evaluatedList, err := evalAST(t, env)
				if err != nil {
					return nil, err
				}
				switch f := evaluatedList.(MalList)[0].(type) {
				case MalFunction: // functions defined in core
					return f(evaluatedList.(MalList)[1:]...)
				case MalFunctionTCO: // functions defined with fn*
					ast = f.AST
					env, err = environment.CreateEnv(f.Env, f.Params, evaluatedList.(MalList)[1:])
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid function calling")
				}
			}
		default:
			return evalAST(ast, env)
		}
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

func runInitCommands(env MalEnv) {
	for _, command := range core.InitCommands {
		_ = rep(command, env) // output is ignored
	}
}

func main() {
	defer readline.Close()
	replEnv := environment.GetInitEnv()
	// it breaks the program's structure to add 'eval' function here
	// but doing so is the simplest way
	_ = replEnv.Set(MalSymbol{Value: "eval"}, MalFunction(func(args ...MalType) (MalType, error) {
		return EVAL(args[0], replEnv)
	}))
	runInitCommands(replEnv)
	for { // infinite REPL loop
		input, err := readline.PromptAndRead("user> ")
		if err != nil { // EOF or something unexpected
			break
		}
		fmt.Println(rep(input, replEnv))
	}
}
