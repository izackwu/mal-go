package core

import (
	"fmt"
	"github.com/keithnull/mal-go/types"
)

// AssertLength asserts the length of a list
func AssertLength(args []types.MalType, expect int) error {
	if actual := len(args); actual != expect {
		return fmt.Errorf("incorrect number of arguments: expect %d but get %d", expect, actual)
	}
	return nil
}

func add(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 2); err != nil {
		return nil, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid operand(s) for '+'")
	}
	return types.MalNumber{Value: a.Value + b.Value}, nil
}

func sub(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 2); err != nil {
		return nil, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid operand(s) for '-'")
	}
	return types.MalNumber{Value: a.Value - b.Value}, nil
}

func mul(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 2); err != nil {
		return nil, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid operand(s) for '*'")
	}
	return types.MalNumber{Value: a.Value * b.Value}, nil
}

func div(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 2); err != nil {
		return nil, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid operand(s) for '/'")
	}
	if b.Value == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return types.MalNumber{Value: a.Value / b.Value}, nil
}
