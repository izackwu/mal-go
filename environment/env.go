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
	// first pass: check `binds` type and whether it's variadic
	variadic := false
	for i, k := range binds {
		symbol, ok := k.(types.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("invalid symbol(s) in variable bindings")
		}
		if symbol.Value == "&" { // variadic function parameters
			if i != len(binds)-2 {
				return nil, fmt.Errorf("invalid position for '&' in bindings")
			}
			variadic = true
		}
	}
	// check the length of `binds` and `exps`
	if variadic && len(binds)-2 > len(exps) {
		return nil, fmt.Errorf("not enough expressions for a variadic function")
	} else if !variadic && len(binds) != len(exps) {
		return nil, fmt.Errorf(
			"different numbers of bindings and expressions for a non-variadic function")
	}
	// second pass: do variable bindings actually
	for i, k := range binds {
		// skip "&"
		if variadic && i == len(binds)-2 {
			continue
		}
		symbol, _ := k.(types.MalSymbol)
		var v types.MalType
		if variadic && i == len(binds)-1 {
			v = exps[i-1:]
		} else { // in case of index out of bound
			v = exps[i]
		}
		err := env.Set(symbol, v)
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
