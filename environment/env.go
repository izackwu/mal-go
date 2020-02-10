package environment

import (
	"fmt"
	"github.com/keithnull/mal-go/core"
	"github.com/keithnull/mal-go/types"
)

// Env implements types.MalEnv interface
type Env struct {
	outer types.MalEnv
	data  map[string]types.MalType
}

// Set takes a symbol `key` and a mal `value`, then set the pair in environment data
func (e *Env) Set(key types.MalSymbol, value types.MalType) error {
	if e == nil {
		return fmt.Errorf("set value in nil environment")
	}
	e.data[key.Value] = value
	return nil
}

// Find takes a symbol `key` and returns the closest environment where `key` is
// If `key` doesn't exist in any outer environment, nil will be returned
func (e *Env) Find(key types.MalSymbol) types.MalEnv {
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

// Get takes a symbol `key`, and returns its value in the closest environment where it exists
// If `key` is not found, an error will be returned
func (e *Env) Get(key types.MalSymbol) (types.MalType, error) {
	targetEnv := e.Find(key)
	if targetEnv == nil {
		return nil, fmt.Errorf("failed to look up '%s' in environments", key.Value)
	}
	return targetEnv.(*Env).data[key.Value], nil
}

// CreateEnv creates a new environment, with `outer` as its outer environment, `binds` and `exps`
// for variable bindings, i.e., `binds[i]` will be bound to `exps[i]`
// Note that `binds` and `exps` should be two types.MalList of equal length
func CreateEnv(outer types.MalEnv, binds types.MalList, exps types.MalList) (*Env, error) {
	env := &Env{
		outer: outer,
		data:  make(map[string]types.MalType),
	}
	if len(binds) != len(exps) {
		return nil, fmt.Errorf("different numbers of bindings and expressions")
	}
	for i := 0; i < len(binds); i++ {
		k, v := binds[i], exps[i]
		if _, ok := k.(types.MalSymbol); !ok {
			return nil, fmt.Errorf("invalid symbol(s) in variable bindings")
		}
		err := env.Set(k.(types.MalSymbol), v)
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

// GetInitEnv creates an initial environment (with only builtin variable bindings)
func GetInitEnv() (e *Env) {
	e, _ = CreateEnv(nil, nil, nil)
	for k, v := range core.NameSpace {
		err := e.Set(types.MalSymbol{Value: k}, v)
		if err != nil {
			return nil
		}
	}
	return
}
