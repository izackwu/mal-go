package reader

import (
	"fmt"
	"github.com/keithnull/mal-go/types"
	"regexp"
	"strconv"
)

var (
	tokenRegexp = `[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`
)

type Reader interface {
	Next() (string, error) // returns the token at the current position and increments the position
	Peek() (string, error) // returns the token at the current position
}

type TokenReader struct {
	tokens   []string
	position int
}

func (tr *TokenReader) anyError() error {
	if tr == nil {
		return fmt.Errorf("Nil TokenReader")
	}
	if tr.position >= len(tr.tokens) {
		return fmt.Errorf("Position out of bound")
	}
	return nil
}

func (tr *TokenReader) Peek() (string, error) {
	if err := tr.anyError(); err != nil {
		return "", err
	}
	return tr.tokens[tr.position], nil
}

func (tr *TokenReader) Next() (string, error) {
	if err := tr.anyError(); err != nil {
		return "", err
	}
	token := tr.tokens[tr.position]
	tr.position += 1
	return token, nil
}

func ReadStr(input string) (types.MalType, error) {
	// call tokenize()
	tokens, err := tokenize(input)
	if err != nil{
		return nil, err
	}
	// create a new Reader instance
	tr := TokenReader{tokens, 0}
	// call read_form() with the Reader instance
	return readForm(&tr)
}

// take a single string and return a slice of all the tokens (strings) in it
func tokenize(input string) ([]string, error) {
	re, err := regexp.Compile(tokenRegexp)
	if err != nil {
		return nil, err
	}
	tokens := make([]string, 0)
	for _, group := range re.FindAllStringSubmatch(input, -1) {
		token := group[1]
		if token == "" || token[0] == ';' { // ignore empty tokens and comments
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func readForm(rd Reader) (types.MalType, error) {
	token, err := rd.Peek()
	if err != nil {
		return nil, err
	}
	if token == "(" {
		return readList(rd)
	} else {
		return readAtom(rd)
	}
}

func readAtom(rd Reader) (types.MalType, error) {
	token, err := rd.Next()
	if err != nil {
		return nil, err
	}
	// there should be no error in MatchString()
	if matched, _ := regexp.MatchString(`^[-+]?\d+$`, token); matched { // number
		number, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		return types.MalNumber{number}, nil
	} else { // so far, only symbols are left
		return types.MalSymbol{token}, nil
	}
}

func readList(rd Reader) (types.MalList, error) {
	// sanity check as last peek we already saw "("
	start, _ := rd.Next()
	if start != "(" {
		return nil, fmt.Errorf("Unexpect starting token for a list: %s", start)
	}
	astList := types.MalList{}
	for token, err := rd.Peek(); token != ")"; token, err = rd.Peek() {
		if err != nil {
			return nil, err
		}
		ast, err := readForm(rd)
		if err != nil {
			return nil, err
		}
		astList = append(astList, ast)
	}
	// last peek, we saw ")"
	rd.Next()
	return astList, nil
}
