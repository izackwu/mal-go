package core

import (
	"github.com/keithnull/mal-go/types"
)

// NameSpace is the initial namespace for mal
var NameSpace = map[string]types.MalFunction{
	"+":      add,
	"-":      sub,
	"*":      mul,
	"/":      div,
	"prn":    printValue,
	"list":   createList,
	"list?":  isList,
	"empty?": isEmptyList,
	"count":  getListSize,
	"=":      isEqual,
	"<":      isLess,
	"<=":     isLessEqual,
	">":      isGreater,
	">=":     isGreaterEqual,
}
