package environment

import (
	"fmt"
	. "github.com/keithnull/mal-go/types"
)

var initialBindings = map[string]MalType{
	"+": func(args ...MalType) (MalType, error) {
		if err := AssertLength(args, 2); err != nil {
			return nil, err
		}
		a, ok1 := args[0].(MalNumber)
		b, ok2 := args[1].(MalNumber)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid operands")
		}
		return MalNumber{Value: a.Value + b.Value}, nil
	},
	"-": func(args ...MalType) (MalType, error) {
		if err := AssertLength(args, 2); err != nil {
			return nil, err
		}
		a, ok1 := args[0].(MalNumber)
		b, ok2 := args[1].(MalNumber)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid operands")
		}
		return MalNumber{Value: a.Value - b.Value}, nil
	},
	"*": func(args ...MalType) (MalType, error) {
		if err := AssertLength(args, 2); err != nil {
			return nil, err
		}
		a, ok1 := args[0].(MalNumber)
		b, ok2 := args[1].(MalNumber)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid operands")
		}
		return MalNumber{Value: a.Value * b.Value}, nil
	},
	"/": func(args ...MalType) (MalType, error) {
		if err := AssertLength(args, 2); err != nil {
			return nil, err
		}
		a, ok1 := args[0].(MalNumber)
		b, ok2 := args[1].(MalNumber)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid operands")
		}
		if b.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return MalNumber{Value: a.Value / b.Value}, nil
	},
}

// *Env implements MalEnv interface
type Env struct {
	outer MalEnv
	data  map[string]MalType
}

func (e *Env) Set(key MalSymbol, value MalType) error {
	if e == nil {
		return fmt.Errorf("set value in nil environment")
	}
	e.data[key.Value] = value
	return nil
}

func (e *Env) Find(key MalSymbol) MalEnv {
	if e == nil {
		return nil
	}
	if _, ok := e.data[key.Value]; ok {
		return e
	} else if e.outer != nil {
		return e.outer.Find(key)
	} else {
		return nil
	}
}

func (e *Env) Get(key MalSymbol) (MalType, error) {
	targetEnv := e.Find(key)
	if targetEnv == nil {
		return nil, fmt.Errorf("failed to look up '%s' in environments", key.Value)
	}
	return targetEnv.(*Env).data[key.Value], nil
}

func CreateEnv(outer MalEnv) *Env {
	return &Env{
		outer: outer,
		data:  make(map[string]MalType),
	}
}

func GetInitEnv() (e *Env) {
	e = CreateEnv(nil)
	for k, v := range initialBindings {
		err := e.Set(MalSymbol{Value: k}, v)
		if err != nil {
			return nil
		}
	}
	return
}
