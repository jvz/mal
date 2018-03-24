package reader

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	. "types"
)

func ReadStr(str string) (MalType, error) {
	tr := NewReader(str)
	return tr.ReadForm()
}

type Reader interface {
	ReadForm() (MalType, error)
}

var tokenPattern = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

func tokenizer(str string) []string {
	tokens := make([]string, 0, 1)
	for _, groups := range tokenPattern.FindAllStringSubmatch(str, -1) {
		if groups[1] == "" || groups[1][0] == ';' {
			// ignore comments and blank lines
			continue
		}
		tokens = append(tokens, groups[1])
	}
	return tokens
}

func NewReader(text string) *TokenReader {
	tokens := tokenizer(text)
	reader := TokenReader{tokens: tokens}
	return &reader
}

type TokenReader struct {
	pos    uint
	tokens []string
}

func (tr *TokenReader) next() *string {
	if len(tr.tokens) == 0 {
		return nil
	}
	next := &tr.tokens[0]
	tr.tokens = tr.tokens[1:]
	return next
}

func (tr *TokenReader) peek() *string {
	if len(tr.tokens) == 0 {
		return nil
	}
	return &tr.tokens[0]
}

func (tr *TokenReader) ReadForm() (MalType, error) {
	tok := tr.peek()
	if tok == nil {
		return nil, errors.New("ReadForm underflow")
	}
	switch *tok {
	case "(":
		list, err := tr.readList("(", ")")
		if err != nil {
			return nil, err
		}
		return MalList{Value: list, StartStr: "(", EndStr: ")"}, nil
	case ")":
		return nil, errors.New("unexpected )")
	case "[":
		vec, err := tr.readList("[", "]")
		if err != nil {
			return nil, err
		}
		return MalList{Value: vec, StartStr: "[", EndStr: "]"}, nil
	case "]":
		return nil, errors.New("unexpected ]")
	case "{":
		keyValues, err := tr.readList("{", "}")
		if err != nil {
			return nil, err
		}
		if len(keyValues)&1 != 0 {
			return nil, errors.New("expected an even number of params to a map literal")
		}
		m := make(map[MalType]MalType)
		for i := 0; i < len(keyValues); i += 2 {
			m[keyValues[i]] = keyValues[i+1]
		}
		return MalMap{Value: m}, nil
	case "}":
		return nil, errors.New("unexpected }")
	default:
		return tr.readAtom()
	}
}

var stringEscapesReplacer = strings.NewReplacer(`\"`, `"`, `\n`, "\n", `\\`, `\`)
var intPattern = regexp.MustCompile(`^-?[0-9]+$`)

func (tr *TokenReader) readAtom() (MalType, error) {
	tok := tr.next()
	if tok == nil {
		return nil, errors.New("readAtom underflow")
	}
	switch {
	case (*tok)[0] == '"':
		end := strings.LastIndex(*tok, `"`)
		if end <= 0 {
			return nil, errors.New("unbalanced quotes")
		}
		contents := stringEscapesReplacer.Replace((*tok)[1:end])
		return MalString{Value: contents}, nil
	case (*tok)[0] == ':':
		keyword := (*tok)[1:]
		return MalKeyword{Value: keyword}, nil
	case intPattern.MatchString(*tok):
		i, err := strconv.Atoi(*tok)
		if err != nil {
			return nil, err
		} else {
			return MalInt{Value: i}, nil
		}
	}
	switch *tok {
	case "'":
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "quote"}, form), nil
	case "`":
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "quasiquote"}, form), nil
	case "~":
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "unquote"}, form), nil
	case "~@":
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "splice-unquote"}, form), nil
	case "^":
		meta, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "with-meta"}, form, meta), nil
	case "@":
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		return NewListOf(MalSymbol{Value: "deref"}, form), nil
	case "nil":
		return MalNil{}, nil
	case "true":
		return MalTrue, nil
	case "false":
		return MalFalse, nil
	default:
		return MalSymbol{Value: *tok}, nil
	}
}

func (tr *TokenReader) readList(start, end string) ([]MalType, error) {
	tok := tr.next() // (
	if tok == nil {
		return nil, errors.New("readList underflow")
	}
	if *tok != start {
		return nil, fmt.Errorf("expected %s", start)
	}
	list := make([]MalType, 0, 1)
	for {
		tok = tr.peek()
		if tok == nil {
			return nil, fmt.Errorf("expected %s", end)
		}
		if *tok == end {
			break
		}
		form, err := tr.ReadForm()
		if err != nil {
			return nil, err
		}
		list = append(list, form)
	}
	tr.next() // )
	return list, nil
}
