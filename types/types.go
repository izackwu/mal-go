package types

type MalType interface{}

type MalNumber struct {
	Value int
}

type MalList []MalType

type MalSymbol struct {
	Value string
}
