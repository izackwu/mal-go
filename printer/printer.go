package printer

import (
	"github.com/keithnull/mal-go/types"
	"strconv"
)

func PrintStr(ast types.MalType) string {
	switch t := ast.(type) {
	case types.MalNumber:
		return strconv.Itoa(t.Value)
	case types.MalSymbol:
		return t.Value
	case types.MalList:
		result := "("
		for i, element := range t {
			if i != 0 {
				result += " "
			}
			result += PrintStr(element)
		}
		result += ")"
		return result
	default:
		return "[UNKNOWN VALUE]"
	}
}
