package main

import (
	"bufio"
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

type envMap = map[string]func([]MalType) (MalType, error)

var replEnv = envMap{
	`+`: func(args []MalType) (MalType, error) {
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
	},
	`-`: func(args []MalType) (MalType, error) {
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
	},
	`*`: func(args []MalType) (MalType, error) {
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
	},
	`/`: func(args []MalType) (MalType, error) {
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
	},
}

func evalAst(ast MalType, env envMap) (MalType, error) {
	switch ast := ast.(type) {
	case MalSymbol:
		if expr, ok := env[ast.Value]; ok {
			return expr, nil
		} else {
			return nil, fmt.Errorf("unknown env key: %v", ast.Value)
		}
	case MalList:
		evals := make([]MalType, len(ast.Value))
		for i, arg := range ast.Value {
			res, err := EVAL(arg, env)
			if err != nil {
				return nil, err
			}
			evals[i] = res
		}
		return ast.New(evals), nil
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

func EVAL(ast MalType, env envMap) (MalType, error) {
	switch {
	case IsList(ast):
		if len(ast.(MalList).Value) == 0 {
			return ast, nil
		}
		eval, err := evalAst(ast, env)
		if err != nil {
			return nil, err
		}
		list := eval.(MalList).Value
		fn, ok := list[0].(func([]MalType) (MalType, error))
		if !ok {
			return RaiseTypeError("function", list[0])
		}
		return fn(list[1:])
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
