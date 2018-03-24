package core

import (
	"fmt"
	"printer"
	"strings"
	. "types"
)

func monoFunc(f func(MalType) MalType) func([]MalType) (MalType, error) {
	return monoErrFunc(func(a MalType) (MalType, error) {
		return f(a), nil
	})
}

func monoErrFunc(f func(MalType) (MalType, error)) func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid args: %v", args)
		}
		return f(args[0])
	}
}

func biFunc(f func(MalType, MalType) MalType) func([]MalType) (MalType, error) {
	return biErrFunc(func(a MalType, b MalType) (MalType, error) {
		return f(a, b), nil
	})
}

func biErrFunc(f func(MalType, MalType) (MalType, error)) func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid args: %v", args)
		}
		return f(args[0], args[1])
	}
}

func intBiFunc(f func(int, int) int) func([]MalType) (MalType, error) {
	return biErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		a, err := GetInt(a1)
		if err != nil {
			return nil, err
		}
		b, err := GetInt(a2)
		if err != nil {
			return nil, err
		}
		return MalInt{Value: f(a.Value, b.Value)}, nil
	})
}

func intBiPred(f func(int, int) bool) func([]MalType) (MalType, error) {
	return biErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		a, err := GetInt(a1)
		if err != nil {
			return nil, err
		}
		b, err := GetInt(a2)
		if err != nil {
			return nil, err
		}
		return MalBool{Value: f(a.Value, b.Value)}, nil
	})
}

var NS = map[string]MalType{
	`+`: intBiFunc(func(a int, b int) int {
		return a + b
	}),
	`-`: intBiFunc(func(a int, b int) int {
		return a - b
	}),
	`*`: intBiFunc(func(a int, b int) int {
		return a * b
	}),
	`/`: intBiFunc(func(a int, b int) int {
		return a / b
	}),
	`list`: func(args []MalType) (MalType, error) {
		return MalList{Value: args, StartStr: "(", EndStr: ")"}, nil
	},
	`list?`: monoFunc(func(a MalType) MalType {
		return MalBool{Value: IsList(a)}
	}),
	`empty?`: monoErrFunc(func(a MalType) (MalType, error) {
		list, err := GetSlice(a)
		if err != nil {
			return nil, err
		}
		return MalBool{Value: len(list) == 0}, nil
	}),
	`count`: monoErrFunc(func(a MalType) (MalType, error) {
		switch arg := a.(type) {
		case MalNil:
			return MalInt{Value: 0}, nil
		case MalList:
			return MalInt{Value: len(arg.Value)}, nil
		default:
			return RaiseTypeError("list", arg)
		}
	}),
	`=`: biFunc(func(a MalType, b MalType) MalType {
		return MalBool{Value: equal(a, b)}
	}),
	`<`: intBiPred(func(a int, b int) bool {
		return a < b
	}),
	`<=`: intBiPred(func(a int, b int) bool {
		return a <= b
	}),
	`>`: intBiPred(func(a int, b int) bool {
		return a > b
	}),
	`>=`: intBiPred(func(a int, b int) bool {
		return a >= b
	}),
	`pr-str`: func(args []MalType) (MalType, error) {
		prints := make([]string, len(args))
		for i, arg := range args {
			prints[i] = printer.PrintStr(arg, true)
		}
		return MalString{Value: strings.Join(prints, " ")}, nil
	},
	`str`: func(args []MalType) (MalType, error) {
		str := strings.Builder{}
		for _, arg := range args {
			_, err := str.WriteString(printer.PrintStr(arg, false))
			if err != nil {
				return nil, err
			}
		}
		return MalString{Value: str.String()}, nil
	},
	`prn`: func(args []MalType) (MalType, error) {
		prints := make([]string, len(args))
		for i, arg := range args {
			prints[i] = printer.PrintStr(arg, true)
		}
		fmt.Println(strings.Join(prints, " "))
		return MalNil{}, nil
	},
	`println`: func(args []MalType) (MalType, error) {
		prints := make([]string, len(args))
		for i, arg := range args {
			prints[i] = printer.PrintStr(arg, false)
		}
		fmt.Println(strings.Join(prints, " "))
		return MalNil{}, nil
	},
}

func equal(a, b MalType) bool {
	switch a := a.(type) {
	case MalList:
		as := a.Value
		bs, err := GetSlice(b)
		if err != nil {
			return false
		}
		if len(as) != len(bs) {
			return false
		}
		for i := range as {
			if !equal(as[i], bs[i]) {
				return false
			}
		}
		return true

	case MalMap:
		if b, ok := b.(MalMap); ok {
			if len(a.Value) != len(b.Value) {
				return false
			}
			for key, x := range a.Value {
				y, ok := b.Value[key]
				if !ok || !equal(x, y) {
					return false
				}
			}
			return true
		} else {
			return false
		}

	case MalSymbol:
		if b, ok := b.(MalString); ok {
			return a.Value == b.Value
		} else {
			return false
		}

	case MalString:
		if b, ok := b.(MalString); ok {
			return a.Value == b.Value
		} else {
			return false
		}

	case MalInt:
		if b, ok := b.(MalInt); ok {
			return a.Value == b.Value
		} else {
			return false
		}

	case MalBool:
		if b, ok := b.(MalBool); ok {
			return a.Value == b.Value
		} else {
			return false
		}

	case MalNil:
		_, ok := b.(MalNil)
		return ok

	default:
		return a == b
	}
}
