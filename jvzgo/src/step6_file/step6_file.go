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
	//fmt.Println(ast)
	for {
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
				env = inner
				ast = a2
				continue

			case "do":
				// evaluate all arguments and return the last one's result
				switch len(list) {
				case 1:
					return MalNil{}, nil
				case 2:
					ast = list[1]
					continue
				default:
					if _, err := evalAst(NewList(list[1:len(list)-1]), env); err != nil {
						return nil, err
					}
					ast = list[len(list)-1]
					continue
				}

			case "if":
				// check first arg, if not nil or false, evaluates and returns second arg
				// otherwise, the third arg is evaluated and returned if provided or nil otherwise
				if len(list) < 3 || len(list) > 4 {
					return nil, fmt.Errorf("if invalid args: %v", list)
				}
				expr, err := EVAL(a1, env)
				switch {
				case err != nil:
					return nil, err
				case IsTruthy(expr):
					ast = a2
					continue
				case len(list) < 4:
					return MalNil{}, nil
				default:
					ast = list[3]
					continue
				}

			case "fn*":
				// create a new function closure
				if len(list) != 3 {
					return nil, fmt.Errorf("fn* invalid args: %v", list)
				}
				binds, err := GetSlice(a1)
				if err != nil {
					return nil, err
				}
				return NewFunc(EVAL, binds, a2, env), nil

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
}

func PRINT(exp MalType) (string, error) {
	return printer.PrintStr(exp, true), nil
}

var replEnv = NewEnv()

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
	for sym, fn := range core.NS {
		replEnv.Set(sym, fn)
	}
	replEnv.Set("eval", core.MonoErrFunc(func(a MalType) (MalType, error) {
		return EVAL(a, replEnv)
	}))
	rep(`(def! not (fn* (a) (if a false true)))`)
	rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) ")")))))`)
	if len(os.Args) > 1 {
		filename := os.Args[1]
		argv := make([]MalType, len(os.Args)-2)
		if len(os.Args) > 2 {
			for i, arg := range os.Args[2:] {
				argv[i] = MalString{Value: arg}
			}
		}
		replEnv.Set("*ARGV*", NewList(argv))
		rep(`(load-file "` + filename + `")`)
		return
	}
	replEnv.Set("*ARGV*", NewListOf())
	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("user> ")
		if in.Scan() {
			read := strings.TrimSpace(in.Text())
			result, err := rep(read)
			if err != nil {
				fmt.Println("Error:", err)
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
