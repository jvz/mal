package core

import (
	"fmt"
	"io/ioutil"
	"printer"
	"reader"
	"strings"
	. "types"
)

func MonoFunc(f func(MalType) MalType) func([]MalType) (MalType, error) {
	return MonoErrFunc(func(a MalType) (MalType, error) {
		return f(a), nil
	})
}

func MonoErrFunc(f func(MalType) (MalType, error)) func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("invalid args: %v", args)
		}
		return f(args[0])
	}
}

func BiFunc(f func(MalType, MalType) MalType) func([]MalType) (MalType, error) {
	return BiErrFunc(func(a MalType, b MalType) (MalType, error) {
		return f(a, b), nil
	})
}

func BiErrFunc(f func(MalType, MalType) (MalType, error)) func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("invalid args: %v", args)
		}
		return f(args[0], args[1])
	}
}

func intBiFunc(f func(int, int) int) func([]MalType) (MalType, error) {
	return BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
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
	return BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
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
	`list?`: MonoFunc(func(a MalType) MalType {
		return MalBool{Value: IsList(a)}
	}),
	`empty?`: MonoErrFunc(func(a MalType) (MalType, error) {
		list, err := GetSlice(a)
		if err != nil {
			return nil, err
		}
		return MalBool{Value: len(list) == 0}, nil
	}),
	`count`: MonoErrFunc(func(a MalType) (MalType, error) {
		switch arg := a.(type) {
		case MalNil:
			return MalInt{Value: 0}, nil
		case MalList:
			return MalInt{Value: len(arg.Value)}, nil
		default:
			return RaiseTypeError("list", arg)
		}
	}),
	`=`: BiFunc(func(a MalType, b MalType) MalType {
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
	`read-string`: MonoErrFunc(func(a MalType) (MalType, error) {
		str, err := GetString(a)
		if err != nil {
			return nil, err
		}
		return reader.ReadStr(str.Value)
	}),
	`slurp`: MonoErrFunc(func(a MalType) (MalType, error) {
		str, err := GetString(a)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadFile(str.Value)
		if err != nil {
			return nil, err
		}
		return MalString{Value: string(content)}, nil
	}),
	`atom`: MonoFunc(func(a MalType) MalType {
		return NewAtom(a)
	}),
	`atom?`: MonoFunc(func(a MalType) MalType {
		_, ok := a.(*MalAtom)
		return MalBool{Value: ok}
	}),
	`deref`: MonoErrFunc(func(a MalType) (MalType, error) {
		atom, err := GetAtom(a)
		if err != nil {
			return nil, err
		}
		return atom.Value, nil
	}),
	`reset!`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		atom, err := GetAtom(a1)
		if err != nil {
			return nil, err
		}
		atom.Set(a2)
		return a2, nil
	}),
	`swap!`: func(args []MalType) (MalType, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid args: %v", args)
		}
		atom, err := GetAtom(args[0])
		if err != nil {
			return nil, err
		}
		fn, err := GetFn(args[1])
		if err != nil {
			return nil, err
		}
		fnArgs := make([]MalType, len(args)-1)
		fnArgs[0] = atom.Value
		if len(args) > 2 {
			copy(fnArgs[1:], args[2:])
		}
		res, err := fn(fnArgs)
		if err != nil {
			return nil, err
		}
		atom.Set(res)
		return res, nil
	},
	`cons`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		tail, err := GetSlice(a2)
		if err != nil {
			return nil, err
		}
		list := make([]MalType, len(tail)+1)
		list[0] = a1
		copy(list[1:], tail)
		return NewList(list), nil
	}),
	`concat`: func(args []MalType) (MalType, error) {
		concat := make([]MalType, 0, 1)
		for _, arg := range args {
			list, err := GetSlice(arg)
			if err != nil {
				return nil, err
			}
			concat = append(concat, list...)
		}
		return NewList(concat), nil
	},
	`nth`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		list, err := GetSlice(a1)
		if err != nil {
			return nil, err
		}
		index, err := GetInt(a2)
		if err != nil {
			return nil, err
		}
		i := index.Value
		if i < 0 || i >= len(list) {
			return nil, fmt.Errorf("index out of ranges: %v", i)
		}
		return list[i], nil
	}),
	`first`: MonoErrFunc(func(a MalType) (MalType, error) {
		switch list := a.(type) {
		case MalNil:
			return MalNil{}, nil
		case MalList:
			if len(list.Value) == 0 {
				return MalNil{}, nil
			}
			return list.Value[0], nil
		default:
			return RaiseTypeError("list", a)
		}
	}),
	`rest`: MonoErrFunc(func(a MalType) (MalType, error) {
		switch list := a.(type) {
		case MalNil:
			return NewListOf(), nil
		case MalList:
			if len(list.Value) <= 1 {
				return NewListOf(), nil
			}
			return NewList(list.Value[1:]), nil
		default:
			return RaiseTypeError("list", a)
		}
	}),
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

	case *MalAtom:
		if b, ok := b.(*MalAtom); ok {
			return equal(a.Value, b.Value)
		} else {
			return false
		}

	case MalSymbol:
		if b, ok := b.(MalSymbol); ok {
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

	case MalKeyword:
		if b, ok := b.(MalKeyword); ok {
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
