package printer

import (
	"github.com/keithnull/mal-go/types"
	"strconv"
)

func printList(lst types.MalList, start, end string, readable bool) string {
	result := start
	for i, element := range lst {
		if i != 0 {
			result += " "
		}
		result += PrintStr(element, readable)
	}
	result += end
	return result
}

func printHashmap(hm types.MalHashmap, readable bool) string {
	result := "{"
	isFirstPair := true
	for k, v := range hm {
		if !isFirstPair {
			result += " "
		}
		isFirstPair = false
		result += PrintStr(k, readable)
		result += " "
		result += PrintStr(v, readable)
	}
	result += "}"
	return result
}

// PrintStr converts a Mal AST to string
// If readable is set to true, then MalString will get escaped properly
func PrintStr(ast types.MalType, readable bool) string {
	switch t := ast.(type) {
	case types.MalNumber:
		return strconv.Itoa(t.Value)
	case types.MalSymbol:
		return t.Value
	case types.MalString:
		if readable {
			return strconv.Quote(t.Value)
		}
		return t.Value
	case types.MalLiteral: // nil, true, false
		return string(t)
	case types.MalKeyword:
		return ":" + t.Value
	case types.MalList: // (foo bar baz)
		return printList(t, "(", ")", readable)
	case types.MalVector: // [foo bar baz]
		return printList(types.MalList(t), "[", "]", readable)
	case types.MalHashmap: // {foo bar}
		return printHashmap(t, readable)
	case types.MalFunction:
		return "#<function>"
	case types.MalFunctionTCO:
		return "#<functionTCO>"
	default:
		return "/UNKNOWN VALUE/"
	}
}
