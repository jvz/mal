package core

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"printer"
	"reader"
	"strings"
	"time"
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

func MonoPred(f func(MalType) bool) func([]MalType) (MalType, error) {
	return MonoFunc(func(a MalType) MalType {
		return MalBool{Value: f(a)}
	})
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
		return NewList(args), nil
	},
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
		return atom.Value(), nil
	}),
	`reset!`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		atom, err := GetAtom(a1)
		if err != nil {
			return nil, err
		}
		atom.SetValue(a2)
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
		fnArgs[0] = atom.Value()
		if len(args) > 2 {
			copy(fnArgs[1:], args[2:])
		}
		res, err := fn(fnArgs)
		if err != nil {
			return nil, err
		}
		atom.SetValue(res)
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
	`throw`: MonoErrFunc(func(a MalType) (MalType, error) {
		return nil, MalError{Value: a}
	}),
	`apply`: func(args []MalType) (MalType, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("apply invalid args: %v", args)
		}
		fn, err := GetFn(args[0])
		if err != nil {
			return nil, err
		}
		last, err := GetSlice(args[len(args)-1])
		if err != nil {
			return nil, err
		}
		fnArgs := make([]MalType, len(last)+len(args)-2)
		copy(fnArgs[:len(args)-2], args[1:len(args)-1])
		copy(fnArgs[len(args)-2:], last)
		return fn(fnArgs)
	},
	`map`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		fn, err := GetFn(a1)
		if err != nil {
			return nil, err
		}
		list, err := GetSlice(a2)
		if err != nil {
			return nil, err
		}
		ret := make([]MalType, len(list))
		for i, v := range list {
			res, err := fn([]MalType{v})
			if err != nil {
				return nil, err
			}
			ret[i] = res
		}
		return NewList(ret), nil
	}),
	`nil?`:     MonoPred(IsNil),
	`true?`:    MonoPred(IsTrue),
	`false?`:   MonoPred(IsFalse),
	`symbol?`:  MonoPred(IsSymbol),
	`keyword?`: MonoPred(IsKeyword),
	`string?`:  MonoPred(IsString),
	`number?`:  MonoPred(IsInt),
	`fn?`:      MonoPred(IsFn),
	`macro?`:   MonoPred(IsMacro),
	`list?`:    MonoPred(IsList),
	`vector?`:  MonoPred(IsVec),
	`map?`:     MonoPred(IsMap),
	`symbol`: MonoErrFunc(func(a MalType) (MalType, error) {
		str, err := GetString(a)
		if err != nil {
			return nil, err
		}
		return MalSymbol{Value: str.Value}, nil
	}),
	`keyword`: MonoErrFunc(func(a MalType) (MalType, error) {
		str, err := GetString(a)
		if err != nil {
			return nil, err
		}
		return MalKeyword{Value: str.Value}, nil
	}),
	`vector`: func(args []MalType) (MalType, error) {
		return NewVec(args), nil
	},
	`hash-map`: func(args []MalType) (MalType, error) {
		if len(args)&1 == 1 {
			return nil, fmt.Errorf("hash-map invalid number of args: %v", args)
		}
		m := make(map[MalType]MalType)
		for i := 0; i < len(args); i += 2 {
			m[args[i]] = args[i+1]
		}
		return MalMap{Value: m}, nil
	},
	`assoc`: func(args []MalType) (MalType, error) {
		if len(args)&1 != 1 {
			return nil, fmt.Errorf("hash-map invalid number of args: %v", args)
		}
		m, err := GetMap(args[0])
		if err != nil {
			return nil, err
		}
		updated := CopyMap(m)
		for i := 1; i < len(args); i += 2 {
			updated.Value[args[i]] = args[i+1]
		}
		return updated, nil
	},
	`dissoc`: func(args []MalType) (MalType, error) {
		if len(args) == 0 {
			return nil, errors.New("dissoc invalid args")
		}
		m, err := GetMap(args[0])
		if err != nil {
			return nil, err
		}
		updated := CopyMap(m)
		for _, key := range args[1:] {
			delete(updated.Value, key)
		}
		return updated, nil
	},
	`get`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		if IsNil(a1) {
			return MalNil{}, nil
		}
		m, err := GetMap(a1)
		if err != nil {
			return nil, err
		}
		if val, ok := m.Value[a2]; ok {
			return val, nil
		}
		return MalNil{}, nil
	}),
	`contains?`: BiErrFunc(func(a1 MalType, a2 MalType) (MalType, error) {
		m, err := GetMap(a1)
		if err != nil {
			return nil, err
		}
		_, ok := m.Value[a2]
		return MalBool{Value: ok}, nil
	}),
	`keys`: MonoErrFunc(func(a MalType) (MalType, error) {
		m, err := GetMap(a)
		if err != nil {
			return nil, err
		}
		keys := make([]MalType, 0, len(m.Value))
		for key := range m.Value {
			keys = append(keys, key)
		}
		return NewList(keys), nil
	}),
	`vals`: MonoErrFunc(func(a MalType) (MalType, error) {
		m, err := GetMap(a)
		if err != nil {
			return nil, err
		}
		vals := make([]MalType, 0, len(m.Value))
		for _, val := range m.Value {
			vals = append(vals, val)
		}
		return NewList(vals), nil
	}),
	`sequential?`: MonoPred(func(a MalType) bool {
		_, ok := a.(MalList)
		return ok
	}),
	`readline`: MonoErrFunc(func(a MalType) (MalType, error) {
		prompt, err := GetString(a)
		if err != nil {
			return nil, err
		}
		fmt.Print(prompt.Value)
		in := bufio.NewScanner(os.Stdin)
		if in.Scan() {
			line := strings.TrimSpace(in.Text())
			return MalString{Value: line}, nil
		}
		if in.Err() != nil {
			return nil, in.Err()
		}
		return MalNil{}, nil
	}),
	`meta`:      MonoFunc(GetMeta),
	`with-meta`: BiErrFunc(WithMeta),
	`time-ms`: func(args []MalType) (MalType, error) {
		nanos := time.Duration(time.Now().UnixNano())
		millis := nanos.Truncate(time.Millisecond).Nanoseconds() / int64(time.Millisecond)
		return MalInt{Value: int(millis)}, nil
	},
	`conj`: func(args []MalType) (MalType, error) {
		switch {
		case len(args) < 2:
			return nil, fmt.Errorf("conj invalid args: %v", args)
		case IsList(args[0]):
			list, _ := GetSlice(args[0])
			conj := make([]MalType, 0, len(list)+len(args)-1)
			for i := len(args) - 1; i > 0; i-- {
				conj = append(conj, args[i])
			}
			conj = append(conj, list...)
			return NewList(conj), nil
		case IsVec(args[0]):
			vec, _ := GetSlice(args[0])
			conj := make([]MalType, 0, len(vec)+len(args)-1)
			conj = append(conj, vec...)
			conj = append(conj, args[1:]...)
			return NewVec(conj), nil
		case IsNil(args[0]):
			return NewList(args[1:]), nil
		default:
			return RaiseTypeError("collection", args[0])
		}
	},
	`seq`: MonoErrFunc(func(a MalType) (MalType, error) {
		switch {
		case IsNil(a):
			return MalNil{}, nil
		case IsList(a):
			list, _ := GetSlice(a)
			if len(list) == 0 {
				return MalNil{}, nil
			}
			return a, nil
		case IsVec(a):
			list, _ := GetSlice(a)
			if len(list) == 0 {
				return MalNil{}, nil
			}
			return NewList(list), nil
		case IsString(a):
			str := a.(MalString).Value
			if len(str) == 0 {
				return MalNil{}, nil
			}
			r := strings.NewReader(str)
			chars := make([]MalType, 0, len(str))
			for ch, _, err := r.ReadRune(); err == nil; ch, _, err = r.ReadRune() {
				char := fmt.Sprintf("%c", ch)
				chars = append(chars, MalString{Value: char})
			}
			return NewList(chars), nil
		default:
			return RaiseTypeError("sequence", a)
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
			return equal(a.Value(), b.Value())
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
