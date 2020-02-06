package types

type MalType interface{}

type MalNumber struct {
	Value int
}

type MalList []MalType

type MalSymbol struct {
	Value string
}

type MalEnv interface {
	Set(key MalSymbol, value MalType) error
	Find(key MalSymbol) MalEnv
	Get(key MalSymbol) (MalType, error)
}
