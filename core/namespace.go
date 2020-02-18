package core

import (
	"github.com/keithnull/mal-go/types"
)

// NameSpace is the initial namespace for mal
var NameSpace = map[string]types.MalFunction{
	// arithmetic operators
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
	// string functions
	"pr-str":      strReadable,
	"str":         strUnreadable,
	"prn":         printReadable,
	"println":     printUnreadable,
	"read-string": readString,
	"slurp":       slurp,
	// list related operations
	"list":   createList,
	"list?":  isList,
	"empty?": isEmptyList,
	"count":  getListSize,
	// comparision
	"=":  isEqual,
	"<":  isLess,
	"<=": isLessEqual,
	">":  isGreater,
	">=": isGreaterEqual,
}

// InitCommands contain mal commands to be executed in sequence during initialization
var InitCommands = []string{
	`(def! not (fn* (a) (if a false true)))`,
	`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`,
}
