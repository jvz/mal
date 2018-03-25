package main

import (
	"bufio"
	"core"
	. "env"
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
	env := NewEnv()
	for sym, fn := range core.NS {
		env.Set(sym, fn)
	}
	return env
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

func EVAL(ast MalType, env EnvType) (MalType, error) {
	//fmt.Println(ast)
	switch {
	case IsList(ast):
		list := ast.(MalList).Value
		if len(list) == 0 {
			return ast, nil
		}
		sym := "__<*fn*>__"
		if s, ok := list[0].(MalSymbol); ok {
			sym = s.Value
		}
		var a1 MalType
		var a2 MalType
		switch len(list) {
		case 1:
			a1 = nil
			a2 = nil
		case 2:
			a1 = list[1]
			a2 = nil
		default:
			a1 = list[1]
			a2 = list[2]
		}
		switch sym {
		case "def!":
			// define a symbol in the given env
			if len(list) != 3 {
				return nil, fmt.Errorf("def! invalid args: %v", list)
			}
			key, err := GetSymbol(a1)
			if err != nil {
				return nil, err
			}
			val, err := EVAL(a2, env)
			if err != nil {
				return nil, err
			}
			env.Set(key.Value, val)
			return val, nil

		case "let*":
			// create an inner env with ordered bindings and apply it to an expression
			if len(list) != 3 {
				return nil, fmt.Errorf("let* invalid args: %v", list)
			}
			binds, err := GetSlice(a1)
			if err != nil {
				return nil, err
			}
			if len(binds)&1 == 1 {
				return nil, errors.New("odd number of binds provided to let*")
			}
			inner, err := env.New(nil, nil)
			if err != nil {
				return nil, err
			}
			for i := 0; i < len(binds); i += 2 {
				sym, err := GetSymbol(binds[i])
				if err != nil {
					return nil, err
				}
				expr, err := EVAL(binds[i+1], inner)
				if err != nil {
					return nil, err
				}
				inner.Set(sym.Value, expr)
			}
			return EVAL(a2, inner)

		case "do":
			// evaluate all arguments and return the last one's result
			eval, err := evalAst(NewList(list[1:]), env)
			if err != nil {
				return nil, err
			}
			evals, err := GetSlice(eval)
			if err != nil {
				return nil, err
			}
			if len(evals) == 0 {
				return MalNil{}, nil
			}
			return evals[len(evals)-1], nil

		case "if":
			// check first arg, if not nil or false, evaluates and returns second arg
			// otherwise, the third arg is evaluated and returned if provided or nil otherwise
			if len(list) < 3 || len(list) > 4 {
				return nil, fmt.Errorf("if invalid args: %v", list)
			}
			expr, err := EVAL(a1, env)
			if err != nil {
				return nil, err
			}
			if IsTruthy(expr) {
				return EVAL(a2, env)
			}
			if len(list) < 4 {
				return MalNil{}, nil
			}
			return EVAL(list[3], env)

		case "fn*":
			// create a new function closure
			if len(list) != 3 {
				return nil, fmt.Errorf("fn* invalid args: %v", list)
			}
			binds, err := GetSlice(a1)
			if err != nil {
				return nil, err
			}
			return func(args []MalType) (MalType, error) {
				inner, err := env.New(binds, args)
				if err != nil {
					return nil, err
				}
				return EVAL(a2, inner)
			}, nil

		default:
			// evaluate functions
			eval, err := evalAst(ast, env)
			if err != nil {
				return nil, err
			}
			evals, err := GetSlice(eval)
			if err != nil {
				return nil, err
			}
			fn, err := GetFn(evals[0])
			if err != nil {
				return nil, err
			}
			return fn(evals[1:])
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
	rep(`(def! not (fn* (a) (if a false true)))`)
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
