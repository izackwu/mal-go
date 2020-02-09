package types

type MalType interface{}

type MalNumber struct {
	Value int
}

type MalString struct {
	Value string
}

type MalKeyword struct {
	Value string
}

type MalLiteral string

const (
	MalNil   MalLiteral = "nil"
	MalTrue  MalLiteral = "true"
	MalFalse MalLiteral = "false"
)

type MalList []MalType

type MalVector []MalType

type MalHashmap map[MalType]MalType

type MalSymbol struct {
	Value string
}

type MalFunction func(args ...MalType) (MalType, error)

type MalEnv interface {
	Set(key MalSymbol, value MalType) error
	Find(key MalSymbol) MalEnv
	Get(key MalSymbol) (MalType, error)
}
