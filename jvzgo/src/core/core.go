package core

import (
	"fmt"
	"printer"
	"strings"
	. "types"
)

var NS = map[string]MalType{
	`+`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`+`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalInt{Value: a.Value + b.Value}, nil
	},
	`-`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`-`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalInt{Value: a.Value - b.Value}, nil
	},
	`*`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`*`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalInt{Value: a.Value * b.Value}, nil
	},
	`/`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`/`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalInt{Value: a.Value / b.Value}, nil
	},
	`list`: func(args []MalType) (MalType, error) {
		return MalList{Value: args, StartStr: "(", EndStr: ")"}, nil
	},
	`list?`: func(args []MalType) (MalType, error) {
		if len(args) != 1 {
			return raiseInvalid(`list?`, args)
		}
		return MalBool{Value: IsList(args[0])}, nil
	},
	`empty?`: func(args []MalType) (MalType, error) {
		if len(args) != 1 {
			return raiseInvalid(`empty?`, args)
		}
		list, err := GetSlice(args[0])
		if err != nil {
			return nil, err
		}
		return MalBool{Value: len(list) == 0}, nil
	},
	`count`: func(args []MalType) (MalType, error) {
		if len(args) != 1 {
			return raiseInvalid(`count`, args)
		}
		switch arg := args[0].(type) {
		case MalList:
			return MalInt{Value: len(arg.Value)}, nil
		case MalNil:
			return MalInt{Value: 0}, nil
		default:
			return RaiseTypeError("list", arg)
		}
	},
	`=`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`=`, args)
		}
		if equal(args[0], args[1]) {
			return MalTrue, nil
		} else {
			return MalFalse, nil
		}
	},
	`<`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`<`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalBool{Value: a.Value < b.Value}, nil
	},
	`<=`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`<=`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalBool{Value: a.Value <= b.Value}, nil
	},
	`>`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`>`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalBool{Value: a.Value > b.Value}, nil
	},
	`>=`: func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return raiseInvalid(`>=`, args)
		}
		a, ok := args[0].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[0])
		}
		b, ok := args[1].(MalInt)
		if !ok {
			return RaiseTypeError("int", args[1])
		}
		return MalBool{Value: a.Value >= b.Value}, nil
	},
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

func raiseInvalid(sym string, args []MalType) (MalType, error) {
	return nil, fmt.Errorf("%v invalid args: %v", sym, args)
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
