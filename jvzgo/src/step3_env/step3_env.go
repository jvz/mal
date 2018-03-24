package main

import (
	"bufio"
	"env"
	"errors"
	"fmt"
	"os"
	"printer"
	"reader"
	"strings"
	. "types"
)

func READ(str string) (MalType, error) {
	return reader.ReadStr(str)
}

var replEnv = newReplEnv()

func newReplEnv() EnvType {
	e := env.NewEnv()
	e.Set(`+`, func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid args")
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
	})
	e.Set(`-`, func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid args")
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
	})
	e.Set(`*`, func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid args")
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
	})
	e.Set(`/`, func(args []MalType) (MalType, error) {
		if len(args) != 2 {
			return nil, errors.New("invalid args")
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
	})
	return e
}

func evalAst(ast MalType, env EnvType) (MalType, error) {
	switch ast := ast.(type) {
	case MalSymbol:
		return env.Get(ast.Value)
	case MalList:
		evals := make([]MalType, len(ast.Value))
		for i, arg := range ast.Value {
			res, err := EVAL(arg, env)
			if err != nil {
				return nil, err
			}
			evals[i] = res
		}
		return MalList{Value: evals, StartStr: ast.StartStr, EndStr: ast.EndStr}, nil
	case MalMap:
		evals := make(map[MalType]MalType)
		for k, v := range ast.Value {
			res, err := EVAL(v, env)
			if err != nil {
				return nil, err
			}
			evals[k] = res
		}
		return MalMap{Value: evals}, nil
	default:
		return ast, nil
	}
}

func EVAL(ast MalType, env EnvType) (MalType, error) {
	switch {
	case IsList(ast):
		list := ast.(MalList).Value
		if len(list) == 0 {
			return ast, nil
		}
		sym, ok := list[0].(MalSymbol)
		if !ok {
			return RaiseTypeError("symbol", list[0])
		}
		switch sym.Value {
		case "def!":
			if len(list) != 3 {
				return nil, fmt.Errorf("def! invalid args: %v", list)
			}
			key, ok := list[1].(MalSymbol)
			if !ok {
				return RaiseTypeError("symbol", list[1])
			}
			val, err := EVAL(list[2], env)
			if err != nil {
				return nil, err
			}
			env.Set(key.Value, val)
			return val, nil

		case "let*":
			if len(list) != 3 {
				return nil, fmt.Errorf("let* invalid args: %v", list)
			}
			binds, err := GetSlice(list[1])
			if err != nil {
				return nil, err
			}
			if len(binds)&1 == 1 {
				return nil, errors.New("odd number of binds provided to let*")
			}
			inner, _ := env.New(nil, nil)
			for i := 0; i < len(binds); i += 2 {
				sym, ok := binds[i].(MalSymbol)
				if !ok {
					return RaiseTypeError("symbol", binds[i])
				}
				expr, err := EVAL(binds[i+1], inner)
				if err != nil {
					return nil, err
				}
				inner.Set(sym.Value, expr)
			}
			return EVAL(list[2], inner)

		default:
			eval, err := evalAst(ast, env)
			if err != nil {
				return nil, err
			}
			list = eval.(MalList).Value
			fn, ok := list[0].(func([]MalType) (MalType, error))
			if !ok {
				return RaiseTypeError("function", list[0])
			}
			return fn(list[1:])
		}

	default:
		return evalAst(ast, env)
	}
}

func PRINT(exp MalType) (string, error) {
	return printer.PrintStr(exp, true), nil
}

func rep(str string) (string, error) {
	ast, err := READ(str)
	if err != nil {
		return "", err
	}
	exp, err := EVAL(ast, replEnv)
	if err != nil {
		return "", err
	}
	return PRINT(exp)
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("user> ")
		if in.Scan() {
			read := strings.TrimSpace(in.Text())
			result, err := rep(read)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		} else {
			err := in.Err()
			if err != nil {
				panic(err)
			}
			return
		}
	}
}
