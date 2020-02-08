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

// ReadStr builds a Mal AST with the given string
func ReadStr(input string) (types.MalType, error) {
	// call tokenize()
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	// create a new Reader instance
	tr := TokenReader{tokens, 0}
	// call readForm() with the Reader instance
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
	switch token {
	case "(":
		return readList(rd)
	case ")":
		return nil, fmt.Errorf("unexpected ']'")
	case "[":
		return readVector(rd)
	case "]":
		return nil, fmt.Errorf("unexpected ']")
	case "{":
		return readHashmap(rd)
	case "}":
		return nil, fmt.Errorf("unexpected '}")
	default:
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
		return types.MalNumber{Value: number}, nil
	} else if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"?$`, token); matched { // string
		if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"$`, token); !matched {
			return nil, fmt.Errorf("unclosed string: %s", token)
		}
		unquoted, err := strconv.Unquote(token) // unquote and handle escape chars gracefully
		if err != nil {
			return nil, err
		}
		return types.MalString{Value: unquoted}, nil
	} else if token == "nil" {
		return types.MalNil, nil
	} else if token == "true" {
		return types.MalTrue, nil
	} else if token == "false" {
		return types.MalFalse, nil
	} else if token[0] == ':' {
		return types.MalKeyword{Value: token[1:]}, nil
	} else { // so far, only symbols are left
		return types.MalSymbol{Value: token}, nil
	}
}

// readStartEnd assumes the next token is `start` and reads till `end`
// It returns a list of Mal objects
// If any error encountered, it will stop reading immediately and return that error
func readStartEnd(rd Reader, start, end string) (types.MalList, error) {
	// sanity check as last peek we already saw the starting token
	first, _ := rd.Next()
	if first != start {
		return nil, fmt.Errorf("incorrect starting token: expect '%s' but get '%s'", start, first)
	}
	astList := types.MalList{}
	for token, err := rd.Peek(); token != end; token, err = rd.Peek() {
		if err != nil {
			return nil, err
		}
		ast, err := readForm(rd)
		if err != nil {
			return nil, err
		}
		astList = append(astList, ast)
	}
	// last peek, we saw the ending token
	_, _ = rd.Next()
	return astList, nil
}

func readList(rd Reader) (types.MalList, error) {
	return readStartEnd(rd, "(", ")")
}

func readVector(rd Reader) (types.MalType, error) {
	list, err := readStartEnd(rd, "[", "]")
	if err != nil {
		return nil, err
	}
	return types.MalVector(list), nil
}

func readHashmap(rd Reader) (types.MalType, error) {
	list, err := readStartEnd(rd, "{", "}")
	if err != nil {
		return nil, err
	}
	if len(list)%2 != 0 {
		return nil, fmt.Errorf("incorrect number of elements for a hashmap")
	}
	hashmap := make(types.MalHashmap)
	for i := 0; i < len(list); i += 2 {
		switch t := list[i].(type) {
		case types.MalKeyword, types.MalString:
			hashmap[t] = list[i+1]
		default:
			return nil, fmt.Errorf("hashmap keys only accept string or keyword")
		}
	}
	return hashmap, nil
}
