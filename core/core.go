package core

import (
	"fmt"
	"github.com/keithnull/mal-go/printer"
	"github.com/keithnull/mal-go/types"
)

// AssertLength asserts the length of a list
func AssertLength(args []types.MalType, expect int) error {
	if actual := len(args); actual != expect {
		return fmt.Errorf("incorrect number of arguments: expect %d but get %d", expect, actual)
	}
	return nil
}

// assertTwoNumbers asserts that `args` are a list of two numbers
func assertTwoNumbers(args []types.MalType) (int, int, error) {
	if err := AssertLength(args, 2); err != nil {
		return 0, 0, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("invalid operand(s)")
	}
	return a.Value, b.Value, nil
}

func add(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{Value: a + b}, nil
}

func sub(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{Value: a - b}, nil
}

func mul(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{Value: a * b}, nil
}

func div(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	if b == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return types.MalNumber{Value: a / b}, nil
}

func printValue(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 1); err != nil {
		return nil, err
	}
	fmt.Println(printer.PrintStr(args[0], true))
	return types.MalNil, nil
}

func createList(args ...types.MalType) (types.MalType, error) {
	// it's necessary to force cast to MalList (though nothing is done actually)
	return types.MalList(args), nil
}

func isList(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 1); err != nil {
		return nil, err
	}
	_, ok := args[0].(types.MalList)
	return types.ToMalBool(ok), nil
}

func isEmptyList(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 1); err != nil {
		return nil, err
	}
	lst, ok := args[0].(types.MalList)
	if !ok {
		return nil, fmt.Errorf("can't check whether a non-list is empty")
	}
	return types.ToMalBool(len(lst) == 0), nil
}

func getListSize(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 1); err != nil {
		return nil, err
	}
	// MalNil is a special case
	literal, ok := args[0].(types.MalLiteral)
	if ok && literal == types.MalNil {
		return types.MalNumber{Value: 0}, nil
	}
	// MalList
	lst, ok := args[0].(types.MalList)
	if !ok {
		return nil, fmt.Errorf("can't count the number of elements in a non-list")
	}
	return types.MalNumber{Value: len(lst)}, nil
}

func isEqual(args ...types.MalType) (types.MalType, error) {
	if err := AssertLength(args, 2); err != nil {
		return nil, err
	}
	same := true
	switch first := args[0].(type) {
	case types.MalList:
		second, ok := args[1].(types.MalList)
		// necessary condition for equality: second should be a list of the same length as first
		if !ok || len(first) != len(second) {
			same = false
			break
		}
		for i := 0; i < len(first); i++ {
			if b, _ := isEqual(first[i], second[i]); b == types.MalFalse {
				same = false
				break
			}
		}
	case types.MalNumber:
		second, ok := args[1].(types.MalNumber)
		same = ok && first.Value == second.Value
	case types.MalLiteral:
		second, ok := args[1].(types.MalLiteral)
		same = ok && first == second
	case types.MalString:
		second, ok := args[1].(types.MalString)
		same = ok && first.Value == second.Value
	case types.MalVector:
		second, ok := args[1].(types.MalVector)
		if ok { // convert both to MalList and then compare
			return isEqual(types.MalList(first), types.MalList(second))
		}
		same = false
	default:
		return nil, fmt.Errorf("sorry but unimplemented yet")
	}
	return types.ToMalBool(same), nil
}

// for <,<=,>,>=, actually only one of them needs implementing as others can be derived from it
// and =. But considering the complexity of =, I implement both < and >.

func isLess(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.ToMalBool(a < b), nil

}

func isGreater(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.ToMalBool(a > b), nil
}

func isLessEqual(args ...types.MalType) (types.MalType, error) {
	greater, err := isGreater(args...) // note the ... here
	if err != nil {
		return nil, err
	}
	return types.NotMalBool(greater.(types.MalLiteral)), nil
}
func isGreaterEqual(args ...types.MalType) (types.MalType, error) {
	less, err := isLess(args...)
	if err != nil {
		return nil, err
	}
	return types.NotMalBool(less.(types.MalLiteral)), nil
}
